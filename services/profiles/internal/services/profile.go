package services

import (
	"marketplace/services/profiles/internal/models"
	"marketplace/services/profiles/internal/repositories"
)

type ProfileService struct {
	ProfileRepository *repositories.ProfileRepository
}

func (ps *ProfileService) GetProfileByID(id string) (*models.Profile, error) {
	return ps.ProfileRepository.GetProfileById(id)
}

func (ps *ProfileService) CreateProfile(profile *models.Profile) (*models.Profile, error) {
	return ps.ProfileRepository.CreateProfile(profile)
}
