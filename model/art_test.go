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

func ArtDBInit(t *testing.T, gormDB *gorm.DB) model.ArtDB {
	db := model.DB{GormDB: gormDB}

	var artDB model.ArtDB
	err := artDB.Init(&db)
	require.NoError(t, err)

	return artDB
}

func CreateArt(t *testing.T, artDB model.ArtDB, artDto dto.ArtDto) dto.ArtDto {
	dto, err := artDB.CreateArt(artDto)
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, dto.Quantity, artDto.Quantity)
	assert.Equal(t, dto.Title, artDto.Title)
	assert.Equal(t, dto.AuthorId, artDto.AuthorId)

	return *dto
}

func createUserAndArt(t *testing.T, accountDB model.AccountDB, artDB model.ArtDB, username string) (dto.AccountDto, dto.ArtDto) {
	accountDto := CreateAccount(t, accountDB, username, "password")
	artDto := CreateArt(t, artDB, dto.ArtDto{
		Id:       0,
		Quantity: 1,
		Title:    "title",
		AuthorId: accountDto.Id,
	})

	return accountDto, artDto
}

func TestArtDBInit(t *testing.T) {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	ArtDBInit(t, gormDB)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()
}

func TestCreateArt(t *testing.T) {
	accountDB, gormDB := AccountDBInit(t)
	artDB := ArtDBInit(t, gormDB)
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
	dto1 := CreateArt(t, artDB, artDto)

	// successfully create art with duplicate title
	artDto = dto.ArtDto{
		Id:       0,
		Quantity: 2,
		Title:    "title",
		AuthorId: accountDto.Id,
	}
	dto2 := CreateArt(t, artDB, artDto)
	assert.NotEqual(t, dto1.Id, dto2.Id)

	// fail to create art for non-existent account
	dto3, err := artDB.CreateArt(dto.ArtDto{
		Id:       0,
		Quantity: 2,
		Title:    "title",
		AuthorId: 420,
	})
	require.Error(t, err)
	require.Nil(t, dto3)

	// create art for another account successfully
	_, dto4 := createUserAndArt(t, accountDB, artDB, "username2")
	assert.NotEqual(t, dto4.AuthorId, dto1.AuthorId)
}

func TestGetArt(t *testing.T) {
	accountDB, gormDB := AccountDBInit(t)
	artDB := ArtDBInit(t, gormDB)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()

	dto, err := artDB.GetArt(0)
	assert.Error(t, err)
	assert.Nil(t, dto)

	_, artDto := createUserAndArt(t, accountDB, artDB, "username")

	// get art successfully
	dto, err = artDB.GetArt(artDto.Id)
	if assert.NoError(t, err) && assert.NotNil(t, dto) {
		assert.Equal(t, *dto, artDto)
	}

	// get art successfully from another account
	_, artDto2 := createUserAndArt(t, accountDB, artDB, "username2")
	dto, err = artDB.GetArt(artDto2.Id)
	if assert.NoError(t, err) && assert.NotNil(t, dto) {
		assert.Equal(t, *dto, artDto2)
		assert.NotEqual(t, dto.AuthorId, artDto.AuthorId)
	}

	// don't get art
	dto, err = artDB.GetArt(420)
	assert.Error(t, err)
	assert.Nil(t, dto)
}

