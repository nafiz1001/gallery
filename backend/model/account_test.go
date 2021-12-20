package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nafiz1001/gallery-go/dto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func AccountDBInit(t *testing.T) AccountDB {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	db := DB{GormDB: gormDB}

	var accountDB AccountDB
	err = accountDB.Init(&db)
	require.NoError(t, err)

	return accountDB
}

func CreateAccount(t *testing.T, db AccountDB, username string, password string) dto.AccountDto {
	dto := dto.AccountDto{
		Username: username,
		Password: password,
	}

	account, err := db.CreateAccount(dto)
	require.NoError(t, err)
	require.NotNil(t, account)

	return *account
}

func TestAccountDBInit(t *testing.T) {
	db := AccountDBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()
}

func TestCreateAccount(t *testing.T) {
	db := AccountDBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()

	// successful create
	account1 := CreateAccount(t, db, "username", "password")
	assert.Equal(t, account1.Username, "username")

	// duplicate username
	account2, err := db.CreateAccount(dto.AccountDto{
		Username: "username",
		Password: "password",
	})
	assert.Error(t, err)
	assert.Nil(t, account2)

	// successful second create
	account2, err = db.CreateAccount(dto.AccountDto{
		Username: "username2",
		Password: "password2",
	})
	if assert.NoError(t, err) && assert.NotNil(t, account2) {
		assert.NotEqual(t, account2.Id, account1.Id)
		assert.NotEqual(t, account2.Username, account1.Username)
	}
}

func TestGetAccountById(t *testing.T) {
	db := AccountDBInit(t)
	account := CreateAccount(t, db, "username", "password")
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()

	// successful get account by id
	if dto, err := db.GetAccountById(account.Id); assert.NoError(t, err) && assert.NotNil(t, dto) {
		assert.Equal(t, dto.Id, account.Id)
		assert.Equal(t, dto.Username, account.Username)
	}

	// get account by non-existent id
	dto, err := db.GetAccountById(420)
	assert.Error(t, err)
	assert.Nil(t, dto)
}

func TestGetAccountByUsername(t *testing.T) {
	db := AccountDBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()
	account := CreateAccount(t, db, "username", "password")

	// successful get account by username
	dto, err := db.GetAccountByUsername(account.Username)
	if assert.NoError(t, err) && assert.NotNil(t, dto) {
		assert.Equal(t, dto.Id, account.Id)
		assert.Equal(t, dto.Username, dto.Username)
	}

	// get account by non-existent username
	dto, err = db.GetAccountByUsername("username2")
	assert.Error(t, err)
	assert.Nil(t, dto)
}
