package model

import (
	"database/sql"

	"github.com/nafiz1001/gallery-go/dto"
)

type ArtDB struct {
	sqlDB *sql.DB
}

func (db *ArtDB) Init(sqlDB *sql.DB) error {
	db.sqlDB = sqlDB

	_, err := sqlDB.Exec(
		`CREATE TABLE arts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			quantity INTEGER NOT NULL
		  );`,
	)

	return err
}

func (db *ArtDB) CreateArt(art dto.ArtDto) (*dto.ArtDto, error) {
	if res, err := db.sqlDB.Exec(`INSERT INTO arts (title, quantity) VALUES (?, ?)`, art.Title, art.Quantity); err != nil {
		return nil, err
	} else {
		id, _ := res.LastInsertId()
		art.Id = int(id)
		return &art, nil
	}
}

func (db *ArtDB) GetArt(id int) (*dto.ArtDto, error) {
	var title string
	var quantity int
	if err := db.sqlDB.QueryRow(`SELECT id, title, quantity FROM arts WHERE id = ?`, id).Scan(&id, &title, &quantity); err != nil {
		return nil, err
	} else {
		return &dto.ArtDto{
			Id:       id,
			Title:    title,
			Quantity: quantity,
		}, nil
	}
}

func (db *ArtDB) GetArts() ([]dto.ArtDto, error) {
	arts := []dto.ArtDto{}

	if rows, err := db.sqlDB.Query(`SELECT id, title, quantity FROM arts`); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int
			var title string
			var quantity int

			if err := rows.Scan(&id, &title, &quantity); err != nil {
				return nil, err
			} else {
				arts = append(arts, dto.ArtDto{
					Id:       id,
					Title:    title,
					Quantity: quantity,
				})
			}
		}
	}

	return arts, nil
}

func (db *ArtDB) UpdateArt(art dto.ArtDto) (*dto.ArtDto, error) {
	if _, err := db.sqlDB.Exec(`UPDATE arts SET title=?, quantity=? WHERE id = ?`, art.Title, art.Quantity, art.Id); err != nil {
		return nil, err
	} else {
		return &art, nil
	}
}

func (db *ArtDB) DeleteArt(id int) (*dto.ArtDto, error) {
	if art, err := db.GetArt(id); err != nil {
		return nil, err
	} else {
		if _, err := db.sqlDB.Exec(`DELETE FROM arts WHERE id = ?`, art.Id); err != nil {
			return nil, err
		} else {
			return art, nil
		}
	}
}
