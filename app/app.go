package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/database"
	views "github.com/rbc33/views/index"
	"github.com/rbc33/views/tailwind"
	"github.com/rs/zerolog/log"
)

func Run(database database.Database) error {
	r := gin.Default()
	r.MaxMultipartMemory = 1
	// r.LoadHTMLGlob("views/**/*")

	r.GET("/", makeHomeHandler(database))
	r.GET("/tailwind", makeHomeHandlertailwind(database))

	// Contact form related endpoints
	r.GET("/contact", makeContactPageHandler())
	r.POST("/contact-send", makeContactFormHandler())

	// Post related endpoints
	r.GET("/post/:id", makePostHandler(database))

	// Service
	r.GET("/services", makeServiceHandler())

	r.Static("/static", "./static")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	return nil
}

// / This function will act as the handler for
// / the home page
func makeHomeHandler(db database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		posts, err := db.GetPosts()
		if err != nil {
			log.Error().Msgf("error loading posts: %v\n", err)
			return
		}
		TemplRender(c, http.StatusOK, views.MakeIndex(posts))
		// c.HTML(http.StatusOK, "", views.MakeIndex(posts))
	}
}
func makeHomeHandlertailwind(db database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		posts, err := db.GetPosts()
		if err != nil {
			log.Error().Msgf("error loading posts: %v\n", err)
			return
		}
		TemplRender(c, http.StatusOK, tailwind.MakeIndex(posts))
		// c.HTML(http.StatusOK, "", views.MakeIndex(posts))
	}
}
