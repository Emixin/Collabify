package database

import (
	"collabify/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data/db.sqlite"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	DB = db

	return DB
}
