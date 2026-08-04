package main

import (
	"fmt"
	"os"

	"github.com/Emixin/Collabify/internal/database"
	"github.com/Emixin/Collabify/internal/handlers"
	"github.com/Emixin/Collabify/internal/middlewares"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	_ = database.InitDB()

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env files")
	}

	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		panic("SESSION_SECRET is not set")
	}

	// Default router
	router := gin.Default()
	store := cookie.NewStore([]byte(secret))
	router.Use(sessions.Sessions("collabify-session", store))

	router.GET("/", handlers.HomepageHandler)
	router.GET("/login", handlers.LoginpageHandler)
	router.GET("/users_list", handlers.UserslistHandler)
	router.GET("/teams_list", handlers.TeamslistHandler)
	router.POST("/login", handlers.LoginpageHandler)
	router.Any("/signup", handlers.SignuppageHandler)

	// Protected Group
	protected := router.Group("/")
	protected.Use(middlewares.RequireAuthMiddleware())

	protected.GET("/tasks_list", handlers.TasklistHandler)
	protected.GET("/delete_team", handlers.DeleteTeamHandler)
	protected.GET("/delete_task", handlers.DeleteTaskHandler)
	protected.GET("/create_task", handlers.CreateTaskHandler)
	protected.GET("/create_team", handlers.CreateTeamHandler)
	protected.POST("/tasks_list", handlers.TasklistHandler)
	protected.POST("/create_team", handlers.CreateTeamHandler)
	protected.POST("/delete_team", handlers.DeleteTeamHandler)
	protected.POST("/create_task", handlers.CreateTaskHandler)
	protected.POST("/delete_task", handlers.DeleteTaskHandler)
	protected.Any("/logout", handlers.LogoutHandler)
	protected.Any("/dashboard", handlers.UserDashboardHandler)
	protected.Any("/delete_account", handlers.DeleteAccountHandler)

	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("/statics", "./web/statics")

	fmt.Println("Server started on http://localhost:8080")
	router.Run("0.0.0.0:8080")
}
