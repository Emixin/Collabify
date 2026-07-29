package handlers

import (
	"fmt"
	"strconv"

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

	var redirected any
	redirected = session.Get("redirect")
	redirected_str, actually_redirected := redirected.(string)

	if actually_redirected && redirected_str != "" {
		session.Delete("redirect")

		if err := session.Save(); err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "home.html", "failed to save session!")
			return
		}

		context.HTML(http.StatusOK, "home.html", gin.H{
			"username":         "Anonymous User",
			"redirect_message": "Logged out successfully!",
			"is_authenticated": false,
		})
		return
	} else {

		user_teams_ids, err := utils.UserTeamIDs(session, context, "home.html")
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "home.html", err.Error())
			return
		}

		// Note: Tried to use one query instead of two queries!
		all_tasks := database.DB.Model(&models.Task{}).Where("team_id IN ?", user_teams_ids)
		var all_tasks_count int64
		all_tasks.Count(&all_tasks_count)

		all_tasks.Where("status = ?", models.StatusPending)
		var pending_tasks_count int64
		all_tasks.Count(&pending_tasks_count)

		log.Printf("all tasks: %v\npending tasks: %v", all_tasks_count, pending_tasks_count)

		message := fmt.Sprintf("You have %d tasks and %d of them are pending", all_tasks_count, pending_tasks_count)

		if username == nil {
			context.HTML(http.StatusOK, "home.html", gin.H{
				"username":         "Anonymous User",
				"is_authenticated": false,
			})
		} else {
			context.HTML(http.StatusOK, "home.html", gin.H{
				"message":          message,
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

	username, ok := session.Get("username").(string)
	log.Printf("username:%v", username)

	if !ok {
		utils.ErrorCatcher(nil, context, http.StatusOK, "home.html", "failed to fetch your username!")
		return
	}

	if username == "" {
		context.HTML(http.StatusBadRequest, "home.html", gin.H{
			"message": "You are not logged in!",
		})
		return
	}

	session.Delete("user_id")
	session.Delete("username")
	session.Delete("email")
	session.Set("redirect", "not nil")

	if err := session.Save(); err != nil {
		utils.ErrorCatcher(err, context, http.StatusInternalServerError, "home.html", "logout failed!")
		return
	}

	context.Redirect(http.StatusFound, "/")
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
			Status:   models.StatusPending,
		}
		err2 := database.DB.Create(&task).Error
		if err2 != nil {
			utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "users_list.html", "team not created!")
			return
		}

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
		username, ok := session.Get("username").(string)
		if !ok {
			utils.ErrorCatcher(nil, context, http.StatusInternalServerError, "tasks_list.html", "username is not correct!")
			return
		}

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
	user_id, ok := session.Get("user_id").(int)

	if !ok {
		utils.ErrorCatcher(nil, context, http.StatusInternalServerError, "tasks_list.html", "failed to fetch user id")
		return
	}

	42switch context.Request.Method {
	case "GET":
		user_teams_ids, err := utils.UserTeamIDs(session, context, "tasks_list.html")
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "tasks_list", err.Error())
			return
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

	case "POST":
		task_id_str := context.Request.FormValue("task_id")
		task_id_int, err := strconv.Atoi(task_id_str)
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusInternalServerError, "tasks_list.html", "failed to find task")
			return
		}

		var task_obj models.Task
		err2 := database.DB.Preload("Team").Where(&models.Task{ID: task_id_int}).First(&task_obj).Error
		if err2 != nil {
			utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "tasks_list.html", "failed to change task status")
			return
		}

		leader_id := task_obj.Team.LeaderID

		changed := task_obj.MarkAsCompleted(user_id, leader_id)
		if changed {
			// err3 := database.DB.Where(&models.Task{ID: task_id_int}).Update("Status", models.StatusCompleted).Error
			// if err3 != nil {
			// 	utils.ErrorCatcher(err3, context, http.StatusInternalServerError, "tasks_list.html", "failed to change task status")
			// 	return
			// }
			database.DB.Save(&task_obj)
			context.HTML(http.StatusOK, "tasks_list.html", gin.H{
				"message": "task marked as completed!",
			})
			return
		} else if leader_id != user_id {
			context.HTML(http.StatusBadRequest, "tasks_list.html", gin.H{
				"message": "only team leader can perform this action!",
			})
			return
		}

		context.HTML(http.StatusBadRequest, "tasks_list.html", gin.H{
			"message": "task was already marked as completed!",
		})
	}
}

// TODO: Complete UserDashboardHandler function later!
func UserDashboardHandler(context *gin.Context) {
	session := sessions.Default(context)
	username := session.Get("username")

	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "dashboard.html", gin.H{
			"username": username,
		})
	}
}

func DeleteAccountHandler(context *gin.Context) {
	session := sessions.Default(context)
	username := session.Get("username")

	username_string, ok := username.(string)
	if !ok {
		utils.ErrorCatcher(nil, context, http.StatusBadRequest, "delete_account.html", "you must be logged in to delete your account!")
		return
	}

	switch context.Request.Method {
	case "GET":
		context.HTML(http.StatusOK, "delete_account.html", gin.H{
			"username": username,
		})
	case "POST":
		err := context.Request.ParseForm()
		if err != nil {
			utils.ErrorCatcher(err, context, http.StatusBadRequest, "delete_account.html", "confirmation failed!")
		}
		confirmation := context.Request.FormValue("confirmation")
		if confirmation != "" {
			if confirmation == "i want to delete my account" {
				err2 := database.DB.Where(&models.User{Username: username_string}).Delete(&models.User{}).Error
				if err2 != nil {
					utils.ErrorCatcher(err2, context, http.StatusInternalServerError, "delete_account.html", "could not find the user!")
					return
				}

				session.Delete("user_id")
				session.Delete("username")
				session.Delete("email")
				session.Save()
				context.HTML(http.StatusOK, "delete_account.html", gin.H{
					"message": "user has been deleted!",
				})
			} else {
				context.HTML(http.StatusBadRequest, "delete_account.html", gin.H{
					"message": "confirmation failed!",
				})
			}
		} else {
			context.HTML(http.StatusBadRequest, "delete_account.html", gin.H{
				"message": "empty confirmation field!",
			})
		}
	}
}
