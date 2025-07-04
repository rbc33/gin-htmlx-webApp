package app

import (
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func makeWebHookHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		if c.GetHeader("X-Webhook-Secret") != os.Getenv("GIT_SECRET") {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		prc := exec.Command("git", "pull", "origin", "main")
		err := prc.Run()
		if err != nil {
			log.Error().Msgf("Error pulling the repo: %v", err)
		}

		if prc.ProcessState.Success() {
			log.Info().Msgf("Process success with output: \n%s", prc.ProcessState.String())
		}
	}
}
