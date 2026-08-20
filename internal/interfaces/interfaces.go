package interfaces

import (
	"net/http"

	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//TODO: Write your APIs and REST APIs here!

type CreateUserAPIRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Create
func CreateUserAPIHandler(context *gin.Context) {
	//TODO: Complete it later
	var request CreateUserAPIRequest

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	passwordHash, err3 := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err3 != nil {
		return
	}

	err2 := database.DB.Create(&models.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		Type:         models.NoType,
		Score:        0,
	}).Error

	if err2 != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err2.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created!",
	})
}

// Read
func GetUsersAPIHandler(context *gin.Context) {
	//TODO: Complete it later
	var users []models.User
	err := database.DB.Find(&users).Error
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var userStrs []string

	//NOTE: Added this part to describe the users!
	for _, user := range users {
		userStr := user.String()
		userStrs = append(userStrs, userStr)

	}

	context.JSON(http.StatusOK, gin.H{
		"message":    "Users list fetched successfully!",
		"users":      userStrs,
		"usersCount": len(users),
	})
}

func GetUserAPIHandler(context *gin.Context) {
	//TODO: Complete it later
}

// Update
func UpdateUserAPIHandler(context *gin.Context) {
	//TODO: Complete it later
}

type DeleteUserAPIRequest struct {
	Username string `json:"username"`
}

// Delete
func DeleteUserAPIHandler(context *gin.Context) {
	//TODO: Complete it later
	var request DeleteUserAPIRequest

	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//FIX: It seems the Delete function does not fail when i try to delete a non-existing user
	err2 := database.DB.Where(&models.User{Username: request.Username}).Delete(&models.User{}).Error
	if err2 != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":    err2.Error(),
			"username": request.Username,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "user deleted!",
	})

}
