package database

import (
	"diandian/background/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() error {
	db, err := gorm.Open(sqlite.Open("data.dd"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db

	db.AutoMigrate(
		&model.Conversation{},
		&model.Message{},
		&model.Task{},
		&model.Step{},
		&model.Setting{},
	)

	return nil
}
