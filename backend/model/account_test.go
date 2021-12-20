package model

import (
	"testing"

	"github.com/nafiz1001/gallery-go/dto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBInit(t *testing.T) AccountDB {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	db := DB{GormDB: gormDB}

	var accountDB AccountDB
	err = accountDB.Init(&db)
	if err != nil {
		t.Fatal(err)
	}

	return accountDB
}

func CreateAccount(t *testing.T, db AccountDB, username string, password string) *dto.AccountDto {
	dto := dto.AccountDto{
		Username: username,
		Password: password,
	}

	account, err := db.CreateAccount(dto)
	if err != nil {
		t.Fatal(err)
	} else if account == nil {
		t.Fatalf("account nil when creating account with username '%s'", dto.Username)
	}

	return account
}

func TestAccountDBInit(t *testing.T) {
	db := DBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()
}

func TestCreateAccount(t *testing.T) {
	db := DBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()

	// successful create
	account1 := CreateAccount(t, db, "username", "password")

	// duplicate username
	account2, err := db.CreateAccount(dto.AccountDto{
		Username: "username",
		Password: "password",
	})
	if err == nil {
		t.Fatalf("err should be nil for identical username '%s'", account2.Username)
	} else if account2 != nil {
		t.Fatalf("account2 should be nil for identical username '%s'", account2.Username)
	}

	// successful second create
	account2, _ = db.CreateAccount(dto.AccountDto{
		Username: "username2",
		Password: "password2",
	})
	if account2 != nil {
		if account2.Id == account1.Id {
			t.Fatalf("expected to not create account with identical id '%d'", account2.Id)
		} else if account2.Username == account1.Username {
			t.Fatalf("expected to not create account with identical username '%s'", account2.Username)
		}
	}
}

func TestGetAccountById(t *testing.T) {
	db := DBInit(t)
	account := CreateAccount(t, db, "username", "password")
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()

	// successful get account by id
	dto, err := db.GetAccountById(account.Id)
	if err != nil {
		t.Fatal(err)
	} else if dto == nil {
		t.Fatalf("dto returned nil for id '%d'", account.Id)
	} else if account.Id != dto.Id {
		t.Fatalf("id returned '%d' does not match dto's id '%d'", account.Id, dto.Id)
	}

	// get account by non-existent id
	dto, err = db.GetAccountById(420)
	if err == nil {
		t.Fatalf("error not nil for non-existent id '%d'", 420)
	} else if dto != nil {
		t.Fatalf("dto returned for non-existent id '%d'", 420)
	}
}

func TestGetAccountByUsername(t *testing.T) {
	db := DBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()
	account := CreateAccount(t, db, "username", "password")

	// successful get account by username
	dto, err := db.GetAccountByUsername(account.Username)
	if err != nil {
		t.Fatal(err)
	} else if account == nil {
		t.Fatalf("account nil for existing username '%s'", account.Username)
	} else if account.Username != dto.Username {
		t.Fatalf("dto's username '%s' does not match username '%s", account.Username, dto.Username)
	}

	// get account by non-existent username
	dto, err = db.GetAccountByUsername("username2")
	if err == nil {
		t.Fatalf("error is nil for non-existent username")
	} else if dto != nil {
		t.Fatalf("dto returned for non-existent username")
	}
}
