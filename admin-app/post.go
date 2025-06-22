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
	lua "github.com/yuin/gopher-lua"
)

func getPostsHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lee offset y limit de la query (?offset=0&limit=10)
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

		posts, err := database.GetPosts(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, GetPostsResponse{Posts: posts})
	}
}

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

// func postPostHandler(database database.Database, shortcode_handlers map[string]*lua.LState, post_hook plugins.PostHook) func(*gin.Context) {
func postPostHandler(database database.Database, shortcode_handlers map[string]*lua.LState) func(*gin.Context) {
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

		// TODO : Move this into the plugin
		// altered_post := post_hook.UpdatePost(add_post_request.Title, add_post_request.Content, add_post_request.Excerpt, shortcode_handlers)

		transformed_content, err := transformContent(add_post_request.Content, shortcode_handlers)
		if err != nil {
			log.Warn().Msgf("could not transform post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		fmt.Print("Title: ", add_post_request.Title)
		fmt.Print("Excerpt: ", add_post_request.Excerpt)
		fmt.Print("Content: ", transformed_content)

		id, err := database.AddPost(
			add_post_request.Title,
			add_post_request.Excerpt,
			transformed_content,
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

func deletePostHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_post_request DeletePostBinding
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&delete_post_request)
		if err != nil {
			log.Warn().Msgf("could not delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = database.DeletePost(delete_post_request.Id)
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
