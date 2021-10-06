package dto

import (
	"encoding/json"
	"io"
)

type AccountDto struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func DecodeAccount(r io.Reader) (*AccountDto, error) {
	var account AccountDto
	if err := json.NewDecoder(r).Decode(&account); err != nil {
		return nil, err
	} else {
		return &account, err
	}
}
