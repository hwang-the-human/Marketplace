package repositories

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace/services/profiles/internal/models"
	"marketplace/shared/db"
)

type ProfileRepository interface {
	GetProfileById(id string) (*models.Profile, error)
	CreateProfile(profile *models.Profile) (*models.Profile, error)
	DeleteProfileById(id string) error
}

type profileRepository struct {
	Database db.Database
}

func NewProfileRepository(database db.Database) ProfileRepository {
	return &profileRepository{Database: database}
}

func (pr *profileRepository) GetProfileById(id string) (*models.Profile, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	var profile models.Profile
	result := pr.Database.GetDB().First(&profile, "id = ?", uid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

func (pr *profileRepository) CreateProfile(profile *models.Profile) (*models.Profile, error) {
	result := pr.Database.GetDB().Create(profile)
	if result.Error != nil {
		return nil, result.Error
	}
	return profile, nil
}

func (pr *profileRepository) DeleteProfileById(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %v", err)
	}

	result := pr.Database.GetDB().Delete(&models.Profile{}, "id = ?", uid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
