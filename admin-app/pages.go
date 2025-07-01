package admin_app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rs/zerolog/log"
)

// GET /pages?offset=10&limit=10
func getPagesHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Un valor de 0 para el límite significa "sin límite".
		offsetStr := c.DefaultQuery("offset", "0")
		limitStr := c.DefaultQuery("limit", "0")

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}

		// Si no se especifica un límite, se obtienen todas las páginas.
		// Si se especifica, se usa para la paginación.
		pages, err := database.GetPages(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pages": pages})
	}
}

// @Summary      Add a new page
// @Description  Adds a new page to the database.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        page body AddPageRequest true "Page to add"
// @Success      200 {object} PageResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or data"
// @Router       /pages [post]
func postPageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_page_request AddPageRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_page_request)

		if err != nil {
			log.Warn().Msgf("invalid page request: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		err = checkRequiredPageData(add_page_request)
		if err != nil {
			log.Error().Msgf("failed to add post required data is missing: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
			return
		}

		id, err := database.AddPage(
			add_page_request.Title,
			add_page_request.Content,
			add_page_request.Link,
		)
		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusCreated, PageResponse{
			Id:   id,
			Link: add_page_request.Link,
		})
	}
}

// @Summary      Add a new page
// @Description  Adds a new page to the database.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        page body ChangePageRequest true "Page to Update"
// @Success      200 {object} PageResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or data"
// @Router       /pages [post]
func putPageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var change_page_request ChangePageRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&change_page_request)
		if err != nil {
			log.Warn().Msgf("could not get page from DB: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = checkRequiredPageData(AddPageRequest{change_page_request.Title, change_page_request.Link, change_page_request.Content})
		if err != nil {
			log.Error().Msgf("failed to add post required data is missing: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
			return
		}

		err = database.ChangePage(
			change_page_request.Id,
			change_page_request.Title,
			change_page_request.Content,
			change_page_request.Link,
		)
		if err != nil {
			log.Error().Msgf("failed to change post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not change post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": change_page_request.Id,
		})
	}
}

// @Summary      Delete a page
// @Description  Deletes a page by its Link.
// @Tags         pages
// @Produce      json
// @Param        link path string true "Page Link"
// @Success      200 {object} DeletePageRequest
// @Failure      400 {object} common.ErrorResponse "Invalid link provided"
// @Failure      404 {object} common.ErrorResponse "Page not found"
// @Router       /posts/{id} [delete]
func deletePageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_page_request DeletePageRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&delete_page_request)
		if err != nil {
			log.Warn().Msgf("could not delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = database.DeletePage(delete_page_request.Link)
		if err != nil {
			log.Error().Msgf("failed to delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not delete post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"link": delete_page_request.Link,
		})
	}
}

func checkRequiredPageData(add_page_request AddPageRequest) error {
	if strings.TrimSpace(add_page_request.Title) == "" {
		return fmt.Errorf("missing required data 'Title'")
	}

	if strings.TrimSpace(add_page_request.Content) == "" {
		return fmt.Errorf("missing required data 'Content'")
	}

	err := validateLinkRegex(add_page_request.Link)
	if err != nil {
		return err
	}
	if len(add_page_request.Link) > 255 {
		return fmt.Errorf("link must have less than 255 chars")
	}

	return nil
}

func validateLinkRegex(link string) error {
	match, err := regexp.MatchString("^[a-zA-Z0-9_\\-]+$", link)
	if err != nil {
		return fmt.Errorf("could not match the string: %v", err)
	}
	if !match {
		return fmt.Errorf("could not match the string")
	}
	return nil
}
