package dto

type ArtDto struct {
	Id       int    `json:"id"`
	Quantity int    `json:"quantity"`
	Title    string `json:"title"`
	AuthorId int    `json:"author_id"`
}
