package utils

import (
	"log"
	"net/http"

	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"
	"github.com/gin-contrib/sessions"
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

func UserTeamIDs(session sessions.Session, context *gin.Context, page string) ([]uint, bool) {
	// Note: Added a util to don't repeat user teams ids query!

	userID := session.Get("user_id")

	if userID == nil {
		context.HTML(http.StatusBadRequest, page, gin.H{
			"message": "you are not logged in!",
		})
		return []uint{}, true
	}

	user_teams := []models.Team{}
	err := database.DB.Preload("Members").Joins("JOIN team_users ON team_users.team_id=teams.id").Where("team_users.user_id=?", userID).Find(&user_teams).Error

	if err != nil {
		ErrorCatcher(err, context, http.StatusInternalServerError, page, "failed to query db")
		return []uint{}, true
	}

	user_teams_ids := []uint{}
	for _, team := range user_teams {
		user_teams_ids = append(user_teams_ids, uint(team.ID))
	}

	return user_teams_ids, false
}
