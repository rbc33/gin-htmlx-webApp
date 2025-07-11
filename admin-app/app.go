package admin_app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/plugins"
	lua "github.com/yuin/gopher-lua"

	_ "github.com/rbc33/gocms/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func SetupRoutes(settings common.AppSettings, shortcode_handlers map[string]*lua.LState, database database.Database, hooks map[string]plugins.Hook) *gin.Engine {
	// func SetupRoutes(settings common.AppSettings, shortcode_handlers map[string]*lua.LState, database database.Database) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.MaxMultipartMemory = 1
	r.Use(CORSMiddleware())

	post_hook, ok := hooks["add_post"]
	if !ok {
		log.Fatalf("could not find add_post hook")
	}
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
		posts.POST("", postPostHandler(database, shortcode_handlers, post_hook.(*plugins.PostHook)))
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
	// Swagger routes
	r.GET("/swagger/*any", func(c *gin.Context) {
		// Don't serve index.html again if it's already handled above
		if strings.HasSuffix(c.Request.URL.Path, "/swagger/") {
			c.Request.RequestURI = "/swagger/index.html"
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	// CRUD Images
	// r.GET("/images/:id", getImageHandler(&database))
	r.POST("/images", postImageHandler())
	r.DELETE("/images/:name", deleteImageHandler())

	r.GET("/cards/:schema", getCardHandler(database))
	r.GET("/cards/:schema/:limit/:page", getCardHandler(database))
	r.POST("/card-schemas", postSchemaHandler(database))
	r.GET("/card-schemas", getSchemasHandler(database))
	r.DELETE("/card-schemas", deleteCardSchemaHandler(database))
	r.GET("/card-schemas/:id", getSchemaHandler(database))
	// Card related stuff
	// r.GET("/card/:id", getCardHandler(database))
	r.POST("/cards", postCardHandler(database))
	r.PUT("/card", putCardHandler(database))
	r.DELETE("/card", deleteCardHandler(database))
	r.POST("/permalinks/:permalink/:post_id", postPermalinkHandler(database))

	return r
}

// CORS middleware function
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
