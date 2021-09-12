package model

import (
	"database/sql"
	"fmt"

	"github.com/nafiz1001/gallery-go/dto"
)

type AccountDB struct {
	sqlDB *sql.DB
}

func (db *AccountDB) Init(sqlDB *sql.DB) error {
	db.sqlDB = sqlDB

	_, err := sqlDB.Exec(
		`CREATE TABLE accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL
		  );`,
	)

	return err
}

func (db *AccountDB) CreateAccount(account dto.AccountDto) (*dto.AccountDto, error) {
	if a, _ := db.GetAccount(account.Username); a != nil {
		return nil, fmt.Errorf("user '%s' already created", a.Username)
	} else {
		if res, err := db.sqlDB.Exec(`insert into "accounts"("username", "password") values($1, $2)`, account.Username, account.Password); err != nil {
			return nil, err
		} else {
			id, _ := res.LastInsertId()
			account.Id = int(id)
			return &account, nil
		}
	}
}

func (db *AccountDB) GetAccount(username string) (*dto.AccountDto, error) {
	var id int
	var password string
	if err := db.sqlDB.QueryRow(`SELECT id, username, password FROM accounts WHERE username = $1`, username).Scan(&id, &username, &password); err != nil {
		return nil, err
	} else {
		return &dto.AccountDto{
			Id:       id,
			Username: username,
			Password: password,
		}, nil
	}
}
