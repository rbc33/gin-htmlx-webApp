package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbc33/gocms/app"
	"github.com/rbc33/gocms/common"
	"github.com/test-go/testify/assert"
)

type DatabaseMock struct{}

func (db DatabaseMock) GetPosts(offset int, limit int) ([]common.Post, error) {
	return []common.Post{
		{
			Title:   "TestPost",
			Content: "TestContent",
			Excerpt: "TestExcerpt",
			Id:      0,
		},
	}, nil
}
func (db DatabaseMock) GetPost(post_id int) (common.Post, error) {
	return common.Post{}, fmt.Errorf("not implemented")
}
func (db DatabaseMock) AddPost(title string, excerpt string, content string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
func (db DatabaseMock) ChangePost(id int, title string, excerpt string, content string) error {
	return nil
}
func (db DatabaseMock) DeletePost(id int) error {
	return fmt.Errorf("not implemented")
}

func (db DatabaseMock) AddImage(postID string, imageData string, imageType string) error {
	return fmt.Errorf("not implemented")
}
func (db DatabaseMock) GetCard(uuid string) (common.Card, error) {
	return common.Card{}, fmt.Errorf("not implemented")
}
func (db DatabaseMock) AddCard(uuid string, image_location string, json_data string, schema_data string) error {
	return fmt.Errorf("not implemented")
}
func (db DatabaseMock) ChangeCard(uuid string, image_location string, json_data string, schema_data string) error {
	return fmt.Errorf("not implemented")
}
func (db DatabaseMock) DeleteCard(uuid string) error {
	return fmt.Errorf("not implemented")
}
func (db DatabaseMock) DeleteImage(uuid string) error {
	return fmt.Errorf("not implemented")
}
func (db DatabaseMock) AddPage(title string, content string, link string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
func (db DatabaseMock) GetPage(link string) (common.Page, error) {
	return common.Page{}, fmt.Errorf("not implemented")
}

func TestIndexPing(t *testing.T) {

	database_mock := DatabaseMock{}
	r := app.SetupRoutes(common.Settings, &database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
	assert.Contains(t, w.Body.String(), "TestExcerpt")
}
