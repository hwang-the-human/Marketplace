package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"marketplace/services/profiles/internal/dto"
	"marketplace/services/profiles/internal/models"
	"marketplace/services/profiles/internal/services"
	"net/http"
)

type ProfileHandler struct {
	ProfileService *services.ProfileService
	Validator      *validator.Validate
}

func (ph *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var createProfileDTO dto.CreateProfileDTO

	// Декодируем тело запроса в DTO
	if err := json.NewDecoder(r.Body).Decode(&createProfileDTO); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if err := ph.Validator.Struct(createProfileDTO); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	profile := models.Profile{
		FirstName: createProfileDTO.FirstName,
		LastName:  createProfileDTO.LastName,
	}

	if err := ph.ProfileService.CreateProfile(&profile); err != nil {
		http.Error(w, "Error creating profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}
