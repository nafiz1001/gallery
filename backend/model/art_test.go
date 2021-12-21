package model_test

import (
	"testing"

	"github.com/nafiz1001/gallery-go/dto"
	"github.com/nafiz1001/gallery-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ArtDBInit(t *testing.T) (model.ArtDB, *gorm.DB) {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	db := model.DB{GormDB: gormDB}

	var artDB model.ArtDB
	err = artDB.Init(&db)
	require.NoError(t, err)

	return artDB, gormDB
}

func CreateArt(t *testing.T, artDB model.ArtDB, artDto dto.ArtDto, accountDto dto.AccountDto) dto.ArtDto {
	dto, err := artDB.CreateArt(artDto, accountDto)
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, dto.Quantity, artDto.Quantity)
	assert.Equal(t, dto.Title, artDto.Title)
	assert.Equal(t, dto.AuthorId, accountDto.Id)

	return *dto
}

func TestArtDBInit(t *testing.T) {
	_, db := ArtDBInit(t)
	defer func() {
		otherDb, _ := db.DB()
		otherDb.Close()
	}()
}

func TestCreateArt(t *testing.T) {
	accountDB, _ := AccountDBInit(t)
	artDB, gormDB := ArtDBInit(t)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()

	accountDto := CreateAccount(t, accountDB, "username", "password")

	// successfully create art
	artDto := dto.ArtDto{
		Id:       0,
		Quantity: 1,
		Title:    "title",
		AuthorId: accountDto.Id,
	}
	dto1 := CreateArt(t, artDB, artDto, accountDto)

	// successfully create art with duplicate title
	artDto = dto.ArtDto{
		Id:       0,
		Quantity: 2,
		Title:    "title",
		AuthorId: accountDto.Id,
	}
	dto2 := CreateArt(t, artDB, artDto, accountDto)
	assert.NotEqual(t, dto1.Id, dto2.Id)
}

func TestGetArt(t *testing.T) {
	accountDB, _ := AccountDBInit(t)
	artDB, gormDB := ArtDBInit(t)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()
	accountDto := CreateAccount(t, accountDB, "username", "password")
	artDto := dto.ArtDto{
		Id:       0,
		Quantity: 1,
		Title:    "title",
		AuthorId: accountDto.Id,
	}
	artDto = CreateArt(t, artDB, artDto, accountDto)

	dto, err := artDB.GetArt(artDto.Id)
	if assert.NoError(t, err) && assert.NotNil(t, dto) {
		assert.Equal(t, dto.Quantity, artDto.Quantity)
		assert.Equal(t, dto.Title, artDto.Title)
		assert.Equal(t, dto.AuthorId, accountDto.Id)
	}
}
