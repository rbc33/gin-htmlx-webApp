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
	r.GET("/posts/:id", getPostHandler(database))
	// r.POST("/posts", postPostHandler(database, shortcode_handlers, post_hook.(plugins.PostHook)))
	r.POST("/posts", postPostHandler(database, shortcode_handlers))
	r.PUT("/posts", putPostHandler(database))
	r.DELETE("/posts", deletePostHandler(database))

	// CRUD Images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler())
	r.DELETE("/images/:name", deleteImageHandler())

	r.POST("/pages", postPageHandler(database))

	// Card related stuff
	r.GET("/card/:id", getCardHandler(database))
	r.POST("/card", postCardHandler(database))
	r.PUT("/card", putCardHandler(database))
	r.DELETE("/card", deleteCardHandler(database))

	return r
}
