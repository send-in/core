package database

import (
	model "core/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Account{},
		&model.Connection{},
		&model.Template{},
		&model.Message{},
	)
}