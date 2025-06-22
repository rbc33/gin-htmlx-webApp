package app_system_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbc33/gocms/app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestIndexPing(t *testing.T) {

	database_mock := mocks.DatabaseMock{
		GetPostsHandler: func(offset int, limit int) ([]common.Post, error) {
			return []common.Post{
				{
					Title:   "TestPost",
					Content: "TestContent",
					Excerpt: "TestExcerpt",
					Id:      0,
				},
			}, nil
		},
	}
	r := app.SetupRoutes(common.Settings, &database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
	assert.Contains(t, w.Body.String(), "TestExcerpt")
}
