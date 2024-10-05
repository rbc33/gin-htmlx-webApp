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

type Post struct {
	Title   string
	Content string
}
type Database struct {
	address    string
	connection *sql.DB
}

func (db Database) getPost() ([]Post, error) {
	rows, err := db.connection.Query("SELECT title, content FROM posts")
	if err != nil {
		return make([]Post, 0), err
	}
	defer rows.Close()
	all_posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		// Usa & para pasar los punteros a las variables
		if err = rows.Scan(&post.Title, &post.Content); err != nil {
			return make([]Post, 0), err
		}
		all_posts = append(all_posts, post)
	}

	// Verifica si hubo errores en el iterador
	if err = rows.Err(); err != nil {
		return make([]Post, 0), err
	}

	return all_posts, nil
}

func makeHomeHandler(db Database) func(c *gin.Context) {
	return func(c *gin.Context) {
		posts, err := db.getPost()
		if err != nil {
			log.Fatal("Error ejecutando db.getPost()", err)
		}
		c.HTML(http.StatusAccepted, "index", gin.H{"posts": posts})
	}
}
func makeDatabase() Database {
	// load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// set the db conf
	MY_SQL_URL := os.Getenv("MY_SQL_URL")
	db, err := sql.Open("mysql", MY_SQL_URL)
	if err != nil {
		log.Fatalf("couldn't connect to DB: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return Database{
		address:    MY_SQL_URL,
		connection: db,
	}
}

func main() {
	db_connection := makeDatabase()
	r := gin.Default()
	r.MaxMultipartMemory = 1

	// Load HTML templates
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", makeHomeHandler(db_connection))

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
