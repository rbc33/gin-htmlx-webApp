package app

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	views "github.com/rbc33/views/contact"
)

func makeContactFormHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Request.ParseForm()
		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")

		// Parse email
		_, err := mail.ParseAddress(email)
		if err != nil {
			TemplRender(c, http.StatusOK, views.MakeContactFailure(email, err.Error()))
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			TemplRender(c, http.StatusOK, views.MakeContactFailure(email, "name too long (200 char max)"))
			return
		}

		if len(message) > 10000 {
			TemplRender(c, http.StatusOK, views.MakeContactFailure(email, "message too long (10000 char max)"))

			return
		}

		TemplRender(c, http.StatusOK, views.MakeContactSuccess(email, name))

	}
}

// TODO : This is a duplicate of the index handler... abstract
// func makeContactPageHandler(db database.Database) func(*gin.Context) {
func makeContactPageHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		// 	posts, err := db.GetPosts()
		// 	if err != nil {
		// 		log.Error().Msgf("error loading posts: %v\n", err)
		// 		return
		// 	}

		// c.HTML(http.StatusAccepted, "contact", gin.H{"posts": posts})
		TemplRender(c, http.StatusOK, views.MakeContactPage())
	}
}
