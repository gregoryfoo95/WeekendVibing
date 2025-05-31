package main

import (
	"log"
	"os"

	"fithero-backend/config"
	"fithero-backend/controllers"
	"fithero-backend/middleware"
	"fithero-backend/repositories"
	"fithero-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize authentication configuration
	authConfig := config.NewAuthConfig()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	achievementRepo := repositories.NewAchievementRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, authConfig)
	userService := services.NewUserService(userRepo, taskRepo, achievementRepo)
	taskService := services.NewTaskService(taskRepo, userRepo, achievementRepo)
	achievementService := services.NewAchievementService(achievementRepo, userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, authConfig)
	userController := controllers.NewUserController(userService)
	taskController := controllers.NewTaskController(taskService)
	achievementController := controllers.NewAchievementController(achievementService)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	corsConfig.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
	}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "fithero-backend",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.GET("/google", authController.GoogleLogin)
			auth.GET("/google/callback", authController.GoogleCallback)
			auth.POST("/logout", authController.Logout)
			auth.GET("/check", middleware.OptionalAuthMiddleware(authService), authController.CheckAuth)
		}

		// Protected routes requiring authentication
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User profile routes
			protected.GET("/me", authController.Me)
			protected.POST("/auth/refresh", authController.RefreshToken)
			protected.GET("/profile", userController.GetCurrentUserProfile)
			protected.PUT("/profile", func(c *gin.Context) {
				// Get current user ID and forward to update method
				userID, exists := middleware.GetCurrentUserID(c)
				if !exists {
					c.JSON(401, gin.H{"error": "Authentication required"})
					return
				}
				c.Param("id")
				c.Set("id", userID)
				userController.UpdateUser(c)
			})

			// User-specific routes with ownership verification
			users := protected.Group("/users")
			{
				users.GET("/:id", userController.GetUserByID)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)
				users.GET("/tasks", userController.GetUserTasks)
				users.GET("/achievements", userController.GetUserAchievements)
			}

			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.GET("/daily", func(c *gin.Context) {
					userID, _ := middleware.GetCurrentUserID(c)
					dailyTasks, err := taskService.GetUserDailyTasks(userID)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"tasks": dailyTasks})
				})
				tasks.POST("/daily/generate", func(c *gin.Context) {
					userID, _ := middleware.GetCurrentUserID(c)
					dailyTasks, err := taskService.GenerateDailyTasks(userID)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"tasks": dailyTasks})
				})
				tasks.POST("/daily/:id/complete", func(c *gin.Context) {
					taskController.CompleteTask(c)
				})
			}

			// Achievement routes
			achievements := protected.Group("/achievements")
			{
				achievements.GET("/", achievementController.GetAllAchievements)
				achievements.POST("/:id/unlock", achievementController.UnlockAchievement)
				achievements.GET("/user", func(c *gin.Context) {
					userID, _ := middleware.GetCurrentUserID(c)
					userAchievements, err := achievementService.GetUserAchievements(userID)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"achievements": userAchievements})
				})
			}
		}

		// Public routes (no authentication required)
		public := api.Group("/public")
		{
			public.GET("/tasks", taskController.GetAllTasks)
			public.GET("/achievements", achievementController.GetAllAchievements)
			public.GET("/leaderboard", userController.GetLeaderboard)
		}

		// Admin routes (for user creation - could be expanded)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		{
			admin.POST("/users", userController.CreateUser)
		}
	}

	// Set port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 