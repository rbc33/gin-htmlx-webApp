package admin_app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rbc33/gocms/database"
	"github.com/rs/zerolog/log"
)

// Since there are no builtin sets in go, we are using a map to improve the performance when checking for valid extensions
// by creating a map with the valid extensions as keys and using an existence check.
var valid_extensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// CRUD Images
// r.GET("/images/:id", getImageHandler(&database))
// r.POST("/images", postImageHandler(&database))
// r.DELETE("/images", deleteImageHandler(&database))
type AddImageRequest struct {
	Alt string `json:"alt"`
}

// func getImageHandler(database *database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		// Get the image from db
// 		()
// 	}
// }

func postImageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {

		// Get the metadata from the request body
		alt_text := c.Request.FormValue("alt")
		if alt_text == "" {
			log.Error().Msgf("alt text is required")
			return
		}

		// Check if MEDIA_DIR is defined
		err := godotenv.Load()
		if err != nil {
			log.Error().Msgf("MEDIA_DIR not defined: %v", err)
		}
		MEDIA_DIR := os.Getenv("MEDIA_DIR")

		// Begging save the file to MEDIA_DIR
		file, err := c.FormFile("file")
		if err != nil {
			log.Error().Msgf("could not upload file: %v", err)
		}

		uuid, err := uuid.New()
		if err != nil {
			log.Error().Msgf("error generating UUID: %v", err)
		}

		ext := filepath.Ext(file.Filename)
		// Check if ext is supported

		if ext == "" {
			log.Error().Msgf("could not get file extension from %s", file.Filename)
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.String(), ext)
		image_path := filepath.Join(MEDIA_DIR, filename)
		err = c.SaveUploadedFile(file, image_path)
		if err != nil {
			log.Error().Msgf("could not save the file: %v", err)
			return
		}
		// {name, alt b64}
		// Add image to db

		err = database.AddImage(uuid.String(), filename, alt_text)
		if err != nil {
			log.Error().Msgf("could not add image to db: %v", err)
			err = os.Remove(image_path)
			if err != nil {
				log.Error().Msgf("could not remove file %s: %v", image_path, err)
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id": uuid.String(),
		})
	}

}

// func deleteImageHandler(database *database.Database) func(*gin.Context) {
// 	// Delete image from db

// }
