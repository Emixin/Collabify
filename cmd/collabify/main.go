package main

import (
	"collabify/internal/database"
	"collabify/internal/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = database.InitDB()
	router := gin.Default()

	router.GET("/", handlers.HomepageHandler)
	router.GET("/login", handlers.LoginpageHandler)
	router.POST("login", handlers.LoginpageHandler)
	router.GET("/signup", handlers.SignuppageHandler)
	router.POST("signup", handlers.SignuppageHandler)
	router.GET("/users_list", handlers.UserslistHandler)
	router.GET("/teams_list", handlers.TeamslistHandler)
	router.GET("/tasks_list", handlers.TasklistHandler)
	router.GET("/create_team", handlers.CreateTeamHandler)
	router.POST("/create_team", handlers.CreateTeamHandler)

	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("/statics", "./web/statics")

	fmt.Println("Server started on http://localhost:8080")

	router.Run("0.0.0.0:8080")
}
