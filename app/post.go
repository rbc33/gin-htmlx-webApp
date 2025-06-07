package app

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

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
		return nil, err
	}

	// Get the post with the ID
	post, err := database.GetPost(post_binding.Id)
	if err != nil {
		return nil, err
	}

	// Generate HTML page
	post.Content = string(mdToHTML([]byte(post.Content)))
	// post_view := views.MakePostPage(post.Title, post.Content)
	post_view := views.MakePostPage(post.Title, post.Content)
	html_buffer := bytes.NewBuffer(nil)
	err = post_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return html_buffer.Bytes(), nil
}
