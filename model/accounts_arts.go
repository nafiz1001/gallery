package model

import (
	"database/sql"

	"github.com/nafiz1001/gallery-go/dto"
)

type AccountArtsDB struct {
	accountDB *AccountDB
	artDB     *ArtDB
	sqlDB     *sql.DB
}

func (db *AccountArtsDB) Init(sqlDB *sql.DB, accountDB *AccountDB, artDB *ArtDB) error {
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

func (db *AccountArtsDB) AddArt(account dto.AccountDto, art dto.ArtDto) (*dto.ArtDto, error) {
	if acc, err := db.accountDB.GetAccountById(account.Id); err != nil {
		return nil, err
	} else {
		if a, err := db.artDB.CreateArt(art); err != nil {
			return nil, err
		} else {
			if _, err := db.sqlDB.Exec("INSERT INTO account_arts (art_id, account_id) VALUES (?, ?)", a.Id, acc.Id); err != nil {
				return nil, err
			} else {
				a.AuthorId = acc.Id
				return a, nil
			}
		}
	}
}

func (db *AccountArtsDB) IsAuthor(account dto.AccountDto, artId int) bool {
	if acc, err := db.accountDB.GetAccountByUsername(account.Username); err != nil {
		return false
	} else {
		var tmpId int
		err := db.sqlDB.QueryRow("SELECT id FROM account_arts WHERE art_id = ? AND account_id = ?", artId, acc.Id).Scan(&tmpId)
		return err == nil
	}
}

func (db *AccountArtsDB) GetArts() ([]dto.ArtDto, error) {
	v := []dto.ArtDto{}

	if rows, err := db.sqlDB.Query("SELECT art_id, account_id FROM account_arts"); err != nil {
		return []dto.ArtDto{}, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var artId int
			var accountId int
			if err := rows.Scan(&artId, &accountId); err != nil {
				return []dto.ArtDto{}, err
			} else {
				if art, err := db.artDB.GetArt(artId); err != nil {
					return []dto.ArtDto{}, err
				} else {
					art.AuthorId = accountId
					v = append(v, *art)
				}
			}
		}
	}

	return v, nil
}

func (db *AccountArtsDB) GetArtsByAccountId(accountId int) ([]dto.ArtDto, error) {
	v := []dto.ArtDto{}

	if rows, err := db.sqlDB.Query("SELECT art_id FROM account_arts WHERE account_id = ?", accountId); err != nil {
		return []dto.ArtDto{}, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var artId int
			if err := rows.Scan(&artId, &accountId); err != nil {
				return []dto.ArtDto{}, err
			} else {
				if art, err := db.artDB.GetArt(artId); err != nil {
					return []dto.ArtDto{}, err
				} else {
					art.AuthorId = accountId
					v = append(v, *art)
				}
			}
		}
	}

	return v, nil
}

func (db *AccountArtsDB) DeleteArt(id int) (*dto.ArtDto, error) {
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
