package handlers

import (
	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"
	"github.com/gin-contrib/sessions"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HomepageHandler(context *gin.Context) {

	session := sessions.Default(context)
	username := session.Get("username")

	if username == "" {
		context.HTML(http.StatusOK, "home.html", gin.H{
			"message":  "",
			"username": "Anonymous User",
		})
	} else {
		context.HTML(http.StatusOK, "home.html", gin.H{
			"message":  "You have 0 task and none of them is pending",
			"username": username,
		})
	}
}

func LoginpageHandler(context *gin.Context) {

	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "login.html", nil)

	case "POST":
		context.Request.ParseForm()
		username := context.Request.FormValue("username")
		password := context.Request.FormValue("password")

		if username == "" || password == "" {
			context.HTML(http.StatusBadRequest, "login.html", gin.H{
				"message": "Please enter both username and password and submit",
			})
			return
		}

		var user_obj models.User
		err := database.DB.Where(&models.User{Username: username}).First(&user_obj).Error

		if err != nil {
			log.Println(err)
			context.HTML(http.StatusNotFound, "login.html", gin.H{
				"message": "either username or password is incorrect!",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user_obj.PasswordHash), []byte(password))
		if err != nil {
			log.Println(err)
			context.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"message": "either username or password is incorrect!",
			})
			return
		}

		session := sessions.Default(context)
		session.Set("user_id", user_obj.ID)
		session.Set("username", user_obj.Username)
		session.Save()

		context.HTML(http.StatusOK, "dashboard.html", gin.H{
			"message":  "logged in successfully!",
			"username": username,
		})

	default:
		context.HTML(http.StatusMethodNotAllowed, "login.html", gin.H{
			"message": "Method not allowed!",
		})
	}
}

func SignuppageHandler(context *gin.Context) {
	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "signup.html", nil)
	case "POST":
		context.Request.ParseForm()
		username := context.Request.FormValue("username")
		email := context.Request.FormValue("email")
		password := context.Request.FormValue("password")
		confirm_password := context.Request.FormValue("confirm_password")

		switch {
		case username == "" || password == "" || confirm_password == "" || email == "":
			context.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"message": "Please enter all fields then submit",
			})
		case password != confirm_password:
			context.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"message": "passwords did not match!",
			})
		default:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			if err != nil {
				log.Println(err)
				context.HTML(http.StatusInternalServerError, "signup.html", gin.H{
					"message": "failed to create new user!",
				})
				return
			}

			user := models.User{
				Username:     username,
				PasswordHash: string(hashedPassword),
				Email:        email,
				Type:         "NoType",
				Score:        0,
			}

			err = database.DB.Create(&user).Error
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
	default:
		context.HTML(http.StatusMethodNotAllowed, "signup.html", gin.H{
			"message": "Method not allowed!",
		})
	}
}

func UserslistHandler(context *gin.Context) {
	users_list := []models.User{}
	err := database.DB.Find(&users_list).Error
	if err != nil {
		log.Println(err)
		context.HTML(http.StatusInternalServerError, "users_list.html", gin.H{
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
			context.HTML(http.StatusInternalServerError, "users_list.html", gin.H{
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

// TODO: Complete this function later!
func UserDashboardHandler(context *gin.Context) {
	session := sessions.Default(context)
	username := session.Get("username")
	context.HTML(http.StatusOK, "dashboard.html", gin.H{
		"username": username,
	})
}
