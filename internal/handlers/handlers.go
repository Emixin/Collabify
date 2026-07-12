package handlers

import (
	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/models"
	"github.com/Emixin/Collabify/internal/utils"
	"github.com/gin-contrib/sessions"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HomepageHandler(context *gin.Context) {

	session := sessions.Default(context)
	username := session.Get("username")
	redirected := session.Get("redirect")

	if redirected != nil {
		session.Delete("redirect")
		session.Save()

		context.HTML(http.StatusOK, "home.html", gin.H{
			"username":         "Anonymous User",
			"redirect_message": "Logged out successfully!",
			"is_authenticated": false,
		})
	} else {
		if username == nil {
			context.HTML(http.StatusOK, "home.html", gin.H{
				"message":          "You have 0 task and none of them is pending",
				"username":         "Anonymous User",
				"is_authenticated": false,
			})
		} else {
			context.HTML(http.StatusOK, "home.html", gin.H{
				"message":          "You have 0 task and none of them is pending",
				"username":         username,
				"is_authenticated": true,
			})
		}
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
			utils.ErrorCatcher(err, context, http.StatusNotFound, "login.html", "either username or password is incorrect!")
			return
		}

		err2 := bcrypt.CompareHashAndPassword([]byte(user_obj.PasswordHash), []byte(password))
		if err2 != nil {
			utils.ErrorCatcher(err2, context, http.StatusUnauthorized, "login.html", "either username or password is incorrect!")
			return
		}

		session := sessions.Default(context)
		session.Set("user_id", user_obj.ID)
		session.Set("username", user_obj.Username)
		session.Set("email", user_obj.Email)
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
	session := sessions.Default(context)
	user_id := session.Get("user_id")

	if user_id != nil {
		context.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"is_authenticated": true,
		})
		return
	}

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
				utils.ErrorCatcher(err, context, http.StatusInternalServerError, "signup.html", "failed to create new user!")
				return
			}

			user := models.User{
				Username:     username,
				PasswordHash: string(hashedPassword),
				Email:        email,
				Type:         "NoType",
				Score:        0,
			}

			err2 := database.DB.Create(&user).Error
			if err2 != nil {
				utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "signup.html", "failed to create new user!")
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

func LogoutHandler(context *gin.Context) {
	log.Println("LogoutHandler hit")

	session := sessions.Default(context)

	username := session.Get("username")
	if username == nil {
		context.HTML(http.StatusBadRequest, "home.html", gin.H{
			"message": "You are not logged in!",
		})
		return
	}

	session.Delete("user_id")
	session.Delete("username")
	session.Delete("email")
	session.Set("redirect", "not nil")
	session.Save()
	context.Redirect(http.StatusTemporaryRedirect, "/")
}

func UserslistHandler(context *gin.Context) {
	users_list := []models.User{}
	err := database.DB.Find(&users_list).Error
	if err != nil {
		utils.ErrorCatcher(err, context, http.StatusInternalServerError, "users_list.html", "failed to query db")
		return
	}

	context.HTML(http.StatusOK, "users_list.html", gin.H{
		"users_list": users_list,
		"message":    "user table fetched successfully!",
	})
}

func CreateTeamHandler(context *gin.Context) {
	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "create_team.html", nil)
	case "POST":
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
			Name:    name,
			Leader:  leader_obj,
			Members: []models.User{leader_obj},
		}
		err = database.DB.Create(&team).Error
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "create_team.html", "falied to create new team!")
			return
		}

		context.HTML(http.StatusOK, "create_team.html", gin.H{
			"message": "new team created!",
		})
	}

}

