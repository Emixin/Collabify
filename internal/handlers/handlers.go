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

func CreateTeamHandler(context *gin.Context) {
	if context.Request.Method == "GET" {
		context.HTML(http.StatusOK, "create_team.html", nil)
	} else if context.Request.Method == "POST" {
		context.Request.ParseForm()
		name := context.Request.FormValue("Name")
		leader_name := context.Request.FormValue("Leader")

		var leader_obj models.User
		err := database.DB.Where(&models.User{Username: leader_name}).First(&leader_obj).Error
		if err != nil {
			log.Println(err)
			context.HTML(http.StatusOK, "users_list.html", gin.H{
				"message": "failed to query db",
			})
			return
		}

		team := models.Team{
			Name:   name,
			Leader: leader_obj,
		}
		err = database.DB.Create(&team).Error
		if err != nil {
			context.HTML(http.StatusInternalServerError, "create_team.html", gin.H{
				"message": "falied to create new team!",
			})
			return
		}

		context.HTML(http.StatusOK, "create_team.html", gin.H{
			"message": "new team created!",
		})
	}

}

func DeleteTeamHandler(context *gin.Context) {
	if context.Request.Method == "GET" {
		context.HTML(http.StatusOK, "delete_team.html", nil)
	} else if context.Request.Method == "POST" {
		context.Request.ParseForm()
		name := context.Request.FormValue("Name")
		var team models.Team
		err := database.DB.Where(&models.Team{Name: name}).First(&team)
		if err != nil {
			context.HTML(http.StatusInternalServerError, "delete_team.html", gin.H{
				"message": "could not find the team!",
			})
			return
		}
		if team.Leader.Username == context.Request.URL.User.Username() {
			err := database.DB.Where(&models.Team{Name: name}).Delete(&models.Team{}).Error
			if err != nil {
				context.HTML(http.StatusInternalServerError, "delete_team.html", gin.H{
					"message": "could not delete the user!",
				})
				return
			}
			context.HTML(http.StatusOK, "delete_team.html", gin.H{
				"message": "deleted the team!",
			})
		} else {
			context.HTML(http.StatusInternalServerError, "delete_team.html", gin.H{
				"message": "only team leader can perform this action!",
			})
		}
	}
}

func UpdateTeamHandler(context *gin.Context) {
	// TODO: complete this function later!
}

func TeamslistHandler(context *gin.Context) {
	teams_list := []models.Team{}
	err := database.DB.Preload("Members").Preload("Leader").Find(&teams_list).Error

	if err != nil {
		log.Println(err)
		context.HTML(http.StatusInternalServerError, "users_list.html", gin.H{
			"message": "failed to query db",
		})
		return
	}

	context.HTML(http.StatusOK, "teams_list.html", gin.H{
		"teams_list": teams_list,
	})
}

func CreateTaskHandler(context *gin.Context) {
	if context.Request.Method == "GET" {
		context.HTML(http.StatusOK, "create_task.html", nil)
	} else if context.Request.Method == "POST" {
		name := context.Request.FormValue("Name")
		team_name := context.Request.FormValue("Team")
		deadline := context.Request.FormValue("Deadline")

		var team_obj models.Team
		err := database.DB.Where(&models.Team{Name: team_name}).First(&team_obj).Error

		if err != nil {
			log.Println(err)
			context.HTML(http.StatusInternalServerError, "users_list.html", gin.H{
				"message": "team not found!",
			})
			return
		}

		task := models.Task{
			Name:     name,
			Team:     team_obj,
			Deadline: deadline,
		}
		database.DB.Create(&task)

		context.HTML(http.StatusOK, "create_task.html", gin.H{
			"message": "new task created!",
		})
	}
}

func DeleteTaskHandler(context *gin.Context) {
	if context.Request.Method == "GET" {
		context.HTML(http.StatusOK, "delete_task.html", nil)
	} else if context.Request.Method == "POST" {
		name := context.Request.FormValue("Name")
		team_name := context.Request.FormValue("Team")

		var team_obj models.Team
		err := database.DB.Where(&models.Team{Name: team_name}).Find(&team_obj).Error
		if err != nil {
			log.Println(err)
			context.HTML(http.StatusInternalServerError, "tasks_list.html", gin.H{
				"message": "team not found!",
			})
			return
		}

		err = database.DB.Where(&models.Task{Name: name, Team: team_obj}).Delete(&models.Task{}).Error
		if err != nil {
			log.Println(err)
			context.HTML(http.StatusInternalServerError, "tasks_list.html", gin.H{
				"message": "failed to delete the task!",
			})
			return
		}

		context.HTML(http.StatusOK, "delete_task.html", gin.H{
			"message": "task has been deleted!",
		})
	}
}

func UpdateTaskHandler(context *gin.Context) {
	// TODO: complete this function later!
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
