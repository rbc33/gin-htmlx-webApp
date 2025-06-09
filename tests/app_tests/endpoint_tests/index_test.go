package endpoint_tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbc33/gocms/app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestIndexSuccess(t *testing.T) {

	database_mock := mocks.DatabaseMock{
		GetPostsHandler: func() ([]common.Post, error) {
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

	r := app.SetupRoutes(database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/page/0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
	assert.Contains(t, w.Body.String(), "TestExcerpt")
}

func TestIndexFailToGetPosts(t *testing.T) {

	database_mock := mocks.DatabaseMock{
		GetPostsHandler: func() ([]common.Post, error) {
			return []common.Post{}, fmt.Errorf("invalid")
		},
	}

	r := app.SetupRoutes(database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
