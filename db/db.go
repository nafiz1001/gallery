package db

import (
	"fmt"

	"github.com/nafiz1001/gallery-go/arts"
)

type DB struct {
	arts []arts.Art
}

func Init() DB {
	db := DB{arts: []arts.Art{}}
	return db
}

func (db *DB) StoreArt(a arts.Art) (*arts.Art, error) {
	a.Id = len(db.arts)
	db.arts = append(db.arts, a)
	return &a, nil
}

func (db *DB) RetrieveArt(id int) (*arts.Art, error) {
	for i := range db.arts {
		if db.arts[i].Id == id {
			return &db.arts[i], nil
		}
	}

	return nil, fmt.Errorf("could not find art with id %d", id)
}

func (db *DB) RetrieveAllArt() ([]arts.Art, error) {
	arts := make([]arts.Art, len(db.arts))
	copy(arts, db.arts)
	return arts, nil
}

func (db *DB) UpdateArt(a arts.Art) (*arts.Art, error) {
	art, err := db.RetrieveArt(a.Id)

	if art != nil {
		art.Genres = a.Genres
		art.Picture = a.Picture
		art.Price = a.Price
		art.Quantity = a.Quantity
		art.Title = a.Title
	}

	return art, err
}

func (db *DB) DeleteArt(id int) (*arts.Art, error) {
	var art *arts.Art
	index := -1

	for i := range db.arts {
		if db.arts[i].Id == id {
			art = &db.arts[i]
			index = i
			break
		}
	}

	if index < 0 {
		return nil, fmt.Errorf("could not find art with id %d", id)
	} else {
		db.arts = append(db.arts[:index], db.arts[index+1:]...)
		return art, nil
	}
}
