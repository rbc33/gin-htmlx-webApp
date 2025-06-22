package mocks

import (
	"fmt"

	"github.com/rbc33/gocms/common"
)

type DatabaseMock struct {
	GetPostHandler  func(int) (common.Post, error)
	GetPostsHandler func(int, int) ([]common.Post, error)
	AddPageHandler  func(string, string, string) (int, error)
	GetPagesHandler func(int, int) ([]common.Page, error)
}

func (db DatabaseMock) GetPosts(offset int, limit int) ([]common.Post, error) {
	if db.GetPostsHandler != nil {
		return db.GetPostsHandler(offset, limit)
	}
	return nil, fmt.Errorf("GetPostsHandler not set")
}
func (db DatabaseMock) GetPages(offset int, limit int) ([]common.Page, error) {
	if db.GetPostsHandler != nil {
		return db.GetPagesHandler(offset, limit)
	}
	return nil, fmt.Errorf("GetPageHandler not set")
}
func (db DatabaseMock) GetPost(post_id int) (common.Post, error) {
	if db.GetPostHandler != nil {
		return db.GetPostHandler(post_id)
	}
	return common.Post{}, fmt.Errorf("not implemented")
}
func (db DatabaseMock) AddPost(title string, excerpt string, content string) (int, error) {
	// Simulate successful post addition with a positive ID.
	// This helps TestCreatePost_Success satisfy the ID check.
	return 0, nil
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
	return db.AddPageHandler(title, content, link)
}
func (db DatabaseMock) GetPage(link string) (common.Page, error) {
	return common.Page{}, fmt.Errorf("not implemented")
}
func (db DatabaseMock) ChangePage(id int, title string, content string, link string) (err error) {
	return fmt.Errorf("not implemented")
}
func (db DatabaseMock) DeletePage(link string) error {
	return fmt.Errorf("not implemented")
}
