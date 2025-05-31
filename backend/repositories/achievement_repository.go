package repositories

import (
	"fithero-backend/models"
	"gorm.io/gorm"
)

type AchievementRepository interface {
	GetAllAchievements() ([]models.Achievement, error)
	GetAchievementByID(id uint) (*models.Achievement, error)
	
	// User Achievements
	CreateUserAchievement(userAchievement *models.UserAchievement) error
	GetUserAchievements(userID uint) ([]models.UserAchievement, error)
	GetUserAchievementByUserAndAchievement(userID, achievementID uint) (*models.UserAchievement, error)
	IsAchievementUnlocked(userID, achievementID uint) (bool, error)
}

type achievementRepository struct {
	db *gorm.DB
}

// NewAchievementRepository creates a new achievement repository
func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

// GetAllAchievements retrieves all achievements
func (r *achievementRepository) GetAllAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Find(&achievements).Error
	return achievements, err
}

// GetAchievementByID retrieves an achievement by ID
func (r *achievementRepository) GetAchievementByID(id uint) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.First(&achievement, id).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// CreateUserAchievement creates a new user achievement
func (r *achievementRepository) CreateUserAchievement(userAchievement *models.UserAchievement) error {
	return r.db.Create(userAchievement).Error
}

// GetUserAchievements retrieves all achievements for a user
func (r *achievementRepository) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	var userAchievements []models.UserAchievement
	err := r.db.Preload("Achievement").
		Where("user_id = ?", userID).
		Find(&userAchievements).Error
	return userAchievements, err
}

// GetUserAchievementByUserAndAchievement retrieves a specific user achievement
func (r *achievementRepository) GetUserAchievementByUserAndAchievement(userID, achievementID uint) (*models.UserAchievement, error) {
	var userAchievement models.UserAchievement
	err := r.db.Preload("Achievement").
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		First(&userAchievement).Error
	if err != nil {
		return nil, err
	}
	return &userAchievement, nil
}

// IsAchievementUnlocked checks if a user has unlocked a specific achievement
func (r *achievementRepository) IsAchievementUnlocked(userID, achievementID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error
	return count > 0, err
} 