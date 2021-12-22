package model

import (
	"fmt"

	"github.com/nafiz1001/gallery-go/dto"
	"gorm.io/gorm"
)

type AccountDB struct {
	db *gorm.DB
}

type Account struct {
	gorm.Model
	Username string
	Password string
	Arts     []Art `gorm:"foreignKey:AccountID"`
}

// Creates Account object from AccountDto.
// It does not add the account to the database.
func DtoToAccount(data dto.AccountDto) Account {
	var model Account
	model.ID = uint(data.Id)
	model.Username = data.Username
	model.Password = data.Password
	model.Arts = []Art{}

	return model
}

// Converts Acccount to AccountDto.
func (model *Account) ToDto() *dto.AccountDto {
	return &dto.AccountDto{
		Id:       uint(model.ID),
		Username: model.Username,
		Password: model.Password,
	}
}

// Creates account table in the database.
func (db *AccountDB) Init(database *DB) error {
	db.db = database.GormDB
	return db.db.AutoMigrate(&Account{})
}

// Creates new account if there is no existing account with identical username.
func (db *AccountDB) CreateAccount(account dto.AccountDto) (*dto.AccountDto, error) {
	account.Id = 0
	if _, err := db.GetAccountByUsername(account.Username); err == nil {
		return nil, fmt.Errorf("username '%s' already exists", account.Username)
	} else {
		model := DtoToAccount(account)
		if err := db.db.Create(&model).Error; err != nil {
			return nil, err
		} else {
			return model.ToDto(), nil
		}
	}
}

// Gets account by id if it exists.
func (db *AccountDB) GetAccountById(id uint) (*dto.AccountDto, error) {
	var model Account
	if err := db.db.First(&model, id).Error; err != nil {
		return nil, err
	} else {
		return model.ToDto(), err
	}
}

// Gets account by username if it exists
func (db *AccountDB) GetAccountByUsername(username string) (*dto.AccountDto, error) {
	var model Account
	if err := db.db.First(&model, "username = ?", username).Error; err != nil {
		return nil, err
	} else {
		return model.ToDto(), err
	}
}
