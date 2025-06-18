package endpoint_tests

import (
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	admin_app "github.com/rbc33/gocms/admin-app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/tests/mocks"
	test "github.com/rbc33/gocms/tests/system_tests/helpers"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type postRequest struct {
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type postResponse struct {
	Id int `json:"id"`
}

var app_settings = test.GetAppSettings()

func TestCreatePost_Success(t *testing.T) {

	shortcode_handlers, err := admin_app.LoadShortcodesHandlers(common.Settings.Shortcodes)
	if err != nil {
		log.Error().Msgf("%s", err)
		os.Exit(-1)
	}

	database_mock := mocks.DatabaseMock{}
	r := admin_app.SetupRoutes(app_settings, shortcode_handlers, database_mock)
	w := httptest.NewRecorder()

	request := postRequest{
		Title:   "New Test Post",
		Excerpt: "A brief summary of the post.",
		Content: "This is the full content of the new test post.",
	}
	request_body, err := json.Marshal(request)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(request_body))
	req.Header.Add("content-type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code) // Expect 201 Created for successful creation

	var response postResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Id, 0) // Expect a positive ID for the new post
}
