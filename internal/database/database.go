package database

import (
	"github.com/Emixin/Collabify/internal/models"

	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data/db.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println(err)
	}

	err = db.AutoMigrate(&models.Team{})
	if err != nil {
		log.Println(err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Println(err)
	}

	DB = db

	return DB
}
