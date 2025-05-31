package controllers

import (
	"net/http"
	"strconv"
	"fithero-backend/models"
	"fithero-backend/services"
	"fithero-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService *services.UserService
	validator   *validator.Validate
}

// NewUserController creates a new user controller
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
		validator:   validator.New(),
	}
}

// CreateUser creates a new user (Admin only or self-registration)
func (uc *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := uc.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	user, err := uc.userService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUserByID retrieves a user by ID (with authorization check)
func (uc *UserController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Authorization check: Users can only view their own profile
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	if currentUserID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: You can only view your own profile",
		})
		return
	}

	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetCurrentUserProfile returns the current authenticated user's profile
func (uc *UserController) GetCurrentUserProfile(c *gin.Context) {
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUser updates a user's information (only the user themselves)
func (uc *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Authorization check: Users can only update their own profile
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	if currentUserID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: You can only update your own profile",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := uc.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := uc.userService.UpdateUser(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	// Return updated user data
	updatedUser, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}

// DeleteUser soft deletes a user (only the user themselves)
func (uc *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Authorization check: Users can only delete their own account
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	if currentUserID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: You can only delete your own account",
		})
		return
	}

	if err := uc.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetUserTasks returns the current user's daily tasks
func (uc *UserController) GetUserTasks(c *gin.Context) {
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	// Get user's daily tasks
	tasks, err := uc.userService.GetUserDailyTasks(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user tasks",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

// GetUserAchievements returns the current user's achievements
func (uc *UserController) GetUserAchievements(c *gin.Context) {
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	// Get user's achievements
	achievements, err := uc.userService.GetUserAchievements(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user achievements",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
	})
}

// GetLeaderboard handles GET /api/leaderboard
func (uc *UserController) GetLeaderboard(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	users, err := uc.userService.GetLeaderboard(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard"})
		return
	}

	c.JSON(http.StatusOK, users)
} 