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

func DtoToArt(data dto.ArtDto) Art {
	var art Art
	art.ID = uint(data.Id)
	art.Quantity = data.Quantity
	art.Title = data.Title

	return art
}

func (model *Art) ToDto() *dto.ArtDto {
	return &dto.ArtDto{
		Id:       uint(model.ID),
		Quantity: model.Quantity,
		Title:    model.Title,
		AuthorId: uint(model.AccountID),
	}
}

func (db *ArtDB) Init(database *DB) error {
	db.db = database.GormDB
	return db.db.AutoMigrate(&Art{})
}

func (db *ArtDB) CreateArt(art dto.ArtDto, account dto.AccountDto) (*dto.ArtDto, error) {
	artModel := DtoToArt(art)

	accModel := DtoToAccount(account)

	if err := db.db.Create(&artModel).Error; err != nil {
		return nil, err
	} else if err := db.db.Model(&accModel).Association("Arts").Append(&artModel); err != nil {
		return nil, err
	} else {
		return artModel.ToDto(), err
	}
}

func (db *ArtDB) GetArt(id int) (*dto.ArtDto, error) {
	var model Art
	if err := db.db.First(&model, id).Error; err != nil {
		return nil, err
	} else {
		return model.ToDto(), err
	}
}

func (db *ArtDB) GetArts() ([]dto.ArtDto, error) {
	var models []Art

	if err := db.db.Find(&models).Error; err != nil {
		return []dto.ArtDto{}, err
	} else {
		arts := []dto.ArtDto{}
		for _, m := range models {
			arts = append(arts, *m.ToDto())
		}
		return arts, nil
	}
}

func (db *ArtDB) UpdateArt(art dto.ArtDto) (*dto.ArtDto, error) {
	model := DtoToArt(art)
	if err := db.db.Model(&model).Omit("account_id").Updates(&model).Error; err != nil {
		return nil, err
	} else {
		return model.ToDto(), err
	}
}

func (db *ArtDB) DeleteArt(id int) (*dto.ArtDto, error) {
	var artModel Art
	var accModel Account

	if err := db.db.First(&artModel, id).Error; err != nil {
		return nil, err
	} else if err := db.db.First(&accModel, artModel.AccountID).Error; err != nil {
		return nil, err
	} else if err := db.db.Model(&accModel).Association("Arts").Delete(&artModel); err != nil {
		return nil, err
	} else if err := db.db.Delete(&artModel, id).Error; err != nil {
		return nil, err
	} else {
		return artModel.ToDto(), err
	}
}
