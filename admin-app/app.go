package admin_app

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/auth"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/middlewares"
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

	// Public routes
	r.GET("/swagger/*any", func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/swagger/") {
			c.Request.RequestURI = "/swagger/index.html"
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})
	r.POST("/register", auth.CreateRegisterHandler(database))
	r.POST("/login", auth.LoginHandler(database))

	// Protected routes group with JWT middleware
	protected := r.Group("/")
	protected.Use(middlewares.JwtAuthMiddleware()) // replace with your actual middleware function

	// Move posts routes inside protected group
	posts := protected.Group("/posts")
	{
		posts.GET("", getPostsHandler(database))
		posts.GET("/:id", getPostHandler(database))
		posts.POST("", postPostHandler(database, shortcode_handlers, post_hook.(*plugins.PostHook)))
		posts.PUT("", putPostHandler(database))
		posts.DELETE("", deletePostHandler(database))
	}

	// Move pages routes inside protected group
	pages := protected.Group("/pages")
	{
		pages.GET("", getPagesHandler(database))
		pages.POST("", postPageHandler(database))
		pages.PUT("", putPageHandler(database))
		pages.DELETE("", deletePageHandler(database))
	}

	// Similarly, move other routes inside protected group

	protected.POST("/images", postImageHandler())
	protected.DELETE("/images/:name", deleteImageHandler())

	protected.GET("/cards/:schema", getCardHandler(database))
	protected.GET("/cards/:schema/:limit/:page", getCardHandler(database))
	protected.POST("/card-schemas", postSchemaHandler(database))
	protected.GET("/card-schemas", getSchemasHandler(database))
	protected.DELETE("/card-schemas", deleteCardSchemaHandler(database))
	protected.GET("/card-schemas/:id", getSchemaHandler(database))

	protected.POST("/cards", postCardHandler(database))
	protected.PUT("/card", putCardHandler(database))
	protected.DELETE("/card", deleteCardHandler(database))
	protected.POST("/permalinks/:permalink/:post_id", postPermalinkHandler(database))
	protected.GET("/user", auth.GetCurrentUserHandler(database))

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
