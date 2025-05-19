package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/database"
	"github.com/rs/zerolog/log"
)

func Run(database database.Database) error {
	r := gin.Default()
	r.MaxMultipartMemory = 1
	//r.LoadHTMLFiles("./templates/contact/contact-success.html", "./templates/contact/contact-failure.html")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", makeHomeHandler(database))

	// Contact form related endpoints
	r.GET("/contact", makeContactPageHandler(database))
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

		c.HTML(http.StatusAccepted, "index", gin.H{"posts": posts})
	}
}
