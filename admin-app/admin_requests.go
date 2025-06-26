package admin_app

import (
	"encoding/json"

	"github.com/rbc33/gocms/common"
)

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
	Excerpt string `json:"excerpt"`
}

type DeleteImageBinding struct {
	Name string `uri:"name" binding:"required"`
}
type DeleteSchemaBinding struct {
	Id string `uri:"Uuid" binding:"required"`
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

type AddCardRequest struct {
	Image   string `json:"image_location"`
	Schema  string `json:"schema"`
	Content string `json:"data"`
}

type GetCardRequest struct {
	Schema string `uri:"schema" binding:"required"`
	Limit  uint32 `uri:"limit"`
	Page   uint32 `uri:"page"`
}

type AddCardSchemaRequest struct {
	JsonTitle  string `json:"title"`
	JsonSchema string `json:"schema"`
}

// UnmarshalJSON is a custom unmarshaller for Content
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