func DeleteTeamHandler(context *gin.Context) {
	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "delete_team.html", nil)
	case "POST":
		context.Request.ParseForm()
		name := context.Request.FormValue("Name")
		var team models.Team
		err := database.DB.Preload("Leader").Where(&models.Team{Name: name}).First(&team).Error
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "delete_team.html", "could not find the team!")
			return
		}

		session := sessions.Default(context)
		username := session.Get("username")

		log.Printf("username: %v", username)
		log.Printf("team: %v", team)
		log.Printf("team leader: %v", team.Leader.Username)

		if team.Leader.Username == username {
			err2 := database.DB.Where(&models.Team{Name: name}).Delete(&models.Team{}).Error
			if err2 != nil {
				utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "delete_team.html", "could not find the team!")
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
		utils.ErrorCatcher(err, context, http.StatusInternalServerError, "users_list.html", "failed to query db")
		return
	}

	context.HTML(http.StatusOK, "teams_list.html", gin.H{
		"teams_list": teams_list,
	})
}

func CreateTaskHandler(context *gin.Context) {
	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "create_task.html", nil)
	case "POST":
		name := context.Request.FormValue("Name")
		team_name := context.Request.FormValue("Team")
		deadline := context.Request.FormValue("Deadline")

		var team_obj models.Team
		err := database.DB.Where(&models.Team{Name: team_name}).First(&team_obj).Error

		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "users_list.html", "team not found!")
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
	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "delete_task.html", nil)
	case "POST":
		session := sessions.Default(context)
		username := session.Get("username")

		name := context.Request.FormValue("Name")
		team_name := context.Request.FormValue("Team")

		var team_obj models.Team
		err := database.DB.Preload("Leader").Where(&models.Team{Name: team_name}).First(&team_obj).Error
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "tasks_list.html", "team not found!")
			return
		}

		if team_obj.Leader.Username != username {
			context.HTML(http.StatusBadRequest, "delete_task.html", gin.H{
				"message": "Only team leader can perform this action!",
			})
			return
		}

		err2 := database.DB.Where(&models.Task{Name: name, Team: team_obj}).Delete(&models.Task{}).Error
		if err2 != nil {
			utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "tasks_list.html", "failed to delete the task!")
			return
		}

		context.HTML(http.StatusOK, "delete_task.html", gin.H{
			"message": "task has been deleted!",
		})
	}
}

func TasklistHandler(context *gin.Context) {
	session := sessions.Default(context)
	userID := session.Get("user_id")
	log.Printf("user id is:%v", userID)

	if userID == nil {
		context.HTML(http.StatusBadRequest, "tasks_list.html", gin.H{
			"message": "you are not logged in!",
		})
		return
	}

	user_teams := []models.Team{}
	err := database.DB.Preload("Members").Preload("Leader").Joins("JOIN team_users ON team_users.team_id=teams.id").Where("team_users.user_id=?", userID).Find(&user_teams).Error

	if err != nil {
		utils.ErrorCatcher(err, context, http.StatusInternalServerError, "tasks_list.html", "failed to query db")
		return
	}

	user_teams_ids := []uint{}
	for _, team := range user_teams {
		user_teams_ids = append(user_teams_ids, uint(team.ID))
	}

	tasks_list := []models.Task{}
	err2 := database.DB.Preload("Team").Preload("Team.Leader").Preload("Team.Members").Where("team_id IN ?", user_teams_ids).Find(&tasks_list).Error

	if err2 != nil {
		utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "tasks_list.html", "failed to find all user's tasks!")
		return
	}

	context.HTML(http.StatusOK, "tasks_list.html", gin.H{
		"tasks_list": tasks_list,
	})
}

// TODO: Complete UserDashboardHandler function later!
func UserDashboardHandler(context *gin.Context) {
	session := sessions.Default(context)
	username := session.Get("username")
	context.HTML(http.StatusOK, "dashboard.html", gin.H{
		"username": username,
	})
}

// TODO: Complete DeleteAccountHandler function later!
func DeleteAccountHandler(context *gin.Context) {
	session := sessions.Default(context)
	username := session.Get("username")
	context.HTML(http.StatusOK, "delete_account.html", gin.H{
		"username": username,
	})
}
