package model

import (
	"gorm.io/gorm"
)

type DB struct {
	GormDB *gorm.DB
}
