package repositories

import (
	"errors"
	"gorm.io/gorm"
	"marketplace/services/profiles/internal/models"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func (pr *ProfileRepository) GetProfileById(id uint) (*models.Profile, error) {
	var profile models.Profile
	result := pr.DB.First(&profile, "id = ?", id)
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
