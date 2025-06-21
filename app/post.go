package app

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

func serveErrorPage(c *gin.Context, err string, error_code int) error {
	error_view := views.MakeErrorPage(err, common.Settings.AppNavbar.Links)
	if err := TemplRender(c, error_code, error_view); err != nil {
		log.Error().Msgf("Could not render: %v", err)
	}
	return nil
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
func postHandler(c *gin.Context, database database.Database) ([]byte, error) {

	var post_binding common.PostIdBinding
	if err := c.ShouldBindUri(&post_binding); err != nil {
		err = serveErrorPage(c, "requested invalid post ID", http.StatusBadRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		}
		return nil, err
	}

	// Get the post with the ID
	post, err := database.GetPost(post_binding.Id)
	if err != nil || post.Content == "" {
		err = serveErrorPage(c, "post not found", http.StatusNotFound)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post not found"})
		}
		return nil, err

	}

	// Generate HTML page
	post.Content = string(mdToHTML([]byte(post.Content)))
	post_view := views.MakePostPage(post.Title, common.Settings.AppNavbar.Links, post.Content)
	html_buffer := bytes.NewBuffer(nil)
	err = post_view.Render(c, html_buffer)
	if err != nil {
		err = serveErrorPage(c, "error generating the HTML", http.StatusInternalServerError)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error generating the HTML"})
		}
		return nil, err
	}

	return html_buffer.Bytes(), nil
}
