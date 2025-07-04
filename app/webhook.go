package app

import (
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func makeWebHookHandler() func(*gin.Context) {
	return func(c *gin.Context) {
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
