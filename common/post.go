package common

type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}
type Card struct {
	Uuid          string `json:"uuid"`
	ImageLocation string `json:"image_location"`
	JsonData      string `json:"json_data"`
	SchemaName    string `json:"json_chema"`
}
