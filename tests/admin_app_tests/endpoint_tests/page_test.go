package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	admin_app "github.com/rbc33/gocms/admin-app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/plugins"
	"github.com/rbc33/gocms/tests/mocks"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestAddPageHappyPath(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPageHandler: func(string, string, string) (int, error) {
			return 0, nil
		},
	}

	page_data := admin_app.AddPageRequest{
		Title:   "Title",
		Content: "Content",
		Link:    "Link",
	}

	shortcode_handlers, err := admin_app.LoadShortcodesHandlers(common.Settings.Shortcodes)
	if err != nil {
		log.Error().Msgf("%s", err)
		os.Exit(-1)
	}

	post_hook := &plugins.PostHook{}
	image_plugin := plugins.Plugin{
		ScriptName: "img",
		Id:         "img-plugin",
	}
	post_hook.Register(image_plugin)
	// img, _ := shortcode_handlers["img"]
	hooks_map := map[string]plugins.Hook{
		"add_post": post_hook,
	}
	router := admin_app.SetupRoutes(app_settings, shortcode_handlers, databaseMock, hooks_map)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(page_data)
	request, _ := http.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	var response admin_app.PageResponse
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.NotNil(t, response.Id)
	assert.NotEmpty(t, response.Link)
	assert.Equal(t, page_data.Link, response.Link)
}
