package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusAccepted, "index", gin.H{})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MY_SQL_URL := os.Getenv("MY_SQL_URL")
	db, err := sql.Open("mysql", MY_SQL_URL)
	if err != nil {
		log.Fatalf("couldn't connect to DB: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	r := gin.Default()
	r.MaxMultipartMemory = 1

	// Load HTML templates
	// r.LoadHTMLFiles("templates/contact-success.html", "templates/contact-failure.html")
	r.LoadHTMLFiles("templates/index.html")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", homeHandler)

	// Contact form endpoint
	r.POST("/contact-send", func(c *gin.Context) {
		c.Request.ParseForm()
		name := c.Request.FormValue("name")
		email := c.Request.FormValue("email")
		message := c.Request.FormValue("message")

		// parse email
		_, err := mail.ParseAddress(email)
		if err != nil {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "Invalid email address",
			})
			return
		}
		//make sure that name and message are reasonable
		if len(name) > 200 {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "Name is too long",
			})
			return
		}
		if len(message) > 1000 {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "message too long",
			})
			return
		}

		c.HTML(http.StatusOK, "contact-success.html", gin.H{
			"name":  name,
			"email": email,
		})

	})
	r.Static("/static", "./static")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
