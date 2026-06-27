package main

import (
	"fmt"

	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = database.InitDB()
	router := gin.Default()

	router.GET("/", handlers.HomepageHandler)
	router.Any("/login", handlers.LoginpageHandler)
	router.Any("/signup", handlers.SignuppageHandler)
	router.GET("/users_list", handlers.UserslistHandler)
	router.GET("/teams_list", handlers.TeamslistHandler)
	router.GET("/tasks_list", handlers.TasklistHandler)
	router.GET("/create_team", handlers.CreateTeamHandler)
	router.POST("/create_team", handlers.CreateTeamHandler)
	router.GET("/delete_team", handlers.DeleteTeamHandler)
	router.POST("/delete_team", handlers.DeleteTeamHandler)
	router.GET("/create_task", handlers.CreateTaskHandler)
	router.POST("/create_task", handlers.CreateTaskHandler)
	router.GET("/delete_task", handlers.DeleteTaskHandler)
	router.POST("/delete_task", handlers.DeleteTaskHandler)

	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("/statics", "./web/statics")

	fmt.Println("Server started on http://localhost:8080")

	router.Run("0.0.0.0:8080")
}
