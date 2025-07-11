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
	"github.com/rbc33/gocms/plugins"
	"github.com/rs/zerolog/log"
	lua "github.com/yuin/gopher-lua"
)

// @Summary      Get a list of posts
// @Description  Retrieves a paginated list of posts.
// @Tags         posts
// @Produce      json
// @Param        offset query int false "Offset for pagination" default(0)
// @Param        limit  query int false "Limit for pagination. 0 means no limit." default(0)
// @Success      200 {object} GetPostsResponse
// @Failure      400 {object} common.ErrorResponse "Invalid offset or limit parameter"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @Router       /posts [get]
func getPostsHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lee offset y limit de la query (?offset=0&limit=10)
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

		// Si no se especifica un límite, se obtienen todos los posts.
		// Si se especifica, se usa para la paginación.
		posts, err := database.GetPosts(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, GetPostsResponse{Posts: posts})
	}
}

// @Summary      Get a single post
// @Description  Retrieves a single post by its ID.
// @Tags         posts
// @Produce      json
// @Param        id path int true "Post ID"
// @Success      200 {object} GetPostResponse
// @Failure      400 {object} common.ErrorResponse "Invalid post ID"
// @Failure      404 {object} common.ErrorResponse "Post not found"
// @Router       /posts/{id} [get]
func getPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// localhost:8080/post/{id}
		var post_binding common.PostIdBinding
		if err := c.ShouldBindUri(&post_binding); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not get post id",
				"msg":   err.Error(),
			})
			return
		}

		post, err := database.GetPost(post_binding.Id)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post id not found",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, GetPostResponse{
			Id:      post.Id,
			Title:   post.Title,
			Excerpt: post.Excerpt,
			Content: post.Content,
		})
	}
}

// @Summary      Add a new post
// @Description  Adds a new post to the database.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post body AddPostRequest true "Post to add"
// @Success      201 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or missing data"
// @Router       /post [post]
func postPostHandler(database database.Database, shortcode_handlers map[string]*lua.LState, post_hook *plugins.PostHook) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_post_request AddPostRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_post_request)

		if err != nil {
			log.Warn().Msgf("invalid post request: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		err = checkRequiredData(add_post_request)
		if err != nil {
			log.Error().Msgf("failed to add post required data is missing: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
		}

		altered_post := post_hook.UpdatePost(add_post_request.Title, add_post_request.Excerpt, add_post_request.Content, shortcode_handlers)

		fmt.Print("Title: ", altered_post.Title)
		fmt.Print("Excerpt: ", altered_post.Excerpt)
		fmt.Print("Content: ", altered_post.Content)

		id, err := database.AddPost(
			altered_post.Title,
			altered_post.Excerpt,
			altered_post.Content,
		)
		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusCreated, PostIdResponse{
			id,
		})
	}
}

// @Summary      Update an existing post
// @Description  Updates an existing post with new data.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post body ChangePostRequest true "Post data to update"
// @Success      200 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or could not change post"
// @Router       /posts [put]
func putPostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var change_post_request ChangePostRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&change_post_request)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = database.ChangePost(
			change_post_request.Id,
			change_post_request.Title,
			change_post_request.Excerpt,
			change_post_request.Content,
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
			"id": change_post_request.Id,
		})
	}
}

// @Summary      Delete a post
// @Description  Deletes a post by its ID.
// @Tags         posts
// @Produce      json
// @Param        id body DeletePostRequest true "Post ID to delete"
// @Success      200 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid ID provided"
// @Failure      404 {object} common.ErrorResponse "Post not found"
// @Router       /posts [delete]
func deletePostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_post_request DeletePostRequest
		if err := c.ShouldBindJSON(&delete_post_request); err != nil {
			log.Warn().Msgf("could not delete post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		err := database.DeletePost(delete_post_request.Id)
		if err != nil {
			log.Error().Msgf("failed to delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not delete post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": delete_post_request.Id,
		})
	}
}

func checkRequiredData(addPostRequest AddPostRequest) error {
	if strings.TrimSpace(addPostRequest.Title) == "" {
		return fmt.Errorf("missing required data 'Title'")
	}

	if strings.TrimSpace(addPostRequest.Excerpt) == "" {
		return fmt.Errorf("missing required data 'Excerpt'")
	}

	if strings.TrimSpace(addPostRequest.Content) == "" {
		return fmt.Errorf("missing required data 'Content'")
	}

	return nil
}

// partitionString will partition the strings by
// removing the given ranges
func partitionString(text string, indexes [][]int) []string {

	if len(text) == 0 {
		return []string{}
	}

	partitions := make([]string, 0)
	start := 0
	for _, window := range indexes {
		partitions = append(partitions, text[start:window[0]])
		start = window[1]
	}

	partitions = append(partitions, text[start:len(text)-1])
	return partitions
}

func shortcodeToMarkdown(shortcode string, shortcode_handlers map[string]*lua.LState) (string, error) {
	key_value := strings.Split(shortcode, ":")

	key := key_value[0]
	values := key_value[1:]

	if handler, found := shortcode_handlers[key]; found {

		// Need to quote all values for a valid lua syntax
		quoted_values := make([]string, 0)
		for _, value := range values {
			quoted_values = append(quoted_values, fmt.Sprintf("%q", value))
		}

		err := handler.DoString(fmt.Sprintf(`result = HandleShortcode({%s})`, strings.Join(quoted_values, ",")))
		if err != nil {
			return "", fmt.Errorf("error running %s shortcode: %v", key, err)
		}

		value := handler.GetGlobal("result")
		if ret_type := value.Type().String(); ret_type != "string" {
			return "", fmt.Errorf("error running %s shortcode: invalid return type %s", key, ret_type)
		} else if ret_type == "" {
			return "", fmt.Errorf("error running %s shortcode: returned empty string", key)
		}

		return value.String(), nil
	}

	return "", fmt.Errorf("unsupported shortcode: %s", key)
}

func transformContent(content string, shortcode_handlers map[string]*lua.LState) (string, error) {
	// Find all the occurences of {{ and }}
	regex, _ := regexp.Compile(`{{[\w.-]+(:[\w.-]+)+}}`)

	shortcodes := regex.FindAllStringIndex(content, -1)
	if len(shortcodes) == 0 {
		return content, nil
	}

	partitions := partitionString(content, shortcodes)

	builder := strings.Builder{}
	i := 0

	for i, shortcode := range shortcodes {
		builder.WriteString(partitions[i])

		markdown, err := shortcodeToMarkdown(content[shortcode[0]+2:shortcode[1]-2], shortcode_handlers)
		if err != nil {
			log.Error().Msgf("%v", err)
			markdown = ""
		}
		builder.WriteString(markdown)
	}

	// Guaranteed to have +1 than the number of
	// shortcodes by algorithm
	builder.WriteString(partitions[i+1])

	return builder.String(), nil
}
