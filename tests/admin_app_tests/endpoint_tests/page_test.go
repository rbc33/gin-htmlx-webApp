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

type pageRequest struct {
	Title   string `json:"title"`
	Link    string `json:"excerpt"`
	Content string `json:"content"`
}

type pageResponse struct {
	Id   int `json:"id"`
	Link int `json:"link"`
}

var app_settings = test.GetAppSettings()

func TestCreatePage_Success(t *testing.T) {

	shortcode_handlers, err := admin_app.LoadShortcodesHandlers(common.Settings.Shortcodes)
	if err != nil {
		log.Error().Msgf("%s", err)
		os.Exit(-1)
	}

	database_mock := mocks.DatabaseMock{}
	r := admin_app.SetupRoutes(app_settings, shortcode_handlers, database_mock)
	w := httptest.NewRecorder()

	request := pageRequest{
		Title:   "Never gonna",
		Content: "<p>give</p",
		Link:    "you-app",
	}
	request_body, err := json.Marshal(request)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(request_body))
	req.Header.Add("content-type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code) // Expect 201 Created for successful creation

	var response pageResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Id, 0) // Expect a positive ID for the new post
}
