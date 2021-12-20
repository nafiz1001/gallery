package model

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ArtDBInit(t *testing.T) ArtDB {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	db := DB{GormDB: gormDB}

	var artDB ArtDB
	err = artDB.Init(&db)
	require.NoError(t, err)

	return artDB
}

func TestArtDBInit(t *testing.T) {
	db := ArtDBInit(t)
	defer func() {
		otherDb, _ := db.db.DB()
		otherDb.Close()
	}()
}
