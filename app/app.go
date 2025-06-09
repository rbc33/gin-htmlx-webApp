package app

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

const CACHE_TIMEOUT = 20 * time.Second

type Generator = func(*gin.Context, database.Database) ([]byte, error)

func SetupRoutes(database database.Database) *gin.Engine {

	cache := makeCache(4, time.Minute*10)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.MaxMultipartMemory = 1

	// Contact form related endpoints
	r.POST("/contact-send", makeContactFormHandler())

	// Service
	r.GET("/services", makeServiceHandler())
	// r.GET("/card/:id", makeCardHandler(database))
	addCacheHandler(r, "GET", "/", homeHandler, &cache, database)
	addCacheHandler(r, "GET", "/contact", contactHandler, &cache, database)
	addCacheHandler(r, "GET", "/post/:id", postHandler, &cache, database)
	addCacheHandler(r, "GET", "/card/:id", cardHandler, &cache, database)

	r.Static("/static", "./static")

	return r
}

func addCacheHandler(e *gin.Engine, method string, endpoint string, generator Generator, cache *Cache, db database.Database) {

	handler := func(c *gin.Context) {
		// if the endpoint is cached
		cached_endpoint, err := (*cache).Get(c.Request.RequestURI)
		if err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", cached_endpoint.contents)
			return
		}

		// Before handler call (retrieve from cache)
		html_buffer, err := generator(c, db)
		if err != nil {
			log.Error().Msgf("could not generate html: %v", err)
		}
		// After handler  (add to cache)
		err = (*cache).Store(c.Request.RequestURI, html_buffer)
		if err != nil {
			log.Warn().Msgf("could not add page to cache: %v", err)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", html_buffer)
	}
	// Hacky
	if method == "GET" {
		e.GET(endpoint, handler)
	}
	if method == "POST" {
		e.POST(endpoint, handler)
	}
	if method == "DELETE" {
		e.DELETE(endpoint, handler)
	}
	if method == "PUT" {
		e.PUT(endpoint, handler)
	}
}

// This function will act as the handler for
// the home page
func homeHandler(c *gin.Context, db database.Database) ([]byte, error) {
	posts, err := db.GetPosts()
	if err != nil {
		return nil, err
	}

	// if not cached, create the cache
	// index_view := views.MakeIndex(posts)
	index_view := views.MakeIndex(posts)
	html_buffer := bytes.NewBuffer(nil)
	err = index_view.Render(c, html_buffer)
	if err != nil {
		return nil, err
	}
	return html_buffer.Bytes(), nil
}
