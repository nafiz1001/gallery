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

func (db *ArtDB) CreateArt(art dto.ArtDto) (*dto.ArtDto, error) {
	art.Id = util.CreateId()
	db.arts[art.Id] = &art
	return &art, nil
}

func (db *ArtDB) GetArt(id string) (*dto.ArtDto, error) {

	if art, ok := db.arts[id]; !ok {
		return nil, fmt.Errorf("could not find art with id %s", id)
	} else {
		return art, nil
	}
}

func (db *ArtDB) GetArts() ([]dto.ArtDto, error) {
	arts := []dto.ArtDto{}

	for id := range db.arts {
		arts = append(arts, *db.arts[id])
	}

	return arts, nil
}

func (db *ArtDB) UpdateArt(art dto.ArtDto) (*dto.ArtDto, error) {
	if a, ok := db.arts[art.Id]; !ok {
		return nil, fmt.Errorf("could not find art with id %s", art.Id)
	} else {
		a.Quantity = art.Quantity
		a.Title = art.Title
		return db.arts[a.Id], nil
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
