package utils

import (
	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"

	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) {
	// Used RAM for storing the test database here
	db, err := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	database.DB = db

}
