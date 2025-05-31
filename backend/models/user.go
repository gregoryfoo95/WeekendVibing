package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"not null;unique" validate:"required,min=3,max=50"`
	Email        string    `json:"email" gorm:"not null;unique" validate:"required,email"`
	Level        int       `json:"level" gorm:"not null;default:1"`
	Points       int       `json:"points" gorm:"not null;default:0"`
	Character    string    `json:"character" gorm:"not null;default:'Rookie Hero'"`
	JobTitle     string    `json:"job_title" gorm:"not null;default:'Fitness Novice'"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	DailyTasks      []DailyTask      `json:"daily_tasks,omitempty" gorm:"foreignKey:UserID"`
	UserAchievements []UserAchievement `json:"user_achievements,omitempty" gorm:"foreignKey:UserID"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Level     *int    `json:"level,omitempty" validate:"omitempty,min=1,max=5"`
	Points    *int    `json:"points,omitempty" validate:"omitempty,min=0"`
	Character *string `json:"character,omitempty"`
	JobTitle  *string `json:"job_title,omitempty"`
} 