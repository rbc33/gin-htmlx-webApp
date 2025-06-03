package app

import (
	"bytes"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/database"
	"github.com/rbc33/views/tailwind"
	"github.com/rs/zerolog/log"
)

func makeContactFormHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		if err != nil {
			log.Error().Msgf("%s", err)
		}

		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")

		// Parse email
		_, err = mail.ParseAddress(email)
		if err != nil {
			// TemplRender(c, http.StatusOK, views.MakeContactFailure(email, err.Error()))
			err = TemplRender(c, http.StatusOK, tailwind.MakeContactFailure(email, err.Error()))
			if err != nil {
				log.Error().Msgf("%s", err)
			}
			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			// TemplRender(c, http.StatusOK, views.MakeContactFailure(email, "name too long (200 char max)"))
			err = TemplRender(c, http.StatusOK, tailwind.MakeContactFailure(email, "name too long (200 char max)"))
			if err != nil {
				log.Error().Msgf("%s", err)
			}
			return
		}

		if len(message) > 10000 {
			// TemplRender(c, http.StatusOK, views.MakeContactFailure(email, "message too long (10000 char max)"))
			err = TemplRender(c, http.StatusOK, tailwind.MakeContactFailure(email, "message too long (10000 char max)"))
			if err != nil {
				log.Error().Msgf("%s", err)
			}

			return
		}

		// TemplRender(c, http.StatusOK, views.MakeContactSuccess(email, name))
		err = TemplRender(c, http.StatusOK, tailwind.MakeContactSuccess(email, name))
		if err != nil {
			log.Error().Msgf("%s", err)
		}

	}
}

// TODO : This is a duplicate of the index handler... abstract

func contactHandler(c *gin.Context, db *database.Database) ([]byte, error) {
	// index_view := views.MakeContactPage()
	index_view := tailwind.MakeContactPage()
	html_buffer := bytes.NewBuffer(nil)
	err := index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return html_buffer.Bytes(), nil
}
