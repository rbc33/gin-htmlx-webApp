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

func getPagesHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		offsetStr := c.DefaultQuery("offset", "")
		limitStr := c.DefaultQuery("limit", "")

		var offset, limit int
		var err error

		// Si offset o limit no están presentes, usa -1 para indicar "sin límite"
		if offsetStr == "" {
			offset = -1
		} else {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
				return
			}
		}

		if limitStr == "" {
			limit = -1
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
				return
			}
		}
		pages, err := database.GetPages(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pages": pages})
	}
}

// postPageHandler if the function handling the endpoint for adding pages
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

// func putPostHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		var change_post_request ChangePostRequest
// 		decoder := json.NewDecoder(c.Request.Body)
// 		decoder.DisallowUnknownFields()

// 		err := decoder.Decode(&change_post_request)
// 		if err != nil {
// 			log.Warn().Msgf("could not get post from DB: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "invalid request body",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		err = database.ChangePost(
// 			change_post_request.Id,
// 			change_post_request.Title,
// 			change_post_request.Excerpt,
// 			change_post_request.Content,
// 		)
// 		if err != nil {
// 			log.Error().Msgf("failed to change post: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "could not change post",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"id": change_post_request.Id,
// 		})
// 	}
// }

// func deletePostHandler(database database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		var delete_post_request DeletePostBinding
// 		decoder := json.NewDecoder(c.Request.Body)
// 		decoder.DisallowUnknownFields()

// 		err := decoder.Decode(&delete_post_request)
// 		if err != nil {
// 			log.Warn().Msgf("could not delete post: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "invalid request body",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		err = database.DeletePost(delete_post_request.Id)
// 		if err != nil {
// 			log.Error().Msgf("failed to delete post: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "could not delete post",
// 				"msg":   err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"id": delete_post_request.Id,
// 		})
// 	}
// }

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
