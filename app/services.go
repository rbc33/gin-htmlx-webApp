package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

// This fumction will act as the handles for the home page
func makeServiceHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		log.Info().Msg("Testing")
		// TemplRender(c, http.StatusOK, views.MakeServicesPage())
		err := TemplRender(c, http.StatusOK, views.MakeServicesPage(common.Settings.AppNavbar.Links))
		if err != nil {
			log.Error().Msgf("%s", err)
		}

	}
}
