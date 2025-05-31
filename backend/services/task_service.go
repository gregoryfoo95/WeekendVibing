package services

import (
	"errors"
	"fithero-backend/models"
	"fithero-backend/repositories"
	"gorm.io/gorm"
)

type TaskService struct {
	taskRepo        repositories.TaskRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
	achievementRepo repositories.AchievementRepositoryInterface
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repositories.TaskRepositoryInterface, userRepo repositories.UserRepositoryInterface, achievementRepo repositories.AchievementRepositoryInterface) *TaskService {
	return &TaskService{
		taskRepo:        taskRepo,
		userRepo:        userRepo,
		achievementRepo: achievementRepo,
	}
}

// GetAllTasks returns all available tasks (public endpoint)
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.taskRepo.GetAll()
}

// GenerateDailyTasks generates daily tasks for a specific user
func (s *TaskService) GenerateDailyTasks(userID uint) ([]models.DailyTask, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if user already has daily tasks for today
	existingTasks, err := s.taskRepo.GetDailyTasksByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// If user already has tasks for today, return them
	if len(existingTasks) > 0 {
		return existingTasks, nil
	}

	// Generate new daily tasks based on user level
	tasks, err := s.taskRepo.GetTasksByLevel(user.Level)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("no tasks available for user level")
	}

	// Create daily task entries for the user
	var dailyTasks []models.DailyTask
	maxTasks := 3 // Generate 3 daily tasks
	if len(tasks) < maxTasks {
		maxTasks = len(tasks)
	}

	for i := 0; i < maxTasks; i++ {
		task := tasks[i%len(tasks)] // Cycle through available tasks
		dailyTask := models.DailyTask{
			UserID:      userID,
			TaskID:      task.ID,
			Task:        task,
			IsCompleted: false,
			Points:      task.Points,
		}

		createdTask, err := s.taskRepo.CreateDailyTask(&dailyTask)
		if err != nil {
			return nil, err
		}
		dailyTasks = append(dailyTasks, *createdTask)
	}

	return dailyTasks, nil
}

// CompleteTask marks a daily task as completed for a specific user
func (s *TaskService) CompleteTask(userID uint, dailyTaskID uint) (*models.DailyTask, error) {
	// Get the daily task
	dailyTask, err := s.taskRepo.GetDailyTaskByID(dailyTaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("daily task not found")
		}
		return nil, err
	}

	// Authorization check: Ensure the task belongs to the requesting user
	if dailyTask.UserID != userID {
		return nil, errors.New("access denied: you can only complete your own tasks")
	}

	// Check if task is already completed
	if dailyTask.IsCompleted {
		return nil, errors.New("task already completed")
	}

	// Mark as completed
	updateReq := &models.UpdateDailyTaskRequest{
		IsCompleted: &[]bool{true}[0], // Pointer to true
	}

	err = s.taskRepo.UpdateDailyTask(dailyTaskID, updateReq)
	if err != nil {
		return nil, err
	}

	// Award points to user
	err = s.awardPointsToUser(userID, dailyTask.Points)
	if err != nil {
		// Log the error but don't fail the task completion
		// In production, you might want to implement a retry mechanism
		return dailyTask, err
	}

	// Get updated task
	updatedTask, err := s.taskRepo.GetDailyTaskByID(dailyTaskID)
	if err != nil {
		return dailyTask, nil // Return original if we can't get updated
	}

	return updatedTask, nil
}

// GetUserDailyTasks returns daily tasks for a specific user
func (s *TaskService) GetUserDailyTasks(userID uint) ([]models.DailyTask, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.taskRepo.GetDailyTasksByUserID(userID)
}

// GetTaskByID returns a specific task (public endpoint)
func (s *TaskService) GetTaskByID(taskID uint) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return task, nil
}

// Private helper methods

// awardPointsToUser awards points to a user and updates their level
func (s *TaskService) awardPointsToUser(userID uint, points int) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	newPoints := user.Points + points
	newLevel := s.calculateLevelFromPoints(newPoints)
	character := s.getCharacterForLevel(newLevel)

	updateReq := &models.UpdateUserRequest{
		Points:    &newPoints,
		Level:     &newLevel,
		Character: &character,
	}

	return s.userRepo.Update(userID, updateReq)
}

// calculateLevelFromPoints calculates user level based on total points
func (s *TaskService) calculateLevelFromPoints(points int) int {
	switch {
	case points < 100:
		return 1
	case points < 300:
		return 2
	case points < 600:
		return 3
	case points < 1000:
		return 4
	default:
		return 5
	}
}

// getCharacterForLevel returns the character name for a given level
func (s *TaskService) getCharacterForLevel(level int) string {
	characters := map[int]string{
		1: "Rookie Hero",
		2: "Bronze Warrior",
		3: "Silver Champion",
		4: "Gold Legend",
		5: "Platinum Master",
	}

	if character, exists := characters[level]; exists {
		return character
	}
	return "Rookie Hero" // Default fallback
} 