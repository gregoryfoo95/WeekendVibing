package models

import (
	"time"
	"gorm.io/gorm"
)

type Achievement struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Icon        string `json:"icon" gorm:"not null"`
	PointsCost  int    `json:"points_cost" gorm:"not null"`
	Type        string `json:"type" gorm:"not null"` // character, upgrade, badge
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	UserAchievements []UserAchievement `json:"user_achievements,omitempty" gorm:"foreignKey:AchievementID"`
}

type UserAchievement struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        uint        `json:"user_id" gorm:"not null;index"`
	AchievementID uint        `json:"achievement_id" gorm:"not null;index"`
	UnlockedAt    time.Time   `json:"unlocked_at"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Achievement Achievement `json:"achievement,omitempty" gorm:"foreignKey:AchievementID"`
	
	// Composite unique index
	// This ensures a user can only unlock each achievement once
	_ struct{} `gorm:"uniqueIndex:idx_user_achievement,composite:user_id,achievement_id"`
}

// UnlockAchievementRequest represents the request to unlock an achievement
type UnlockAchievementRequest struct {
	UserID        uint `json:"user_id" validate:"required"`
	AchievementID uint `json:"achievement_id" validate:"required"`
} 