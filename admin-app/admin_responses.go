package admin_app

import "github.com/rbc33/gocms/common"

// swagger:response PageResponse
type PageResponse struct {
	// ID of the page
	Id int `json:"id"`
	// Link of the page
	Link string `json:"link"`
}

// swagger:response GetPagesResponse
type GetPagesResponse struct {
	// List of pages
	Pages []common.Page `json:"pages"`
}

// swagger:response PostIdResponse
type PostIdResponse struct {
	// ID of the post
	Id int `json:"id"`
}

// swagger:response PermalinkIdResponse
type PermalinkIdResponse struct {
	// ID of the post
	PostId int `json:"post_id"`
}

// swagger:response GetPostsResponse
type GetPostsResponse struct {
	// List of posts
	Posts []common.Post `json:"posts"`
}

// swagger:response GetPostResponse
type GetPostResponse struct {
	// ID of the post
	Id int `json:"id"`
	// Title of the post
	Title string `json:"title"`
	// Excerpt of the post
	Excerpt string `json:"excerpt"`
	// Content of the post
	Content string `json:"content"`
}

// swagger:response ImageIdResponse
type ImageIdResponse struct {
	// ID of the image
	Id string `json:"id"`
}

// swagger:response GetImageResponse
type GetImageResponse struct {
	// UUID of the image
	Id string `json:"uuid"`
	// Name of the image
	Name string `json:"name"`
	// Alternative text for the image
	AltText string `json:"alt_text"`
	// File extension of the image
	Extension string `json:"extension"`
}

// swagger:response CardIdResponse
type CardIdResponse struct {
	// ID of the card
	Id string `json:"id"`
}

// swagger:response CardSchemaResponse
type CardSchemaResponse struct {
	// UUID of the card schema
	Id string `json:"uuid"`
}

// swagger:response GetSchemaasResponse
type GetSchemaasResponse struct {
	// List of card schemas
	Schemas []common.CardSchema `json:"schemas"`
}

// swagger:response DeletePageResponse
type DeletePageResponse struct {
	// UUID of the deleted page
	Id string `json:"uuid"`
}
