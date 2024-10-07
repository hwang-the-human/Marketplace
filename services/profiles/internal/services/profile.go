package services

import (
	"errors"
	"gorm.io/gorm"
	"marketplace/services/profiles/internal/models"
	"marketplace/services/profiles/internal/repositories"
)

type ProfileService struct {
	ProfileRepository *repositories.ProfileRepository
}

func (ps *ProfileService) GetProfileByID(id uint) (*models.Profile, error) {
	var profile models.Profile
	if err := ps.DB.First(&profile, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (ps *ProfileService) CreateProfile(profile *models.Profile) error {
	return ps.DB.Create(profile).Error
}

func (ps *ProfileService) UpdateProfile(profile *models.Profile) error {
	return ps.DB.Save(profile).Error
}
