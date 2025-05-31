package services

import (
	"errors"
	"fithero-backend/models"
	"fithero-backend/repositories"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
	GetLeaderboard(limit int) ([]models.User, error)
	UpdateUserLevel(userID uint) error
	AddPointsToUser(userID uint, points int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with business logic validation
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Check if username already exists
	if existingUser, _ := s.userRepo.GetByUsername(req.Username); existingUser != nil {
		return nil, errors.New("username already exists")
	}

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
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
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
func (s *userService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check for username conflicts
	if req.Username != nil && *req.Username != user.Username {
		if existingUser, _ := s.userRepo.GetByUsername(*req.Username); existingUser != nil {
			return nil, errors.New("username already exists")
		}
		user.Username = *req.Username
	}

	// Check for email conflicts
	if req.Email != nil && *req.Email != user.Email {
		if existingUser, _ := s.userRepo.GetByEmail(*req.Email); existingUser != nil {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	// Update other fields if provided
	if req.Level != nil {
		user.Level = *req.Level
		user.Character = s.getCharacterForLevel(*req.Level)
	}
	if req.Points != nil {
		user.Points = *req.Points
		// Update level based on points
		s.updateLevelBasedOnPoints(user)
	}
	if req.Character != nil {
		user.Character = *req.Character
	}
	if req.JobTitle != nil {
		user.JobTitle = *req.JobTitle
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser soft deletes a user
func (s *userService) DeleteUser(id uint) error {
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
func (s *userService) GetLeaderboard(limit int) ([]models.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	
	return s.userRepo.GetLeaderboard(limit)
}

// UpdateUserLevel updates user level based on their current level
func (s *userService) UpdateUserLevel(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	s.updateLevelBasedOnPoints(user)
	return s.userRepo.Update(user)
}

// AddPointsToUser adds points to a user and updates their level
func (s *userService) AddPointsToUser(userID uint, points int) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Points += points
	s.updateLevelBasedOnPoints(user)
	
	return s.userRepo.Update(user)
}

// updateLevelBasedOnPoints updates the user's level and character based on points
func (s *userService) updateLevelBasedOnPoints(user *models.User) {
	newLevel := s.calculateLevelFromPoints(user.Points)
	if newLevel != user.Level {
		user.Level = newLevel
		user.Character = s.getCharacterForLevel(newLevel)
	}
}

// calculateLevelFromPoints calculates user level based on total points
func (s *userService) calculateLevelFromPoints(points int) int {
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
func (s *userService) getCharacterForLevel(level int) string {
	characters := map[int]string{
		1: "Rookie Hero",
		2: "Fitness Apprentice",
		3: "Health Guardian",
		4: "Wellness Warrior",
		5: "Ultimate Hero",
	}
	
	if character, exists := characters[level]; exists {
		return character
	}
	return "Rookie Hero"
} 