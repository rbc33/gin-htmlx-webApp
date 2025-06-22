package admin_app

import "github.com/rbc33/gocms/common"

// Extracted all bindings and requests structs into a single package to
// organize the data in a simpler way. Every domain object supporting
// CRUD endpoints has their own structures to handle the http methods.

type AddPageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

type DeletePostBinding struct {
	common.IntIdBinding
}
type DeletePageBinding struct {
	Link string `json:"link"`
}

type AddImageRequest struct {
	Alt string `json:"alt"`
}

type DeleteImageBinding struct {
	Name string `uri:"name" binding:"required"`
}
type AddPostRequest struct {
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type ChangePostRequest struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type ChangePageRequest struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Content string `json:"content"`
}
