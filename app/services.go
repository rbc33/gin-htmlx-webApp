package app

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// This fumction will act as the handles for the home page
func makeServiceHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		log.Info().Msg("Testing")
		c.HTML(http.StatusAccepted, "services", gin.H{})
	}
} 
