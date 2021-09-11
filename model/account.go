package model

import (
	"fmt"

	"github.com/nafiz1001/gallery-go/dto"
	"github.com/nafiz1001/gallery-go/util"
)

type AccountDB struct {
	accounts map[string]*dto.AccountDto
}

func (db *AccountDB) Init() error {
	db.accounts = map[string]*dto.AccountDto{}
	return nil
}

func (db *AccountDB) CreateAccount(account dto.AccountDto) (*dto.AccountDto, error) {
	if a, _ := db.GetAccount(account.Username); a != nil {
		return nil, fmt.Errorf("user '%s' already created", a.Username)
	} else {
		account.Id = util.CreateId()
		db.accounts[account.Id] = &account
		return &account, nil
	}
}

func (db *AccountDB) GetAccount(username string) (*dto.AccountDto, error) {
	for _, a := range db.accounts {
		if a.Username == username {
			return a, nil
		}
	}

	return nil, fmt.Errorf("could not find user '%s'", username)
}
