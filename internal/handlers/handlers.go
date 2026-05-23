package handlers

import (
	"collabify/internal/database"
	"collabify/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomepageHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "home.html", nil)
}

func LoginpageHandler(context *gin.Context) {

	if context.Request.Method == "GET" {
		log.Println("entered if block")
		context.HTML(http.StatusOK, "login.html", nil)
	} else if context.Request.Method == "POST" {
		log.Println("entered else if block")
		context.Request.ParseForm()
		username := context.Request.FormValue("username")
		password := context.Request.FormValue("password")

		if username == "" || password == "" {
			context.HTML(http.StatusOK, "login.html", gin.H{
				"message": "Please enter both username and password and submit",
			})
		}
	}
}

func SignuppageHandler(context *gin.Context) {
	if context.Request.Method == "GET" {
		context.HTML(http.StatusOK, "signup.html", nil)
	} else if context.Request.Method == "POST" {
		log.Println("else if block")
		context.Request.ParseForm()
		username := context.Request.FormValue("username")
		email := context.Request.FormValue("email")
		password := context.Request.FormValue("password")
		confirm_password := context.Request.FormValue("confirm_password")
		if username == "" || password == "" || confirm_password == "" {
			context.HTML(http.StatusOK, "signup.html", gin.H{
				"message": "Please enter all fields then submit",
			})
		} else if password != confirm_password {
			context.HTML(http.StatusOK, "signup.html", gin.H{
				"message": "passwords did not match!",
			})
		} else {
			user := models.User{
				Username: username,
				Email:    email,
				Type:     "NoType",
				Score:    0,
			}

			err := database.DB.Create(&user).Error
			if err != nil {
				log.Println(err)
				context.HTML(http.StatusInternalServerError, "signup.html", gin.H{
					"message": "failed to create new user!",
				})
				return
			}

			context.HTML(http.StatusOK, "signup.html", gin.H{
				"message": "new user created!",
			})
		}
	}
}

func UserslistHandler(context *gin.Context) {
	users_list := []models.User{}
	err := database.DB.Find(&users_list).Error
	if err != nil {
		log.Println(err)
		context.HTML(http.StatusOK, "users_list.html", gin.H{
			"message": "failed to query db",
		})
		return
	}

	context.HTML(http.StatusOK, "users_list.html", gin.H{
		"users_list": users_list,
		"message":    "user table fetched successfully!",
	})
}

func TeamslistHandler(context *gin.Context) {
	teams_list := []models.Team{}
	err := database.DB.Preload("Members").Preload("Leader").Find(&teams_list).Error

	if err != nil {
		log.Println(err)
		context.HTML(http.StatusOK, "users_list.html", gin.H{
			"message": "failed to query db",
		})
		return
	}

	context.HTML(http.StatusOK, "teams_list.html", gin.H{
		"teams_list": teams_list,
	})
}

func TasklistHandler(context *gin.Context) {
	tasks_list := []models.Task{}
	err := database.DB.Preload("Team").Preload("Team.Leader").Preload("Team.Members").Find(&tasks_list).Error

	if err != nil {
		log.Println(err)
		context.HTML(http.StatusInternalServerError, "tasks_list.html", gin.H{
			"message": "failed to query db",
		})
		return
	}

	context.HTML(http.StatusOK, "tasks_list.html", gin.H{
		"tasks_list": tasks_list,
	})
}
