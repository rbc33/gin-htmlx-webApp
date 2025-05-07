package app

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/rbc33/database"
	"github.com/rs/zerolog/log"
)

type PostBinding struct {
	Id string `uri:"id" binding:"required"`
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
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
		post.Content = string(mdToHTML([]byte(post.Content)))

		// Serve the templated page here
		log.Warn().Msgf("Post: %v", post)
		c.HTML(http.StatusOK, "post", gin.H{
			"Title":   post.Title,
			"Content": template.HTML(post.Content),
		})
	}
}
