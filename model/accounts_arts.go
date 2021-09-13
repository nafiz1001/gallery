package model

import (
	"database/sql"

	"github.com/nafiz1001/gallery-go/dto"
)

type AccountsArtsDB struct {
	accountDB *AccountDB
	artDB     *ArtDB
	sqlDB     *sql.DB
}

func (db *AccountsArtsDB) Init(sqlDB *sql.DB, accountDB *AccountDB, artDB *ArtDB) error {
	db.accountDB = accountDB
	db.artDB = artDB
	db.sqlDB = sqlDB

	_, err := sqlDB.Exec(
		`CREATE TABLE account_arts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			art_id INTEGER NOT NULL,
			account_id INTEGER NOT NULL,
			FOREIGN KEY (art_id) REFERENCES arts (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
			FOREIGN KEY (account_id) REFERENCES accounts (id) ON UPDATE RESTRICT ON DELETE RESTRICT
		  );`,
	)

	return err
}

func (db *AccountsArtsDB) AddArt(account dto.AccountDto, art dto.ArtDto) (*dto.ArtDto, error) {
	if acc, err := db.accountDB.GetAccount(account.Username); err != nil {
		return nil, err
	} else {
		if a, err := db.artDB.CreateArt(art); err != nil {
			return nil, err
		} else {
			if _, err := db.sqlDB.Exec("INSERT INTO account_arts (art_id, account_id) VALUES (?, ?)", a.Id, acc.Id); err != nil {
				return nil, err
			} else {
				return a, nil
			}
		}
	}
}

func (db *AccountsArtsDB) IsAuthor(account dto.AccountDto, artId int) bool {
	if acc, err := db.accountDB.GetAccount(account.Username); err != nil {
		return false
	} else {
		var tmpId int
		err := db.sqlDB.QueryRow("SELECT id FROM account_arts WHERE art_id = ? AND account_id = ?", artId, acc.Id).Scan(&tmpId)
		return err == nil
	}
}

func (db *AccountsArtsDB) GetArtsByUsername(username string) ([]dto.ArtDto, error) {
	v := []dto.ArtDto{}

	if account, err := db.accountDB.GetAccount(username); err != nil {
		return v, err
	} else {
		if rows, err := db.sqlDB.Query("SELECT art_id FROM account_arts WHERE account_id = ?", account.Id); err != nil {
			return []dto.ArtDto{}, err
		} else {
			defer rows.Close()
			for rows.Next() {
				var artId int
				if err := rows.Scan(&artId); err != nil {
					return []dto.ArtDto{}, err
				} else {
					if art, err := db.artDB.GetArt(artId); err != nil {
						return []dto.ArtDto{}, err
					} else {
						v = append(v, *art)
					}
				}
			}
		}
	}

	return v, nil
}

func (db *AccountsArtsDB) DeleteArt(id int) (*dto.ArtDto, error) {
	if art, err := db.artDB.GetArt(id); err != nil {
		return nil, err
	} else {
		if _, err := db.sqlDB.Exec("DELETE FROM account_arts WHERE art_id = ?", art.Id); err != nil {
			return nil, err
		} else {
			if art, err := db.artDB.DeleteArt(id); err != nil {
				return nil, err
			} else {
				return art, nil
			}
		}
	}
}
