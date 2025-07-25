package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"it-user-service/internal/logger"
	"it-user-service/internal/models"
	"it-user-service/internal/repositories"
	"it-user-service/internal/validator"
)

type ProfileHandler struct {
	profileRepo repositories.ProfileRepositoryInterface
}

func NewProfileHandler(profileRepo repositories.ProfileRepositoryInterface) *ProfileHandler {
	return &ProfileHandler{
		profileRepo: profileRepo,
	}
}

// GetUserProfile maneja GET /users/{id}/profile
func (h *ProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.profileRepo.GetByUserID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to fetch user profile")
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	log.WithField("user_id", id).Info("User profile retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    profile,
		"message": "Profile retrieved successfully",
	})
}

// UpdateUserProfile maneja PUT /users/{id}/profile
func (h *ProfileHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateProfileRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for update profile request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener perfil existente o crear uno nuevo
	profile, err := h.profileRepo.GetByUserID(uint(id))
	if err != nil {
		// Si no existe, crear uno nuevo
		profile = &models.UserProfile{
			UserID: uint(id),
		}
	}

	// Actualizar campos
	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}
	if req.Website != "" {
		profile.Website = req.Website
	}
	if req.Location != "" {
		profile.Location = req.Location
	}
	if req.Birthday != nil {
		profile.Birthday = req.Birthday
	}
	if req.Gender != "" {
		profile.Gender = req.Gender
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.Preferences != nil {
		profile.Preferences = req.Preferences
	}
	if req.Privacy != nil {
		profile.Privacy = req.Privacy
	}

	// Guardar cambios
	if profile.ID == 0 {
		err = h.profileRepo.Create(profile)
	} else {
		err = h.profileRepo.Update(profile)
	}

	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to update profile")
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("Profile updated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    profile,
		"message": "Profile updated successfully",
	})
}

// GetUserSettings maneja GET /users/{id}/settings
func (h *ProfileHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	settings, err := h.profileRepo.GetSettingsByUserID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to fetch user settings")
		http.Error(w, "Settings not found", http.StatusNotFound)
		return
	}

	log.WithField("user_id", id).Info("User settings retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    settings,
		"message": "Settings retrieved successfully",
	})
}

// UpdateUserSettings maneja PUT /users/{id}/settings
func (h *ProfileHandler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateSettingsRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for update settings request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener configuraciones existentes o crear nuevas
	settings, err := h.profileRepo.GetSettingsByUserID(uint(id))
	if err != nil {
		// Si no existe, crear nuevas
		settings = &models.UserSettings{
			UserID:   uint(id),
			Language: "en",
			Timezone: "UTC",
			Theme:    "light",
		}
	}

	// Actualizar campos
	if req.Language != "" {
		settings.Language = req.Language
	}
	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}
	if req.Theme != "" {
		settings.Theme = req.Theme
	}
	if req.Notifications != nil {
		settings.Notifications = req.Notifications
	}
	if req.Privacy != nil {
		settings.Privacy = req.Privacy
	}
	if req.Security != nil {
		settings.Security = req.Security
	}

	// Guardar cambios
	if settings.ID == 0 {
		err = h.profileRepo.CreateSettings(settings)
	} else {
		err = h.profileRepo.UpdateSettings(settings)
	}

	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to update settings")
		http.Error(w, "Error updating settings", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("Settings updated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    settings,
		"message": "Settings updated successfully",
	})
}

// GetUserStats maneja GET /users/{id}/stats
func (h *ProfileHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	stats, err := h.profileRepo.GetStatsByUserID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to fetch user stats")
		http.Error(w, "Stats not found", http.StatusNotFound)
		return
	}

	log.WithField("user_id", id).Info("User stats retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    stats,
		"message": "Stats retrieved successfully",
	})
}