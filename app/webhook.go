package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// test
func makeWebHookHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		secret := os.Getenv("GIT_SECRET")
		signature := c.GetHeader("X-Hub-Signature-256")
		if signature == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing signature"})
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "could not read body"})
			return
		}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(signature), []byte(expected)) {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid signature"})
			return
		}
		prePrc := exec.Command("git", "stash")
		err = prePrc.Run()
		if err != nil {
			log.Error().Msgf("Error pulling the repo: %v", err)
		}

		if prePrc.ProcessState.Success() {
			log.Info().Msgf("Preprocess success: %s", prePrc.ProcessState.String())
		}
		prc := exec.Command("git", "pull", "origin", "master", "--rebase")
		err = prc.Run()
		if err != nil {
			log.Error().Msgf("Error pulling the repo: %v", err)
		}

		if prc.ProcessState.Success() {
			log.Info().Msgf("Process success: %s", prc.ProcessState.String())
		}

		permProc := exec.Command("chmod", "+x", "admin-run.sh")
		err = permProc.Run()
		if err != nil {
			log.Error().Msgf("Error pulling the repo: %v", err)
		}

		if permProc.ProcessState.Success() {
			log.Info().Msgf("Permission changed success: %s", permProc.ProcessState.String())
		}

		c.JSON(200, gin.H{"status": "ok"})
	}
}
