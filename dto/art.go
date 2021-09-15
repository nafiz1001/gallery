package dto

import (
	"encoding/json"
	"io"
)

type ArtDto struct {
	Id       int    `json:"id"`
	Quantity int    `json:"quantity"`
	Title    string `json:"title"`
	AuthorId int    `json:"author_id"`
}

func DecodeArt(r io.Reader) (*ArtDto, error) {
	var art ArtDto
	err := json.NewDecoder(r).Decode(&art)
	return &art, err
}
