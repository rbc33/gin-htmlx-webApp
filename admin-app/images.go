package admin_app

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/metadata"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
)

var allowed_extensions = map[string]bool{
	".jpeg": true, ".jpg": true, ".png": true, ".heic": true,
}

var allowed_content_types = map[string]bool{
	"image/jpeg": true, "image/png": true, "image/gif": true, "image/heic": true,
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

func resizeImage(srcPath string, width int, height int) error {
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

	// Calculate new size
	bounds := img.Bounds()
	ratio := float64(width) / float64(bounds.Dx())
	newHeight := int(float64(bounds.Dy()) * ratio)
	if newHeight > height {
		newHeight = height
		width = int(float64(newHeight) / float64(bounds.Dy()) * float64(bounds.Dx()))
	}

	// Create new image

	dst := image.NewRGBA(image.Rect(0, 0, width, newHeight))

	// Resize using golang.org/x/image/draw
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Create new file
	out, err := os.Create(srcPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer out.Close()

	// Save based on format
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, dst, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(out, dst)
	case "gif":
		// Note: GIF will lose animation
		err = png.Encode(out, dst)
	}

	if err != nil {
		return fmt.Errorf("could not encode resized image: %v", err)
	}

	return nil
}

// @Summary      Upload a new image
// @Description  Uploads an image file, saves it, and creates minified versions.
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file formData file true "The image file to upload"
// @Param        excerpt formData string false "A brief description of the image"
// @Success      200 {object} ImageIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid input, file type, or size"
// @Failure      500 {object} common.ErrorResponse "Server error while saving file"
// @Router       /images [post]
func postImageHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1000000)
		form, err := c.MultipartForm()
		if err != nil {
			log.Error().Msgf("could not create multipart form: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("request type must be `multipart-form`", err))
			return
		}

		file_array := form.File["file"]
		if len(file_array) == 0 || file_array[0] == nil {
			log.Error().Msgf("could not get the file array: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no file provided for image upload",
			})
			return
		}

		file := file_array[0]
		file_content_type := file.Header.Get("content-type")
		_, ok := allowed_content_types[file_content_type]
		if !ok {
			log.Error().Msgf("file type not supported")
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("file type not supported"))
			return
		}

		detected_content_type, err := getContentType(file)
		if err != nil || detected_content_type != file_content_type {
			log.Error().Msgf("the provided file does not match the provided content type")
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("provided file content is not allowed"))
			return
		}

		uuid, err := uuid.New()
		if err != nil {
			log.Error().Msgf("could not create the UUID: %v", err)
			c.JSON(http.StatusInternalServerError, common.ErrorRes("cannot create unique identifier", err))
			return
		}

		ext := filepath.Ext(file.Filename)
		// check ext is supported
		_, ok = allowed_extensions[ext]
		if ext == "" || !ok {
			log.Error().Msgf("file extension is not supported %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("file extension is not supported", err))
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.String(), ext)
		image_path := filepath.Join(common.Settings.ImageDirectory, filename)
		err = c.SaveUploadedFile(file, image_path)
		if err != nil {
			log.Error().Msgf("could not save file: %v", err)
			c.JSON(http.StatusInternalServerError, common.ErrorRes("failed to upload image", err))
			return
		}

		// Generate Json from metadata
		excerpt_text_array := form.Value["excerpt"]
		excerpt := "unknown"
		if len(excerpt_text_array) > 0 {
			excerpt = excerpt_text_array[0]
		}
		name := file.Filename[:len(file.Filename)-len(ext)]
		metadata.GenerateJson(filename, name, excerpt)

		// Resize image to 477px width
		err = resizeImage(image_path, 477, 620)
		if err != nil {
			log.Error().Msgf("could not resize image: %v", err)
			os.Remove(image_path)
			return
		}

		// End saving to filesystem
		c.JSON(http.StatusOK, ImageIdResponse{
			Id: uuid.String(),
		})
	}
}

// @Summary      Delete an image
// @Description  Deletes an image file from the server by its filename.
// @Tags         images
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name path string true "Image filename to delete"
// @Success      200 {object} ImageIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid or missing filename"
// @Router       /images/{name} [delete]
func deleteImageHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_image_binding DeleteImageRequest
		err := c.ShouldBindUri(&delete_image_binding)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorRes("no id provided to delete image", err))
			return
		}

		image_path := filepath.Join(common.Settings.ImageDirectory, delete_image_binding.Name)
		ext := filepath.Ext(image_path)
		// Fix json_name calculation: replace extension with ".json"
		json_name := image_path[:len(image_path)-len(ext)] + ".json"

		err = os.Remove(image_path)
		if err != nil {
			log.Warn().Msgf("could not delete stored image file: %v", err)
			// No return because we have to remove the database entry nonetheless.
		}
		err = os.Remove(json_name)
		if err != nil {
			log.Warn().Msgf("could not delete stored json file: %v", err)
			// No return because we have to remove the database entry nonetheless.
		}

		c.JSON(http.StatusOK, ImageIdResponse{
			delete_image_binding.Name,
		})
	}
}

func getContentType(file_header *multipart.FileHeader) (string, error) {
	// Check if the content matches the provided type.
	image_file, err := file_header.Open()
	if err != nil {
		log.Error().Msgf("could not open file for check.")
		return "", err
	}

	// According to the documentation only the first `512` bytes are required for verifying the content type
	tmp_buffer := make([]byte, 512)
	_, read_err := image_file.Read(tmp_buffer)
	if read_err != nil {
		log.Error().Msgf("could not read into temp buffer")
		return "", read_err
	}
	return getContentTypeFromData(tmp_buffer), nil
}

func getContentTypeFromData(data []byte) string {
	return http.DetectContentType(data[:512])
}
