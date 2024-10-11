package services

import (
	"marketplace/services/profiles/internal/models"
	"marketplace/services/profiles/internal/repositories"
)

type ProfileService interface {
	GetProfileByID(id string) (*models.Profile, error)
	CreateProfile(profile *models.Profile) (*models.Profile, error)
}

type profileService struct {
	ProfileRepository repositories.ProfileRepository
}

func NewProfileService(profileRepository repositories.ProfileRepository) ProfileService {
	return &profileService{ProfileRepository: profileRepository}
}

func (ps *profileService) GetProfileByID(id string) (*models.Profile, error) {
	return ps.ProfileRepository.GetProfileById(id)
}

func (ps *profileService) CreateProfile(profile *models.Profile) (*models.Profile, error) {
	return ps.ProfileRepository.CreateProfile(profile)
}
