package models

import (
	"time"
	"gorm.io/gorm"
)

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Points      int    `json:"points" gorm:"not null"`
	Category    string `json:"category" gorm:"not null"` // cardio, strength, flexibility, wellness
	Difficulty  string `json:"difficulty" gorm:"not null"` // easy, medium, hard
	Level       int    `json:"level" gorm:"not null;default:1"` // Required user level
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	DailyTasks []DailyTask `json:"daily_tasks,omitempty" gorm:"foreignKey:TaskID"`
}

type DailyTask struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	TaskID      uint      `json:"task_id" gorm:"not null;index"`
	IsCompleted bool      `json:"is_completed" gorm:"not null;default:false"`
	Points      int       `json:"points" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Task Task `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}

// UpdateDailyTaskRequest represents the request to update a daily task
type UpdateDailyTaskRequest struct {
	IsCompleted *bool `json:"is_completed,omitempty"`
}

// GenerateDailyTasksRequest represents the request to generate daily tasks
type GenerateDailyTasksRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

// CompleteTaskRequest represents the request to complete a task
type CompleteTaskRequest struct {
	TaskID uint `json:"task_id" validate:"required"`
} 