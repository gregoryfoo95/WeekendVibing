package controllers

import (
	"net/http"
	"strconv"
	"fithero-backend/services"
	"fithero-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskController struct {
	taskService *services.TaskService
	validator   *validator.Validate
}

// NewTaskController creates a new task controller
func NewTaskController(taskService *services.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
		validator:   validator.New(),
	}
}

// GetAllTasks handles GET /api/public/tasks
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GenerateDailyTasks handles POST /api/tasks/daily/generate
func (tc *TaskController) GenerateDailyTasks(c *gin.Context) {
	// Get current user ID from middleware
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	dailyTasks, err := tc.taskService.GenerateDailyTasks(userID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "no tasks available for user level":
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No tasks available for your level"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate daily tasks"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Daily tasks generated successfully",
		"tasks":   dailyTasks,
	})
}

// CompleteTask handles POST /api/tasks/daily/:id/complete
func (tc *TaskController) CompleteTask(c *gin.Context) {
	// Get current user ID from middleware
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	idParam := c.Param("id")
	dailyTaskID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	dailyTask, err := tc.taskService.CompleteTask(userID, uint(dailyTaskID))
	if err != nil {
		switch err.Error() {
		case "daily task not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Daily task not found"})
		case "task already completed":
			c.JSON(http.StatusConflict, gin.H{"error": "Task already completed"})
		case "access denied: you can only complete your own tasks":
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only complete your own tasks"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Task completed successfully",
		"daily_task":    dailyTask,
		"points_earned": dailyTask.Points,
	})
} 