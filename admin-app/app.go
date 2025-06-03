package admin_app

import (
	"github.com/gin-gonic/gin"
	"github.com/rbc33/database"
	"github.com/rs/zerolog/log"
)

func Run(database database.Database) error {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	// CRUD Posts
	r.GET("/posts/:id", getPostHandler(&database))
	r.POST("/posts", postPostHandler(&database))
	r.PUT("/posts", putPostHandler(&database))
	r.DELETE("/posts", deletePostHandler(&database))

	// CRUD Images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler(&database))
	// r.DELETE("/images", deleteImageHandler(&database))

	err := r.Run(":8081")
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		return err
	}

	return nil
}
