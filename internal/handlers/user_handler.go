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

type UserHandler struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserHandler(userRepo repositories.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// HealthCheck maneja GET /health
func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "user-service"}`))
	
	log.Info("Health check requested")
}

// GetAllUsers maneja GET /users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Parámetros de paginación
	limit := 50 // default
	offset := 0 // default
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	
	users, err := h.userRepo.GetAll(limit, offset)
	if err != nil {
		log.WithError(err).Error("Failed to fetch users")
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(users)).Info("Users retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    users,
		"count":   len(users),
		"limit":   limit,
		"offset":  offset,
		"message": "Users retrieved successfully",
	})
}

// GetUserByID maneja GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to fetch user")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.WithField("user_id", id).Info("User retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

// CreateUser maneja POST /users/create
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.CreateUserRequest

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

	// Establecer valores por defecto
	if req.Status == "" {
		req.Status = "active"
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for create user request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Crear usuario usando el repositorio
	user := &models.User{
		FirebaseID:    req.FirebaseID,
		Email:         req.Email,
		EmailVerified: req.EmailVerified,
		Username:      req.Username,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Provider:      req.Provider,
		ProviderID:    req.ProviderID,
		Status:        req.Status,
	}

	if err := h.userRepo.Create(user); err != nil {
		log.WithError(err).Error("Failed to create user")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", user.ID).Info("User created successfully")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User created successfully",
	})
}

// UpdateUser maneja PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest

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
		log.WithError(err).Warn("Validation failed for update user request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener usuario existente
	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("User not found for update")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Actualizar campos
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.EmailVerified != nil {
		user.EmailVerified = *req.EmailVerified
	}
	if req.Disabled != nil {
		user.Disabled = *req.Disabled
	}

	// Guardar cambios
	if err := h.userRepo.Update(user); err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to update user")
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("User updated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User updated successfully",
	})
}

// DeleteUser maneja DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Verificar que el usuario existe antes de eliminarlo
	_, err = h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("User not found for deletion")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Eliminar usuario
	if err := h.userRepo.Delete(uint(id)); err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("User deleted successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User deleted successfully",
	})
}

// GetUserByFirebaseID maneja GET /users/firebase/{firebase_id}
func (h *UserHandler) GetUserByFirebaseID(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	firebaseID := vars["firebase_id"]
	if firebaseID == "" {
		log.Warn("Firebase ID is required but not provided")
		http.Error(w, "Firebase ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByFirebaseID(firebaseID)
	if err != nil {
		log.WithError(err).WithField("firebase_id", firebaseID).Error("Failed to fetch user by Firebase ID")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.WithField("firebase_id", firebaseID).Info("User retrieved successfully by Firebase ID")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

// SearchUsers maneja GET /users/search
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	query := r.URL.Query().Get("q")
	if query == "" {
		log.Warn("Search query is required")
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Parámetros de paginación
	limit := 20 // default para búsquedas
	offset := 0 // default
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userRepo.SearchUsers(query, limit, offset)
	if err != nil {
		log.WithError(err).WithField("query", query).Error("Failed to search users")
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	log.WithFields(map[string]interface{}{
		"query": query,
		"count": len(users),
	}).Info("Users search completed")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    users,
		"count":   len(users),
		"query":   query,
		"limit":   limit,
		"offset":  offset,
		"message": "Search completed successfully",
	})
}