package utils

import (
	"log"

	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"
	"github.com/gin-gonic/gin"

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

func ErrorCatcher(err error, context *gin.Context, status_code int, template_name string, message string) {
	// Added error catching util to Don't Repeat Myself!
	log.Println(err)
	context.HTML(status_code, template_name, gin.H{
		"message": message,
	})
}
