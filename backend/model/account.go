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

func DtoToAccount(data dto.AccountDto) Account {
	var model Account
	model.ID = uint(data.Id)
	model.Username = data.Username
	model.Password = data.Password

	return model
}

func (model *Account) ToDto() *dto.AccountDto {
	return &dto.AccountDto{
		Id:       int(model.ID),
		Username: model.Username,
		Password: model.Password,
	}
}

func (db *AccountDB) Init(database *DB) error {
	db.db = database.GormDB
	return db.db.AutoMigrate(&Account{})
}

func (db *AccountDB) CreateAccount(account dto.AccountDto) (*dto.AccountDto, error) {
	if _, err := db.GetAccountByUsername(account.Username); err == nil {
		return nil, fmt.Errorf("user '%s' already created", account.Username)
	} else {
		model := DtoToAccount(account)
		model.Arts = []Art{}
		if err := db.db.Create(&model).Error; err != nil {
			return nil, err
		} else {
			return model.ToDto(), err
		}
	}
}

func (db *AccountDB) GetAccountById(id int) (*dto.AccountDto, error) {
	var model Account
	if err := db.db.First(&model, id).Error; err != nil {
		return nil, err
	} else {
		return model.ToDto(), err
	}
}

func (db *AccountDB) GetAccountByUsername(username string) (*dto.AccountDto, error) {
	var model Account
	if err := db.db.First(&model, "username = ?", username).Error; err != nil {
		return nil, err
	} else {
		return model.ToDto(), err
	}
}
