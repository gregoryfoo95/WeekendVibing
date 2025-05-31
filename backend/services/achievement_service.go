package services

import (
	"errors"
	"time"
	"fithero-backend/models"
	"fithero-backend/repositories"
	"gorm.io/gorm"
)

type AchievementService interface {
	GetAllAchievements() ([]models.Achievement, error)
	GetUserAchievements(userID uint) ([]models.UserAchievement, error)
	UnlockAchievement(userID, achievementID uint) (*models.UserAchievement, error)
}

type achievementService struct {
	achievementRepo repositories.AchievementRepository
	userService     UserService
}

// NewAchievementService creates a new achievement service
func NewAchievementService(achievementRepo repositories.AchievementRepository, userService UserService) AchievementService {
	return &achievementService{
		achievementRepo: achievementRepo,
		userService:     userService,
	}
}

// GetAllAchievements returns all available achievements
func (s *achievementService) GetAllAchievements() ([]models.Achievement, error) {
	return s.achievementRepo.GetAllAchievements()
}

// GetUserAchievements retrieves all achievements unlocked by a user
func (s *achievementService) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	// Validate user exists
	_, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.achievementRepo.GetUserAchievements(userID)
}

// UnlockAchievement unlocks an achievement for a user with business logic validation
func (s *achievementService) UnlockAchievement(userID, achievementID uint) (*models.UserAchievement, error) {
	// Validate user exists
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Validate achievement exists
	achievement, err := s.achievementRepo.GetAchievementByID(achievementID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("achievement not found")
		}
		return nil, err
	}

	// Check if achievement is already unlocked
	isUnlocked, err := s.achievementRepo.IsAchievementUnlocked(userID, achievementID)
	if err != nil {
		return nil, err
	}
	if isUnlocked {
		return nil, errors.New("achievement already unlocked")
	}

	// Check if user has enough points
	if user.Points < achievement.PointsCost {
		return nil, errors.New("insufficient points to unlock achievement")
	}

	// Deduct points from user
	user.Points -= achievement.PointsCost
	_, err = s.userService.UpdateUser(userID, &models.UpdateUserRequest{
		Points: &user.Points,
	})
	if err != nil {
		return nil, errors.New("failed to deduct points from user")
	}

	// Create user achievement
	userAchievement := &models.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
		UnlockedAt:    time.Now(),
	}

	err = s.achievementRepo.CreateUserAchievement(userAchievement)
	if err != nil {
		// Rollback points if achievement creation fails
		user.Points += achievement.PointsCost
		s.userService.UpdateUser(userID, &models.UpdateUserRequest{
			Points: &user.Points,
		})
		return nil, errors.New("failed to unlock achievement")
	}

	// Load the achievement relationship
	userAchievement.Achievement = *achievement

	// Update user job title or character if achievement affects them
	s.updateUserBasedOnAchievement(userID, achievement)

	return userAchievement, nil
}

// updateUserBasedOnAchievement updates user character or job title based on achievement type
func (s *achievementService) updateUserBasedOnAchievement(userID uint, achievement *models.Achievement) {
	updateReq := &models.UpdateUserRequest{}
	
	switch achievement.Type {
	case "character":
		updateReq.Character = &achievement.Title
	case "upgrade":
		// Map achievement titles to job titles
		jobTitle := s.mapAchievementToJobTitle(achievement.Title)
		if jobTitle != "" {
			updateReq.JobTitle = &jobTitle
		}
	}

	// Only update if there are changes
	if updateReq.Character != nil || updateReq.JobTitle != nil {
		s.userService.UpdateUser(userID, updateReq)
	}
}

// mapAchievementToJobTitle maps achievement titles to job titles
func (s *achievementService) mapAchievementToJobTitle(achievementTitle string) string {
	jobTitleMap := map[string]string{
		"Personal Trainer":   "Personal Trainer",
		"Fitness Coach":      "Fitness Coach", 
		"Wellness Expert":    "Wellness Expert",
		"Fitness Director":   "Fitness Director",
		"Health Guru":        "Health Guru",
	}
	
	if jobTitle, exists := jobTitleMap[achievementTitle]; exists {
		return jobTitle
	}
	return ""
} 