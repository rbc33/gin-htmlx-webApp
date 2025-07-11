package admin_app

import (
	"encoding/json"

	"github.com/rbc33/gocms/common"
)

// swagger:ignore Extracted all bindings and requests structs into a single package to
// swagger:ignore organize the data in a simpler way. Every domain object supporting
// swagger:ignore CRUD endpoints has their own structures to handle the http methods.

// swagger:parameters addPageRequest AddPageRequest
type AddPageRequest struct {
	// Title of the page
	// in: body
	// required: true
	Title string `json:"title"`
	// Content of the page
	// in: body
	Content string `json:"content"`
	// Link of the page
	// in: body
	Link string `json:"link"`
}

// swagger:parameters deletePostRequest DeletePostRequest
type DeletePostRequest struct {
	common.IntIdBinding
}

// swagger:parameters deletePageRequest DeletePageRequest
type DeletePageRequest struct {
	// Link of the page to delete
	// in: body
	// required: true
	Link string `json:"link"`
}

// swagger:parameters addImageRequest AddImageRequest
type AddImageRequest struct {
	// Excerpt for the image
	// in: body
	Excerpt string `json:"excerpt"`
}

// swagger:parameters deleteImageRequest DeleteImageRequest
type DeleteImageRequest struct {
	// Name of the image to delete
	// in: path
	// required: true
	Name string `uri:"name" binding:"required"`
}

// swagger:parameters deleteSchemaRequest DeleteSchemaBinding
type DeleteSchemaBinding struct {
	// UUID of the schema to delete
	// in: path
	// required: true
	Id string `uri:"Uuid" binding:"required"`
}

// swagger:parameters addPostRequest AddPostRequest
type AddPostRequest struct {
	// Title of the post
	// in: body
	// required: true
	Title string `json:"title"`
	// Excerpt of the post
	// in: body
	Excerpt string `json:"excerpt"`
	// Content of the post
	// in: body
	Content string `json:"content"`
}

// swagger:parameters changePostRequest ChangePostRequest
type ChangePostRequest struct {
	// ID of the post
	// in: body
	// required: true
	Id int `json:"id"`
	// Title of the post
	// in: body
	Title string `json:"title"`
	// Excerpt of the post
	// in: body
	Excerpt string `json:"excerpt"`
	// Content of the post
	// in: body
	Content string `json:"content"`
}

// swagger:parameters changePageRequest ChangePageRequest
type ChangePageRequest struct {
	// ID of the page
	// in: body
	// required: true
	Id int `json:"id"`
	// Title of the page
	// in: body
	Title string `json:"title"`
	// Link of the page
	// in: body
	Link string `json:"link"`
	// Content of the page
	// in: body
	Content string `json:"content"`
}

// swagger:parameters addCardRequest AddCardRequest
type AddCardRequest struct {
	// Image location URL
	// in: body
	Image string `json:"image_location"`
	// Schema name
	// in: body
	Schema string `json:"schema"`
	// Content data of the card
	// in: body
	Content string `json:"data"`
}

// swagger:parameters changeCardRequest ChangeCardRequest
type ChangeCardRequest struct {
	// ID of the card
	// in: body
	Id string `json:"id"`
	// Image location URL
	// in: body
	ImageLocation string `json:"image_location"`
	// JSON data of the card
	// in: body
	JsonData string `json:"json_data"`
	// JSON schema name
	// in: body
	SchemaName string `json:"json_schema"`
}

// swagger:parameters deleteCardRequest DeleteCardRequest
type DeleteCardRequest struct {
	// ID of the card to delete
	// in: body
	Id string `json:"id"`
}

// swagger:parameters getCardRequest GetCardRequest
type GetCardRequest struct {
	// Schema name to filter cards
	// in: path
	// required: true
	Schema string `uri:"schema" binding:"required"`
	// Limit number of cards to return
	// in: query
	Limit uint32 `uri:"limit"`
	// Page number for pagination
	// in: query
	Page uint32 `uri:"page"`
}

// swagger:parameters addCardSchemaRequest AddCardSchemaRequest
type AddCardSchemaRequest struct {
	// Title of the card schema
	// in: body
	// required: true
	JsonTitle string `json:"title"`
	// JSON schema string
	// in: body
	// required: true
	JsonSchema string `json:"schema"`
}

// swagger:parameters addPermalinkRequest AddPermalinkRequest
type AddPermalinkRequest struct {
	// Permalink string
	// in: uri
	// required: true
	Permalink string `uri:"permalink" binding:"required"`
	// Post ID associated with the permalink
	// in: uri
	// required: true
	PostId int `uri:"post_id" binding:"required"`
}

// UnmarshalJSON is a custom unmarshaller for AddCardSchemaRequest that preserves the raw JSON schema string.
func (c *AddCardSchemaRequest) UnmarshalJSON(data []byte) error {

	// Create a map to hold the raw JSON
	var obj_map map[string]*json.RawMessage
	err := json.Unmarshal(data, &obj_map)
	if err != nil {
		return err
	}

	// Extract title as normal
	if title_bytes, ok := obj_map["title"]; ok && title_bytes != nil {
		var title string
		if err := json.Unmarshal(*title_bytes, &title); err != nil {
			return err
		}
		c.JsonTitle = title
	}

	// Extract schema as a raw string
	if schema_bytes, ok := obj_map["schema"]; ok && schema_bytes != nil {
		// Convert the raw schema to a string, preserving its JSON structure
		c.JsonSchema = string(*schema_bytes)
	}

	return nil
}
