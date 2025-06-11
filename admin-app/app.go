package admin_app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
)

func SetupRoutes(settings common.AppSettings, database database.Database) *gin.Engine {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin pong"})
	})

	// CRUD Posts
	r.GET("/posts/:id", getPostHandler(database))
	r.POST("/posts", postPostHandler(database))
	r.PUT("/posts", putPostHandler(database))
	r.DELETE("/posts", deletePostHandler(database))

	// CRUD Images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler(database))
	// r.DELETE("/images", deleteImageHandler(&database))

	// Card related stuff
	r.GET("/card/:id", getCardHandler(database))
	r.POST("/card", postCardHandler(database))
	r.PUT("/card", putCardHandler(database))
	r.DELETE("/card", deleteCardHandler(database))

	return r
}
