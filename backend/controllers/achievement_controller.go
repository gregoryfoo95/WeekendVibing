package controllers

import (
	"net/http"
	"strconv"
	"fithero-backend/services"
	"fithero-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AchievementController struct {
	achievementService *services.AchievementService
	validator         *validator.Validate
}

// NewAchievementController creates a new achievement controller
func NewAchievementController(achievementService *services.AchievementService) *AchievementController {
	return &AchievementController{
		achievementService: achievementService,
		validator:         validator.New(),
	}
}

// GetAllAchievements handles GET /api/achievements and GET /api/public/achievements
func (ac *AchievementController) GetAllAchievements(c *gin.Context) {
	achievements, err := ac.achievementService.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// GetUserAchievements handles GET /api/achievements/user (for current user)
func (ac *AchievementController) GetUserAchievements(c *gin.Context) {
	// Get current user ID from middleware
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userAchievements, err := ac.achievementService.GetUserAchievements(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user achievements"})
		return
	}

	c.JSON(http.StatusOK, userAchievements)
}

// UnlockAchievement handles POST /api/achievements/:id/unlock
func (ac *AchievementController) UnlockAchievement(c *gin.Context) {
	// Get current user ID from middleware
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get achievement ID from URL parameter
	achievementIDParam := c.Param("id")
	achievementID, err := strconv.ParseUint(achievementIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
		return
	}

	userAchievement, err := ac.achievementService.UnlockAchievement(userID, uint(achievementID))
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "achievement not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		case "achievement already unlocked":
			c.JSON(http.StatusConflict, gin.H{"error": "Achievement already unlocked"})
		case "insufficient points to unlock achievement":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient points to unlock achievement"})
		case "failed to deduct points from user", "failed to unlock achievement":
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlock achievement"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Achievement unlocked successfully",
		"user_achievement": userAchievement,
	})
} 