func TestGetArts(t *testing.T) {
	accountDB, gormDB := AccountDBInit(t)
	artDB := ArtDBInit(t, gormDB)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()

	// there should be zero arts present
	artDtos, err := artDB.GetArts()
	assert.NoError(t, err)
	assert.Equal(t, len(artDtos), 0)

	// there should be only 1 art present
	accountDto, artDto := createUserAndArt(t, accountDB, artDB, "username")
	artDtos, err = artDB.GetArts()
	if assert.NoError(t, err) && assert.Equal(t, len(artDtos), 1) {
		if assert.NotNil(t, artDtos[0]) {
			assert.Equal(t, artDtos[0], artDto)
		}
	}

	// there should be only 2 arts present
	artDto2 := CreateArt(t, artDB, dto.ArtDto{
		Id:       0,
		Quantity: 2,
		Title:    "title2",
		AuthorId: accountDto.Id,
	})
	artDtos, err = artDB.GetArts()
	if assert.NoError(t, err) {
		assert.Equal(t, len(artDtos), 2)
		artFound1 := false
		artFound2 := false
		for _, currDto := range artDtos {
			if currDto.Id == artDto.Id {
				assert.Equal(t, currDto, artDto)
				artFound1 = true
			} else if currDto.Id == artDto2.Id {
				assert.Equal(t, currDto, artDto2)
				artFound2 = true
			} else {
				t.Errorf("ArtDto with id '%d' should not exist", currDto.Id)
			}
		}
		assert.Truef(t, artFound1, "%v not found", artDto)
		assert.Truef(t, artFound2, "%v not found", artDto2)
	}
}

func TestUpdateArt(t *testing.T) {
	accountDB, gormDB := AccountDBInit(t)
	artDB := ArtDBInit(t, gormDB)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()
	_, artDto := createUserAndArt(t, accountDB, artDB, "username")

	// successfully update art
	artDto2, err := artDB.UpdateArt(dto.ArtDto{
		Title:    "new_title",
		Quantity: 2,
		AuthorId: artDto.AuthorId,
		Id:       artDto.Id,
	})
	if assert.NoError(t, err) && assert.NotNil(t, artDto2) {
		assert.Equal(t, "new_title", artDto2.Title)
		assert.Equal(t, 2, artDto2.Quantity)
		assert.Equal(t, artDto.Id, artDto2.Id)
		assert.Equal(t, artDto2.AuthorId, artDto.AuthorId)

		artDto22, err := artDB.GetArt(artDto.Id)
		if assert.NoError(t, err) && assert.NotNil(t, artDto22) {
			assert.Equal(t, artDto2, artDto22)
		}
	}

	// don't transfer ownership to an account that does not exist
	artDto3, err := artDB.UpdateArt(dto.ArtDto{
		Title:    "new_title",
		Quantity: 2,
		AuthorId: 420,
		Id:       artDto.Id,
	})
	assert.Error(t, err)
	assert.Nil(t, artDto3)

	// transfer ownership of art to an existing account
	accountDto2 := CreateAccount(t, accountDB, "username2", "password2")
	artDto4, err := artDB.UpdateArt(dto.ArtDto{
		Title:    "new_title",
		Quantity: 2,
		AuthorId: accountDto2.Id,
		Id:       artDto.Id,
	})
	if assert.NoError(t, err) && assert.NotNil(t, artDto4) {
		assert.Equal(t, artDto.Id, artDto4.Id)
		assert.Equal(t, accountDto2.Id, artDto4.AuthorId)
		artDto42, err := artDB.GetArt(artDto.Id)
		if assert.NoError(t, err) && assert.NotNil(t, artDto42) {
			assert.Equal(t, artDto4, artDto42)
		}
	}

	// don't update art that does not exist
	artDto5, err := artDB.UpdateArt(dto.ArtDto{
		Title:    "new_title",
		Quantity: 2,
		AuthorId: accountDto2.Id,
		Id:       420,
	})
	assert.Error(t, err)
	assert.Nil(t, artDto5)
}

func TestDeleteArt(t *testing.T) {
	accountDB, gormDB := AccountDBInit(t)
	artDB := ArtDBInit(t, gormDB)
	defer func() {
		otherDb, _ := gormDB.DB()
		otherDb.Close()
	}()
	_, artDto := createUserAndArt(t, accountDB, artDB, "username")

	// successful delete art
	artDtoTemp, err := artDB.DeleteArt(artDto.Id)
	if assert.NoError(t, err) && assert.Equal(t, artDto, *artDtoTemp) {
		artDto, err := artDB.GetArt(artDto.Id)
		assert.Error(t, err)
		assert.Nil(t, artDto)
	}

	// can't delete non-existent art
	artDto2, err := artDB.DeleteArt(420)
	assert.Error(t, err)
	assert.Nil(t, artDto2)
}
