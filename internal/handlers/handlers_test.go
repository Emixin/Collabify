package handlers

import (
	"collabify/internal/database"
	"collabify/internal/models"
	"collabify/internal/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandlerRejectsEmptyForm(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Any("/login", LoginpageHandler)
	router.LoadHTMLGlob("../../web/templates/*.html")

	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("want status %v got %v", http.StatusBadRequest, response.Code)
	}
}

func TestLoginHandlerRejectsNonExistingUsername(t *testing.T) {
	utils.SetupTestDB(t)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("user1password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to set user password: %v", err)
		return
	}

	database.DB.Create(&models.User{
		Username:     "user1",
		PasswordHash: string(hashedPassword),
	})

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Any("/login", LoginpageHandler)
	router.LoadHTMLGlob("../../web/templates/*.html")

	form := url.Values{}
	form.Set("username", "non_existing_user")
	form.Set("password", "user1password")

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("want status %v got status %v", http.StatusNotFound, response.Code)
	}
}
