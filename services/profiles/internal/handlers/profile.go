package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"marketplace/services/profiles/internal/models"
	"marketplace/services/profiles/internal/services"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	ProfileService *services.ProfileService
}

func (ph *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	profile, err := ph.ProfileService.GetProfileByID(uint(profileID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (ph *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ph.ProfileService.CreateProfile(&profile); err != nil {
		http.Error(w, "Error creating profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (ph *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	profile.ID = uint(profileID)

	if err := ph.ProfileService.UpdateProfile(&profile); err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}
