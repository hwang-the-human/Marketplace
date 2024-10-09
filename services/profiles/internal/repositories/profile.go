package repositories

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace/services/profiles/internal/models"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func (pr *ProfileRepository) GetProfileById(id string) (*models.Profile, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	var profile models.Profile
	result := pr.DB.First(&profile, "id = ?", uid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

func (pr *ProfileRepository) CreateProfile(profile *models.Profile) (*models.Profile, error) {
	result := pr.DB.Create(profile)
	if result.Error != nil {
		return nil, result.Error
	}
	return profile, nil
}
