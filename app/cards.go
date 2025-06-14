package app

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

func cardHandler(c *gin.Context, database database.Database) ([]byte, error) {
	var card_binding common.CardIdBinding
	if err := c.ShouldBindUri(&card_binding); err != nil {
		return nil, err
	}

	// Get the post with the ID
	card, err := database.GetCard(card_binding.Id)
	if err != nil {
		return nil, err
	}
	// log.Info().Msg(card.ImageLocation)
	// Generate HTML page
	log.Info().Msgf("%v", card.JsonData)
	post_view := views.MakeCardPage(card.ImageLocation, common.Settings.AppNavbar.Links, card.JsonData)
	html_buffer := bytes.NewBuffer(nil)
	err = post_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return html_buffer.Bytes(), nil
}
