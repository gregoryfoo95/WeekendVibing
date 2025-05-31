package services

import (
	"errors"
	"time"
	"fithero-backend/models"
	"fithero-backend/repositories"
	"gorm.io/gorm"
)

type AchievementService struct {
	achievementRepo repositories.AchievementRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
}

// NewAchievementService creates a new achievement service
func NewAchievementService(achievementRepo repositories.AchievementRepositoryInterface, userRepo repositories.UserRepositoryInterface) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
	}
}

// GetAllAchievements returns all available achievements
func (s *AchievementService) GetAllAchievements() ([]models.Achievement, error) {
	return s.achievementRepo.GetAll()
}

// GetUserAchievements retrieves all achievements unlocked by a user
func (s *AchievementService) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.achievementRepo.GetUserAchievements(userID)
}

// UnlockAchievement unlocks an achievement for a user with business logic validation
func (s *AchievementService) UnlockAchievement(userID, achievementID uint) (*models.UserAchievement, error) {
	// Validate user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Validate achievement exists
	achievement, err := s.achievementRepo.GetByID(achievementID)
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
	newPoints := user.Points - achievement.PointsCost
	err = s.userRepo.Update(userID, &models.UpdateUserRequest{
		Points: &newPoints,
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

	createdAchievement, err := s.achievementRepo.CreateUserAchievement(userAchievement)
	if err != nil {
		// Rollback points if achievement creation fails
		rollbackPoints := user.Points
		s.userRepo.Update(userID, &models.UpdateUserRequest{
			Points: &rollbackPoints,
		})
		return nil, errors.New("failed to unlock achievement")
	}

	// Update user job title or character if achievement affects them
	s.updateUserBasedOnAchievement(userID, achievement)

	return createdAchievement, nil
}

// updateUserBasedOnAchievement updates user character or job title based on achievement type
func (s *AchievementService) updateUserBasedOnAchievement(userID uint, achievement *models.Achievement) {
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
		s.userRepo.Update(userID, updateReq)
	}
}

// mapAchievementToJobTitle maps achievement titles to job titles
func (s *AchievementService) mapAchievementToJobTitle(achievementTitle string) string {
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