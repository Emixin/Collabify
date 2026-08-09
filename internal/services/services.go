package services

import (
	"errors"
	"log"
	"slices"

	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

type CreateUserRequest struct {
	Username        string `form:"username"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

func (service *UserService) CreateUser(request CreateUserRequest) (*models.User, error) {
	password := request.Password
	confirm := request.ConfirmPassword

	if password != confirm {
		return nil, errors.New("Passwords did not match!")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	//TODO: Check for duplicate username!
	var all_usernames []string
	err2 := database.DB.Model(&models.User{}).Select("username").Find(&all_usernames).Error
	if err2 != nil {
		return nil, err2
	}

	log.Printf("all usernames: %v", all_usernames)

	if is_duplicate := slices.Contains(all_usernames, request.Username); is_duplicate {
		return nil, errors.New("username is already taken!")
	}

	log.Println("username is not duplicate")
	log.Printf("username: %v", request)

	switch {
	case request.Username == "" || request.Password == "" || request.ConfirmPassword == "" || request.Email == "":
		log.Println("case 1")
		return nil, errors.New("Please enter all fields then submit")
	case request.Password != request.ConfirmPassword:
		log.Println("case 2")
		return nil, errors.New("passwords did not match!")
	default:
		log.Println("default")
		user := &models.User{
			Username:     request.Username,
			Email:        request.Email,
			PasswordHash: string(passwordHash),
			Type:         models.NoType,
			Score:        0,
		}

		if err := service.DB.Create(user).Error; err != nil {
			return nil, err
		}
		log.Println("user is created")

		return user, nil
	}
}
