package repositories

import (
	"fithero-backend/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetAll() ([]models.User, error)
	GetTopUsersByPoints(limit int) ([]models.User, error)
	Update(id uint, updates *models.UpdateUserRequest) error
	Delete(id uint) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetTopUsersByPoints(limit int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("is_active = ?", true).
		Order("points DESC").
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(id uint, updates *models.UpdateUserRequest) error {
	user := &models.User{}
	if err := r.db.First(user, id).Error; err != nil {
		return err
	}

	updateData := make(map[string]interface{})
	
	if updates.Username != nil {
		updateData["username"] = *updates.Username
	}
	if updates.Email != nil {
		updateData["email"] = *updates.Email
	}
	if updates.FirstName != nil {
		updateData["first_name"] = *updates.FirstName
	}
	if updates.LastName != nil {
		updateData["last_name"] = *updates.LastName
	}
	if updates.Level != nil {
		updateData["level"] = *updates.Level
	}
	if updates.Points != nil {
		updateData["points"] = *updates.Points
	}
	if updates.Character != nil {
		updateData["character"] = *updates.Character
	}
	if updates.JobTitle != nil {
		updateData["job_title"] = *updates.JobTitle
	}

	if len(updateData) > 0 {
		return r.db.Model(user).Updates(updateData).Error
	}
	
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
} 