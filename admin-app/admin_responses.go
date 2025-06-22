package admin_app

import "github.com/rbc33/gocms/common"

type PageResponse struct {
	Id   int    `json:"id"`
	Link string `json:"link"`
}
type GetPagesResponse struct {
	Pages []common.Page `json:"pages"`
}

type PostIdResponse struct {
	Id int `json:"id"`
}

type GetPostsResponse struct {
	Posts []common.Post `json:"posts"`
}
type GetPostResponse struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type ImageIdResponse struct {
	Id string `json:"id"`
}

type GetImageResponse struct {
	Id        string `json:"uuid"`
	Name      string `json:"name"`
	AltText   string `json:"alt_text"`
	Extension string `json:"extension"`
}
