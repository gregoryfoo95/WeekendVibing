package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	GoogleID     string    `json:"google_id" gorm:"unique;index"`
	Email        string    `json:"email" gorm:"not null;unique" validate:"required,email"`
	Username     string    `json:"username" gorm:"not null;unique" validate:"required,min=3,max=50"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Picture      string    `json:"picture"`
	Level        int       `json:"level" gorm:"not null;default:1"`
	Points       int       `json:"points" gorm:"not null;default:0"`
	Character    string    `json:"character" gorm:"not null;default:'Rookie Hero'"`
	JobTitle     string    `json:"job_title" gorm:"not null;default:'Fitness Novice'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time `json:"last_login_at"`
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
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Level     *int    `json:"level,omitempty" validate:"omitempty,min=1,max=5"`
	Points    *int    `json:"points,omitempty" validate:"omitempty,min=0"`
	Character *string `json:"character,omitempty"`
	JobTitle  *string `json:"job_title,omitempty"`
}

// GoogleUserInfo represents the user info returned from Google OAuth
type GoogleUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
} 