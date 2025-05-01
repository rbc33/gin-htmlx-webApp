package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/database"
	"github.com/rs/zerolog/log"
)

type PostBinding struct {
	Id string `uri:"id" binding:"required"`
}

func makePostHandler(db database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// Get post with the ID
		var post_binding PostBinding
		if err := c.ShouldBindUri(&post_binding); err != nil {
			// TODO redo this error to serve error page
			c.JSON(400, gin.H{"msg": err})
			return
		}
		// Get the post with the ID
		post_id, err := strconv.Atoi(post_binding.Id)
		if err != nil {
			// TODO redo this error to serve error page
			c.JSON(400, gin.H{"msg": err})
			return
		}

		post, err := db.GetPost(post_id)
		if err != nil {
			// TODO redo this error to serve error page
			c.JSON(400, gin.H{"msg": err})
			return
		}

		// Markdown to HTML

		// Serve the templated page here
		log.Warn().Msgf("Post: %v", post)
		c.HTML(http.StatusOK, "post", post)
	}
}
