package admin_app

import (
	"github.com/gin-gonic/gin"
	"github.com/rbc33/database"
)

func SetupRoutes(database database.Database) *gin.Engine {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	// CRUD Posts
	r.GET("/posts/:id", getPostHandler(database))
	r.POST("/posts", postPostHandler(database))
	r.PUT("/posts", putPostHandler(database))
	r.DELETE("/posts", deletePostHandler(database))

	// CRUD Images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler(database))
	// r.DELETE("/images", deleteImageHandler(&database))

	return r
}
