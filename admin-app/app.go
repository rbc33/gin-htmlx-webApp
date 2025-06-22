package admin_app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	lua "github.com/yuin/gopher-lua"
)

func LoadShortcodesHandlers(shortcodes []common.Shortcode) (map[string]*lua.LState, error) {
	shorcodes_handlers := make(map[string]*lua.LState, 0)
	for _, shortcode := range shortcodes {
		// Read the Lua State
		state := lua.NewState()
		err := state.DoFile(shortcode.Plugin)
		if err != nil {
			return map[string]*lua.LState{}, fmt.Errorf("could not load %s: %v", shortcode.Name, err)
		}
		shorcodes_handlers[shortcode.Name] = state
	}
	return shorcodes_handlers, nil
}

// func SetupRoutes(settings common.AppSettings, shortcode_handlers map[string]*lua.LState, database database.Database, hooks map[string]plugins.Hook) *gin.Engine {
func SetupRoutes(settings common.AppSettings, shortcode_handlers map[string]*lua.LState, database database.Database) *gin.Engine {

	r := gin.Default()
	r.MaxMultipartMemory = 1

	// post_hook, ok := hooks["add_post"]
	// if !ok {
	// 	log.Fatalf("could not find add_post hook")
	// }

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin pong"})
	})

	// CRUD Posts
	// Group posts routes and fix ordering
	posts := r.Group("/posts")
	{
		// GET /?offset=10&limit=10
		posts.GET("", getPostsHandler(database))    // GET /posts
		posts.GET("/:id", getPostHandler(database)) // GET /posts/:id
		posts.POST("", postPostHandler(database, shortcode_handlers))
		posts.PUT("", putPostHandler(database))
		posts.DELETE("", deletePostHandler(database))
	}
	// CRUD Pages
	pages := r.Group("/pages")
	{
		// GET /pages?offset=10&limit=10
		pages.GET("", getPagesHandler(database)) // GET /pages
		pages.POST("", postPageHandler(database))
		pages.PUT("", putPageHandler(database))
		pages.DELETE("", deletePageHandler(database))
	}

	// CRUD Images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler())
	r.DELETE("/images/:name", deleteImageHandler())

	// Card related stuff
	r.GET("/card/:id", getCardHandler(database))
	r.POST("/card", postCardHandler(database))
	r.PUT("/card", putCardHandler(database))
	r.DELETE("/card", deleteCardHandler(database))

	return r
}
