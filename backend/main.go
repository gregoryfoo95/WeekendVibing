package main

import (
	"log"
	"os"

	"fithero-backend/config"
	"fithero-backend/controllers"
	"fithero-backend/repositories"
	"fithero-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDatabase()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.GetDB())
	taskRepo := repositories.NewTaskRepository(config.GetDB())
	achievementRepo := repositories.NewAchievementRepository(config.GetDB())

	// Initialize services
	userService := services.NewUserService(userRepo)
	taskService := services.NewTaskService(taskRepo, userService)
	achievementService := services.NewAchievementService(achievementRepo, userService)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	taskController := controllers.NewTaskController(taskService)
	achievementController := controllers.NewAchievementController(achievementService)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// API routes
	api := r.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.GET("", taskController.GetTasks)
			tasks.GET("/daily/:user_id", taskController.GetDailyTasks)
			tasks.POST("/daily", taskController.GenerateDailyTasks)
			tasks.PUT("/daily/:id/complete", taskController.CompleteTask)
		}

		// Achievement routes
		achievements := api.Group("/achievements")
		{
			achievements.GET("", achievementController.GetAchievements)
			achievements.GET("/user/:user_id", achievementController.GetUserAchievements)
			achievements.POST("/unlock", achievementController.UnlockAchievement)
		}

		// Leaderboard route
		api.GET("/leaderboard", userController.GetLeaderboard)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "fithero-backend",
		})
	})

	port := getEnv("PORT", "8080")
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Database connected with GORM")
	log.Printf("🔒 Using ORM for SQL injection prevention")
	log.Printf("✅ Layered architecture: Controller -> Service -> Repository")
	
	r.Run(":" + port)
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 