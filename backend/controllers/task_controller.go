package controllers

import (
	"net/http"
	"strconv"
	"time"
	"fithero-backend/models"
	"fithero-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskController struct {
	taskService services.TaskService
	validator   *validator.Validate
}

// NewTaskController creates a new task controller
func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
		validator:   validator.New(),
	}
}

// GetTasks handles GET /api/tasks
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetDailyTasks handles GET /api/tasks/daily/:user_id
func (tc *TaskController) GetDailyTasks(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse date from query parameter, default to today
	dateParam := c.DefaultQuery("date", "")
	var date time.Time
	if dateParam != "" {
		date, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		date = time.Now().Truncate(24 * time.Hour)
	}

	dailyTasks, err := tc.taskService.GetDailyTasks(uint(userID), date)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve daily tasks"})
		return
	}

	c.JSON(http.StatusOK, dailyTasks)
}

// GenerateDailyTasks handles POST /api/tasks/daily
func (tc *TaskController) GenerateDailyTasks(c *gin.Context) {
	var req models.GenerateDailyTasksRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	if err := tc.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// If no date provided, use today
	if req.Date.IsZero() {
		req.Date = time.Now().Truncate(24 * time.Hour)
	}

	dailyTasks, err := tc.taskService.GenerateDailyTasks(req.UserID, req.Date)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "no tasks available":
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No tasks available to generate"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate daily tasks"})
		}
		return
	}

	c.JSON(http.StatusCreated, dailyTasks)
}

// CompleteTask handles PUT /api/tasks/daily/:id/complete
func (tc *TaskController) CompleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	dailyTask, err := tc.taskService.CompleteTask(uint(id))
	if err != nil {
		switch err.Error() {
		case "daily task not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Daily task not found"})
		case "task already completed":
			c.JSON(http.StatusConflict, gin.H{"error": "Task already completed"})
		case "can only complete tasks for today":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can only complete tasks for today"})
		case "failed to award points to user":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to award points"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task completed successfully",
		"daily_task": dailyTask,
		"points_earned": dailyTask.Task.Points,
	})
} 