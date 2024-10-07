package repositories

import (
	"gorm.io/gorm"
	"marketplace/services/profiles/internal/models"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func (pr *ProfileRepository) CreateProfile(profile *models.Profile) error {
	return pr.DB.Create(profile).Error
}

func (pr *ProfileRepository) GetProfileById(id string) (*models.Profile, error) {
	var profile models.Profile
	result := pr.DB.Where("id = ?", id).First(&profile)

	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}
