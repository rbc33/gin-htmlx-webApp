package app

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
	"github.com/rs/zerolog/log"
)

const CACHE_TIMEOUT = 20 * time.Second

type Generator = func(*gin.Context, database.Database) ([]byte, error)

// func permalinkPostHandler(c *gin.Context, app_settings common.AppSettings, db database.Database) ([]byte, error) {
func permalinkPostHandler(post_id int) func(*gin.Context, database.Database) ([]byte, error) {
	return func(c *gin.Context, database database.Database) ([]byte, error) {
		c.Params = append(c.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", post_id)})
		return postHandler(c, database)
	}
}

func SetupRoutes(settings common.AppSettings, database database.Database) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.MaxMultipartMemory = 1

	// Contact form related endpoints
	r.POST("/contact-send", makeContactFormHandler())
	r.POST("/webhook", makeWebHookHandler())

	// All cache endpoints
	cache := MakeCache(4, time.Minute*10, &TimeValidator{})
	addCacheHandler(r, "GET", "/", homeHandler, &cache, database)
	addCacheHandler(r, "GET", "/contact", contactHandler, &cache, database)
	addCacheHandler(r, "GET", "/about", aboutHandler, &cache, database)
	addCacheHandler(r, "GET", "/services", servicesHandler, &cache, database)
	addCacheHandler(r, "GET", "/post/:id", postHandler, &cache, database)
	addCacheHandler(r, "GET", "/products/:schema", productHandler, &cache, database)
	addCacheHandler(r, "GET", "/products/", getSchemasHandler, &cache, database)
	addCacheHandler(r, "GET", "/images/:name", imageHandler, &cache, database)
	addCacheHandler(r, "GET", "/images", imagesHandler, &cache, database)
	addCacheHandler(r, "GET", "/gallery/:name", galleryHandler, &cache, database)

	// Pages will be querying the page content from the unique
	// link given at the creation of the page step
	addCacheHandler(r, "GET", "/pages", getPagesHandler, &cache, database)
	addCacheHandler(r, "GET", "/pages/:num", getPagesHandler, &cache, database)
	addCacheHandler(r, "GET", "/page/:link", pageHandler, &cache, database)

	// Add the pagination route as a cacheable endpoint
	addCacheHandler(r, "GET", "/posts/:num", homeHandler, &cache, database)

	// Where all the static files (css, js, etc) are served from

	r.Static("/images/data", settings.ImageDirectory)
	r.Static("/static", "./static")
	r.StaticFS("/media", http.Dir(settings.ImageDirectory))
	permalinks, err := database.GetPermalinks()
	if err != nil {
		log.Error().Msgf("could not get permalinks: %v", err)
	} else {
		for _, permalink := range permalinks {
			addCacheHandler(r, "GET", permalink.Path, permalinkPostHandler(permalink.PostId), &cache, database)
		}
	}

	r.NoRoute(notFoundHandler())

	return r
}

func addCacheHandler(e *gin.Engine, method string, endpoint string, generator Generator, cache *Cache, db database.Database) {

	handler := func(c *gin.Context) {
		// if the endpoint is cached
		if common.Settings.CacheEnabled {
			cached_endpoint, err := (*cache).Get(c.Request.RequestURI)
			if err == nil {
				log.Info().Msgf("cache hit for page: %s", c.Request.RequestURI)
				c.Data(http.StatusOK, "text/html; charset=utf-8", cached_endpoint.Contents)
				return
			}
		}

		// Before handler call (retrieve from cache)
		html_buffer, err := generator(c, db)
		if err != nil {
			log.Error().Msgf("could not generate html: %v", err)
			// TODO : Need a proper error page
			c.JSON(http.StatusInternalServerError, common.ErrorRes("could not render HTML", err))
			return
		}

		// After handler  (add to cache)
		if common.Settings.CacheEnabled {
			err = (*cache).Store(c.Request.RequestURI, html_buffer)
			if err != nil {
				log.Warn().Msgf("could not add page to cache: %v", err)
			}
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
	pageNum := 0 // Default to page 0
	if pageNumQuery := c.Param("num"); pageNumQuery != "" {
		num, err := strconv.Atoi(pageNumQuery)
		if err == nil && num > 0 {
			pageNum = num
		} else {
			log.Error().Msgf("Invalid page number: %s", pageNumQuery)
		}
	}
	limit := 10 // or whatever limit you want
	offset := max((pageNum-1)*limit, 0)

	posts, err := db.GetPosts(limit, offset)
	if err != nil {
		log.Error().Msgf("Failed to load posts: %v", err)
		return []byte("error: Failed to load posts"), err
	}
	sticky_posts := make([]common.Post, 0)
	for _, sticky_post_id := range common.Settings.StickyPosts {
		post, err := db.GetPost(sticky_post_id)
		if err != nil {
			log.Error().Msgf("could not find sticky post `%d`: %v", sticky_post_id, err)
			continue
		}
		post.Content = string(mdToHTML([]byte(post.Content)))
		sticky_posts = append(sticky_posts, post)
	}

	// if not cached, create the cache
	index_view := views.MakeIndex(posts, sticky_posts, common.Settings.AppNavbar.Links, common.Settings.AppNavbar.Dropdowns)
	// if not cached, create the cache
	html_buffer := bytes.NewBuffer(nil)
	err = index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil

}

func notFoundHandler() func(*gin.Context) {
	handler := func(c *gin.Context) {
		buffer, err := renderHtml(c, views.MakeNotFoundPage(common.Settings.AppNavbar.Links, common.Settings.AppNavbar.Dropdowns))
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrorRes("could not render HTML", err))
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", buffer)
	}

	return handler
}
