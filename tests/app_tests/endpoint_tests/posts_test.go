package endpoint_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbc33/gocms/app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPostSuccess(t *testing.T) {

	database_mock := mocks.DatabaseMock{
		GetPostHandler: func(post_id int) (common.Post, error) {
			return common.Post{
				Title:   "TestPost",
				Content: "TestContent",
				Excerpt: "TestExcerpt",
				Id:      post_id,
			}, nil
		},
	}

	r := app.SetupRoutes(database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
}
