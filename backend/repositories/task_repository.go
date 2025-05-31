package repositories

import (
	"fithero-backend/models"
	"gorm.io/gorm"
)

type TaskRepositoryInterface interface {
	GetAll() ([]models.Task, error)
	GetByID(id uint) (*models.Task, error)
	GetTasksByLevel(level int) ([]models.Task, error)
	
	// Daily Tasks
	CreateDailyTask(dailyTask *models.DailyTask) (*models.DailyTask, error)
	GetDailyTasksByUserID(userID uint) ([]models.DailyTask, error)
	GetDailyTaskByID(id uint) (*models.DailyTask, error)
	UpdateDailyTask(id uint, updates *models.UpdateDailyTaskRequest) error
}

type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) TaskRepositoryInterface {
	return &TaskRepository{db: db}
}

// GetAll retrieves all tasks
func (r *TaskRepository) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetTasksByLevel retrieves tasks suitable for a specific user level
func (r *TaskRepository) GetTasksByLevel(level int) ([]models.Task, error) {
	var tasks []models.Task
	// Get tasks where required level is less than or equal to user level
	err := r.db.Where("level <= ?", level).Find(&tasks).Error
	return tasks, err
}

// CreateDailyTask creates a new daily task
func (r *TaskRepository) CreateDailyTask(dailyTask *models.DailyTask) (*models.DailyTask, error) {
	if err := r.db.Create(dailyTask).Error; err != nil {
		return nil, err
	}
	// Load the task relationship
	if err := r.db.Preload("Task").First(dailyTask, dailyTask.ID).Error; err != nil {
		return nil, err
	}
	return dailyTask, nil
}

// GetDailyTasksByUserID retrieves daily tasks for a user
func (r *TaskRepository) GetDailyTasksByUserID(userID uint) ([]models.DailyTask, error) {
	var dailyTasks []models.DailyTask
	err := r.db.Preload("Task").
		Where("user_id = ?", userID).
		Find(&dailyTasks).Error
	return dailyTasks, err
}

// GetDailyTaskByID retrieves a daily task by ID
func (r *TaskRepository) GetDailyTaskByID(id uint) (*models.DailyTask, error) {
	var dailyTask models.DailyTask
	err := r.db.Preload("Task").First(&dailyTask, id).Error
	if err != nil {
		return nil, err
	}
	return &dailyTask, nil
}

// UpdateDailyTask updates a daily task
func (r *TaskRepository) UpdateDailyTask(id uint, updates *models.UpdateDailyTaskRequest) error {
	dailyTask := &models.DailyTask{}
	if err := r.db.First(dailyTask, id).Error; err != nil {
		return err
	}

	updateData := make(map[string]interface{})
	
	if updates.IsCompleted != nil {
		updateData["is_completed"] = *updates.IsCompleted
	}

	if len(updateData) > 0 {
		return r.db.Model(dailyTask).Updates(updateData).Error
	}
	
	return nil
} 