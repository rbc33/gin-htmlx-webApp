package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	admin_app "github.com/rbc33/gocms/admin-app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/tests/mocks"
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

var app_settings = common.AppSettings{
	DatabaseUri:   "root:secret@tcp(192.168.0.100:33060)/gocms",
	WebserverPort: "8080",
}

func TestIndexPing(t *testing.T) {

	database_mock := mocks.DatabaseMock{}
	r := admin_app.SetupRoutes(&database_mock)
	w := httptest.NewRecorder()

	request := postRequest{
		Title:   "",
		Excerpt: "",
		Content: "",
	}
	request_body, err := json.Marshal(request)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(request_body))
	req.Header.Add("content-type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response postResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Id, 0)
}

// TODO : Test request without excerpt

// TODO : Test request without content

// TODO : Test request without title

// TODO : Test request that fails to be added to database
