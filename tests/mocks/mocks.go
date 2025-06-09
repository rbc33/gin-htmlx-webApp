package mocks

import (
	"fmt"

	"github.com/rbc33/gocms/common"
)

type DatabaseMock struct {
	GetPostHandler  func(int) (common.Post, error)
	GetPostsHandler func() ([]common.Post, error)
}

func (db DatabaseMock) GetPosts() ([]common.Post, error) {
	return db.GetPostsHandler()
}

func (db DatabaseMock) GetPost(post_id int) (common.Post, error) {
	return db.GetPostHandler(post_id)
}

func (db DatabaseMock) AddPost(title string, excerpt string, content string) (int, error) {
	return 0, nil
}

func (db DatabaseMock) ChangePost(id int, title string, excerpt string, content string) error {
	return nil
}

func (db DatabaseMock) DeletePost(id int) error {
	return fmt.Errorf("not implemented")
}

func (db DatabaseMock) AddImage(string, string, string) error {
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
