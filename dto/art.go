package dto

import (
	"encoding/json"
	"io"
)

type ArtDto struct {
	Id       uint   `json:"id"`
	Quantity int    `json:"quantity"`
	Title    string `json:"title"`
	AuthorId uint   `json:"author_id"`
}

func DecodeArt(r io.Reader) (*ArtDto, error) {
	var art ArtDto
	if err := json.NewDecoder(r).Decode(&art); err != nil {
		return nil, err
	} else {
		return &art, err
	}
}
