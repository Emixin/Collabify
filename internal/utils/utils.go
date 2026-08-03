package utils

import (
	"errors"
	"log"
	"regexp"

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

func UserTeamIDs(session sessions.Session, context *gin.Context, page string) ([]uint, error) {
	// Note: Added a util to don't repeat user teams ids query!

	userID, ok := session.Get("user_id").(string)
	if !ok {
		return []uint{}, errors.New("failed to fetch username")
	}

	if userID == "" {
		return []uint{}, errors.New("failed to fetch username!")
	}

	user_teams := []models.Team{}
	err := database.DB.Preload("Members").Joins("JOIN team_users ON team_users.team_id=teams.id").Where("team_users.user_id=?", userID).Find(&user_teams).Error

	if err != nil {
		return []uint{}, errors.New("failed to find user teams!")
	}

	user_teams_ids := []uint{}
	for _, team := range user_teams {
		user_teams_ids = append(user_teams_ids, uint(team.ID))
	}

	return user_teams_ids, nil
}

// TODO: Define a function to validate emails
func EmailValidator(email string) bool {
	valid := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@gmail\.com$`).MatchString(email)
	if !valid {
		return false
	}
	return true
}

var (
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasLowercase = regexp.MustCompile(`[a-z]`)
	hasSymbol    = regexp.MustCompile(`[!@#$%^&*()_+=.-]`)
)

// TODO: Define a function to validate passwords
func PasswordValidator(password string) (is_valid bool, messages []string) {
	is_valid = true

	if len(password) < 8 {
		messages = append(messages, "Passwords must have at least 8 characters")
		is_valid = false
	}

	var valid bool

	valid = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=.-]{8,}$`).MatchString(password)
	if !valid {
		messages = append(messages, "Password can only contain letters, numbers, and special characters (e.g., !@#)")
		is_valid = false
	}

	valid = hasLowercase.MatchString(password)
	if !valid {
		messages = append(messages, "Password should contain at least one lowercase letter")
		is_valid = false
	}

	valid = hasUppercase.MatchString(password)
	if !valid {
		messages = append(messages, "Password should contain at least one uppercase letter")
		is_valid = false
	}

	valid = hasSymbol.MatchString(password)
	if !valid {
		messages = append(messages, "Password should contain at least one special character (e.g., !@#)")
		is_valid = false
	}

	return
}
