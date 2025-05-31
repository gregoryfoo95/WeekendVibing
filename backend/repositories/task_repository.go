package repositories

import (
	"time"
	"fithero-backend/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	GetRandomTasks(count int) ([]models.Task, error)
	
	// Daily Tasks
	CreateDailyTask(dailyTask *models.DailyTask) error
	GetDailyTasksByUserAndDate(userID uint, date time.Time) ([]models.DailyTask, error)
	GetDailyTaskByID(id uint) (*models.DailyTask, error)
	UpdateDailyTask(dailyTask *models.DailyTask) error
	DeleteDailyTasksByUserAndDate(userID uint, date time.Time) error
}

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// GetAllTasks retrieves all tasks
func (r *taskRepository) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetTaskByID retrieves a task by ID
func (r *taskRepository) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetRandomTasks retrieves random tasks
func (r *taskRepository) GetRandomTasks(count int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Order("RANDOM()").Limit(count).Find(&tasks).Error
	return tasks, err
}

// CreateDailyTask creates a new daily task
func (r *taskRepository) CreateDailyTask(dailyTask *models.DailyTask) error {
	return r.db.Create(dailyTask).Error
}

// GetDailyTasksByUserAndDate retrieves daily tasks for a user on a specific date
func (r *taskRepository) GetDailyTasksByUserAndDate(userID uint, date time.Time) ([]models.DailyTask, error) {
	var dailyTasks []models.DailyTask
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	err := r.db.Preload("Task").
		Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).
		Find(&dailyTasks).Error
	return dailyTasks, err
}

// GetDailyTaskByID retrieves a daily task by ID
func (r *taskRepository) GetDailyTaskByID(id uint) (*models.DailyTask, error) {
	var dailyTask models.DailyTask
	err := r.db.Preload("Task").First(&dailyTask, id).Error
	if err != nil {
		return nil, err
	}
	return &dailyTask, nil
}

// UpdateDailyTask updates a daily task
func (r *taskRepository) UpdateDailyTask(dailyTask *models.DailyTask) error {
	return r.db.Save(dailyTask).Error
}

// DeleteDailyTasksByUserAndDate deletes daily tasks for a user on a specific date
func (r *taskRepository) DeleteDailyTasksByUserAndDate(userID uint, date time.Time) error {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	return r.db.Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).
		Delete(&models.DailyTask{}).Error
} 