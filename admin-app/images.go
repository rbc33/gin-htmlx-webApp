package admin_app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/rbc33/gocms/common"
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

// func getImageHandler(database *database.Database) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		// Get the image from db
// 		()
// 	}
// }

func resizeImage(srcPath string, width uint) error {
	// Open the source file
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("could not open source image: %v", err)
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("could not decode image: %v", err)
	}

	// Calculate height to maintain aspect ratio
	bounds := img.Bounds()
	ratio := float64(bounds.Dy()) / float64(bounds.Dx())
	height := uint(float64(width) * ratio)

	// Resize
	resized := resize.Resize(width, height, img, resize.Lanczos3)

	// Create new file
	out, err := os.Create(srcPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer out.Close()

	// Save based on format
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, resized, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(out, resized)
	case "gif":
		// Note: GIF will lose animation
		err = png.Encode(out, resized)
	}

	if err != nil {
		return fmt.Errorf("could not encode resized image: %v", err)
	}

	return nil
}

func postImageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {

		// Get the metadata from the request body
		alt_text := c.Request.FormValue("alt")
		if alt_text == "" {
			log.Error().Msgf("alt text is required")
			return
		}

		// Check if MEDIA_DIR is defined

		MEDIA_DIR := common.Settings.ImageDirectory

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
		_, ok := valid_extensions[ext]
		if !ok {
			log.Error().Msgf("file extension is not supported %s", ext)
		}

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

		// Resize image to 477px width
		err = resizeImage(image_path, 477)
		if err != nil {
			log.Error().Msgf("could not resize image: %v", err)
			os.Remove(image_path)
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

func deleteImageHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_image_request DeleteImageBinding
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&delete_image_request)
		if err != nil {
			log.Warn().Msgf("could not delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}
		delete_image_split := strings.Split(delete_image_request.Name, ".")
		err = database.DeleteImage(delete_image_split[0])
		if err != nil {
			log.Error().Msgf("failed to delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not delete post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": delete_image_request.Name,
		})
	}
}
