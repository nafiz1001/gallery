package model

import (
	"github.com/nafiz1001/gallery-go/dto"
	"gorm.io/gorm"
)

type ArtDB struct {
	db *gorm.DB
}

type Art struct {
	gorm.Model
	Quantity  int
	Title     string
	AccountID uint
}

func (model *Art) fromDto(data dto.ArtDto) {
	model.ID = uint(data.Id)
	model.Quantity = data.Quantity
	model.Title = data.Title
}

func (model *Art) toDto() *dto.ArtDto {
	return &dto.ArtDto{
		Id:       int(model.ID),
		Quantity: model.Quantity,
		Title:    model.Title,
		AuthorId: int(model.AccountID),
	}
}

func (db *ArtDB) Init(database *DB) error {
	db.db = database.GormDB
	return db.db.AutoMigrate(&Art{})
}

func (db *ArtDB) CreateArt(art dto.ArtDto, account dto.AccountDto) (*dto.ArtDto, error) {
	var artModel Art
	artModel.fromDto(art)

	var accModel Account
	accModel.fromDto(account)

	if err := db.db.Create(&artModel).Error; err != nil {
		return nil, err
	}

	err := db.db.Model(&accModel).Association("Arts").Append(&artModel)
	return artModel.toDto(), err
}

func (db *ArtDB) GetArt(id int) (*dto.ArtDto, error) {
	var model Art
	err := db.db.First(&model, id).Error
	return model.toDto(), err
}

func (db *ArtDB) GetArts() ([]dto.ArtDto, error) {
	var models []Art

	err := db.db.Find(&models).Error
	if err != nil {
		return []dto.ArtDto{}, err
	}

	arts := []dto.ArtDto{}
	for _, m := range models {
		arts = append(arts, *m.toDto())
	}

	return arts, nil
}

func (db *ArtDB) UpdateArt(art dto.ArtDto) (*dto.ArtDto, error) {
	var model Art
	model.fromDto(art)
	err := db.db.Model(&model).Omit("account_id").Updates(&model).Error
	return model.toDto(), err
}

func (db *ArtDB) DeleteArt(id int) (*dto.ArtDto, error) {
	var artModel Art
	if err := db.db.First(&artModel, id).Error; err != nil {
		return nil, err
	}

	var accModel Account
	if err := db.db.First(&accModel, artModel.AccountID).Error; err != nil {
		return nil, err
	}

	if err := db.db.Model(&accModel).Association("Arts").Delete(&artModel); err != nil {
		return nil, err
	}

	err := db.db.Delete(&artModel, id).Error
	return artModel.toDto(), err
}
