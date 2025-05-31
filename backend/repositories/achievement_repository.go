package repositories

import (
	"fithero-backend/models"
	"gorm.io/gorm"
)

type AchievementRepositoryInterface interface {
	GetAll() ([]models.Achievement, error)
	GetByID(id uint) (*models.Achievement, error)
	
	// User Achievements
	CreateUserAchievement(userAchievement *models.UserAchievement) (*models.UserAchievement, error)
	GetUserAchievements(userID uint) ([]models.UserAchievement, error)
	GetUserAchievementByUserAndAchievement(userID, achievementID uint) (*models.UserAchievement, error)
	IsAchievementUnlocked(userID, achievementID uint) (bool, error)
}

type AchievementRepository struct {
	db *gorm.DB
}

// NewAchievementRepository creates a new achievement repository
func NewAchievementRepository(db *gorm.DB) AchievementRepositoryInterface {
	return &AchievementRepository{db: db}
}

// GetAll retrieves all achievements
func (r *AchievementRepository) GetAll() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Find(&achievements).Error
	return achievements, err
}

// GetByID retrieves an achievement by ID
func (r *AchievementRepository) GetByID(id uint) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.First(&achievement, id).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// CreateUserAchievement creates a new user achievement
func (r *AchievementRepository) CreateUserAchievement(userAchievement *models.UserAchievement) (*models.UserAchievement, error) {
	if err := r.db.Create(userAchievement).Error; err != nil {
		return nil, err
	}
	// Load the achievement relationship
	if err := r.db.Preload("Achievement").First(userAchievement, userAchievement.ID).Error; err != nil {
		return nil, err
	}
	return userAchievement, nil
}

// GetUserAchievements retrieves all achievements for a user
func (r *AchievementRepository) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	var userAchievements []models.UserAchievement
	err := r.db.Preload("Achievement").
		Where("user_id = ?", userID).
		Find(&userAchievements).Error
	return userAchievements, err
}

// GetUserAchievementByUserAndAchievement retrieves a specific user achievement
func (r *AchievementRepository) GetUserAchievementByUserAndAchievement(userID, achievementID uint) (*models.UserAchievement, error) {
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
func (r *AchievementRepository) IsAchievementUnlocked(userID, achievementID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error
	return count > 0, err
} 