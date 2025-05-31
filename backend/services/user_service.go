package services

import (
	"errors"
	"fithero-backend/models"
	"fithero-backend/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo        repositories.UserRepositoryInterface
	taskRepo        repositories.TaskRepositoryInterface
	achievementRepo repositories.AchievementRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepositoryInterface, taskRepo repositories.TaskRepositoryInterface, achievementRepo repositories.AchievementRepositoryInterface) *UserService {
	return &UserService{
		userRepo:        userRepo,
		taskRepo:        taskRepo,
		achievementRepo: achievementRepo,
	}
}

// CreateUser creates a new user with business logic validation
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	if existingUser, _ := s.userRepo.GetByEmail(req.Email); existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create new user with default values
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Level:     1,
		Points:    0,
		Character: "Rookie Hero",
		JobTitle:  "Fitness Novice",
		IsActive:  true,
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser updates user information with business logic
func (s *UserService) UpdateUser(id uint, req *models.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check for email conflicts
	if req.Email != nil && *req.Email != user.Email {
		if existingUser, _ := s.userRepo.GetByEmail(*req.Email); existingUser != nil {
			return errors.New("email already exists")
		}
	}

	// Update level based on points if points are being updated
	if req.Points != nil {
		newLevel := s.calculateLevelFromPoints(*req.Points)
		req.Level = &newLevel
		character := s.getCharacterForLevel(newLevel)
		req.Character = &character
	}

	// Update level and character if level is being updated directly
	if req.Level != nil {
		character := s.getCharacterForLevel(*req.Level)
		req.Character = &character
	}

	return s.userRepo.Update(id, req)
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(id)
}

// GetLeaderboard returns top users ordered by points
func (s *UserService) GetLeaderboard(limit int) ([]models.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Sort by points (simple implementation - could be moved to repository)
	// For now, just return all users
	return users, nil
}

// UpdateUserLevel updates user level based on their current level
func (s *UserService) UpdateUserLevel(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	newLevel := s.calculateLevelFromPoints(user.Points)
	if newLevel != user.Level {
		character := s.getCharacterForLevel(newLevel)
		updateReq := &models.UpdateUserRequest{
			Level:     &newLevel,
			Character: &character,
		}
		return s.userRepo.Update(userID, updateReq)
	}
	
	return nil
}

// AddPointsToUser adds points to a user and updates their level
func (s *UserService) AddPointsToUser(userID uint, points int) error {
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

// GetUserDailyTasks retrieves daily tasks for a specific user
func (s *UserService) GetUserDailyTasks(userID uint) ([]models.DailyTask, error) {
	return s.taskRepo.GetDailyTasksByUserID(userID)
}

// GetUserAchievements retrieves achievements for a specific user
func (s *UserService) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	return s.achievementRepo.GetUserAchievements(userID)
}

// calculateLevelFromPoints calculates user level based on total points
func (s *UserService) calculateLevelFromPoints(points int) int {
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
func (s *UserService) getCharacterForLevel(level int) string {
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