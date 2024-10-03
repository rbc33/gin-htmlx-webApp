package main

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLFiles("templates/contact-success.html")
	// , "templates/contact-failure.html")

	r.MaxMultipartMemory = 1
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Contact form endpoint
	r.POST("/contact-send", func(c *gin.Context) {
		c.Request.ParseForm()
		name := c.Request.FormValue("name")
		email := c.Request.FormValue("email")
		message := c.Request.FormValue("message")

		// parse email
		_, err := mail.ParseAddress(email)
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-failure.html", gin.H{
				"email": email,
				"error": "Invalid email address",
			})
		}
		//make sure that name and message are reasonable
		if len(name) > 200 {
			c.HTML(http.StatusBadRequest, "contact-failure.html", gin.H{
				"email": email,
				"error": "Name is too long",
			})
		}
		if len(message) > 1000 {
			c.HTML(http.StatusBadRequest, "contact-failure.html", gin.H{
				"email": email,
				"error": "Invalid email address",
			})
		}

		c.HTML(http.StatusOK, "contact-success.html", gin.H{
			"name":  name,
			"email": email,
		})

	})
	r.Static("/templates", "./templates")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
