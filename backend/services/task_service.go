package services

import (
	"errors"
	"time"
	"fithero-backend/models"
	"fithero-backend/repositories"
	"gorm.io/gorm"
)

type TaskService interface {
	GetAllTasks() ([]models.Task, error)
	GetDailyTasks(userID uint, date time.Time) ([]models.DailyTask, error)
	GenerateDailyTasks(userID uint, date time.Time) ([]models.DailyTask, error)
	CompleteTask(taskID uint) (*models.DailyTask, error)
}

type taskService struct {
	taskRepo repositories.TaskRepository
	userService UserService
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repositories.TaskRepository, userService UserService) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		userService: userService,
	}
}

// GetAllTasks returns all available tasks
func (s *taskService) GetAllTasks() ([]models.Task, error) {
	return s.taskRepo.GetAllTasks()
}

// GetDailyTasks retrieves daily tasks for a user on a specific date
func (s *taskService) GetDailyTasks(userID uint, date time.Time) ([]models.DailyTask, error) {
	// Validate user exists
	_, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.taskRepo.GetDailyTasksByUserAndDate(userID, date)
}

// GenerateDailyTasks generates 3 random tasks for a user for a specific date
func (s *taskService) GenerateDailyTasks(userID uint, date time.Time) ([]models.DailyTask, error) {
	// Validate user exists
	_, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if daily tasks already exist for this date
	existingTasks, err := s.taskRepo.GetDailyTasksByUserAndDate(userID, date)
	if err != nil {
		return nil, err
	}

	// If tasks already exist, delete them first
	if len(existingTasks) > 0 {
		err = s.taskRepo.DeleteDailyTasksByUserAndDate(userID, date)
		if err != nil {
			return nil, err
		}
	}

	// Get 3 random tasks
	randomTasks, err := s.taskRepo.GetRandomTasks(3)
	if err != nil {
		return nil, err
	}

	if len(randomTasks) == 0 {
		return nil, errors.New("no tasks available")
	}

	// Create daily tasks
	var dailyTasks []models.DailyTask
	for _, task := range randomTasks {
		dailyTask := models.DailyTask{
			UserID:    userID,
			TaskID:    task.ID,
			Completed: false,
			Date:      date,
		}

		err = s.taskRepo.CreateDailyTask(&dailyTask)
		if err != nil {
			return nil, err
		}

		// Load the task relationship
		dailyTask.Task = task
		dailyTasks = append(dailyTasks, dailyTask)
	}

	return dailyTasks, nil
}

// CompleteTask marks a daily task as completed and awards points to the user
func (s *taskService) CompleteTask(taskID uint) (*models.DailyTask, error) {
	// Get the daily task
	dailyTask, err := s.taskRepo.GetDailyTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("daily task not found")
		}
		return nil, err
	}

	// Check if task is already completed
	if dailyTask.Completed {
		return nil, errors.New("task already completed")
	}

	// Check if the task is for today (prevent completing old tasks)
	today := time.Now().Truncate(24 * time.Hour)
	taskDate := dailyTask.Date.Truncate(24 * time.Hour)
	if !taskDate.Equal(today) {
		return nil, errors.New("can only complete tasks for today")
	}

	// Mark task as completed
	dailyTask.Completed = true
	err = s.taskRepo.UpdateDailyTask(dailyTask)
	if err != nil {
		return nil, err
	}

	// Award points to user
	err = s.userService.AddPointsToUser(dailyTask.UserID, dailyTask.Task.Points)
	if err != nil {
		// If adding points fails, we should rollback the task completion
		dailyTask.Completed = false
		s.taskRepo.UpdateDailyTask(dailyTask)
		return nil, errors.New("failed to award points to user")
	}

	return dailyTask, nil
} 