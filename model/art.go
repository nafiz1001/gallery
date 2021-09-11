package model

import (
	"fmt"

	"github.com/nafiz1001/gallery-go/dto"
	"github.com/nafiz1001/gallery-go/util"
)

type ArtDB struct {
	arts map[string]*dto.ArtDto
}

func (db *ArtDB) Init() error {
	db.arts = map[string]*dto.ArtDto{}
	return nil
}

func (db *ArtDB) StoreArt(a dto.ArtDto) (*dto.ArtDto, error) {
	a.Id = util.CreateId()
	db.arts[a.Id] = &a
	return &a, nil
}

func (db *ArtDB) RetrieveArt(id string) (*dto.ArtDto, error) {

	if art, ok := db.arts[id]; ok {
		return art, nil
	} else {
		return nil, fmt.Errorf("could not find art with id %s", id)
	}
}

func (db *ArtDB) RetrieveArts() ([]dto.ArtDto, error) {
	arts := []dto.ArtDto{}

	for id := range db.arts {
		arts = append(arts, *db.arts[id])
	}

	return arts, nil
}

func (db *ArtDB) UpdateArt(a dto.ArtDto) (*dto.ArtDto, error) {
	if art, ok := db.arts[a.Id]; ok {
		art.Quantity = a.Quantity
		art.Title = a.Title
		return db.arts[a.Id], nil
	} else {
		return nil, fmt.Errorf("could not find art with id %s", a.Id)
	}
}

func (db *ArtDB) DeleteArt(id string) (*dto.ArtDto, error) {
	if art, ok := db.arts[id]; !ok {
		return nil, fmt.Errorf("could not find art with id %s", id)
	} else {
		delete(db.arts, art.Id)
		return art, nil
	}
}
