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

type RoleHandler struct {
	roleRepo repositories.RoleRepositoryInterface
}

func NewRoleHandler(roleRepo repositories.RoleRepositoryInterface) *RoleHandler {
	return &RoleHandler{
		roleRepo: roleRepo,
	}
}

// GetAllRoles maneja GET /roles
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	roles, err := h.roleRepo.GetAllRoles()
	if err != nil {
		log.WithError(err).Error("Failed to fetch roles")
		http.Error(w, "Error fetching roles", http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(roles)).Info("Roles retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    roles,
		"count":   len(roles),
		"message": "Roles retrieved successfully",
	})
}

// GetRoleByID maneja GET /roles/{id}
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid role ID provided")
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := h.roleRepo.GetRoleByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("role_id", id).Error("Failed to fetch role")
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	log.WithField("role_id", id).Info("Role retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    role,
		"message": "Role retrieved successfully",
	})
}

// CreateRole maneja POST /roles
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.CreateRoleRequest

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
		log.WithError(err).Warn("Validation failed for create role request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Crear rol
	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		Active:      true,
	}

	if err := h.roleRepo.CreateRole(role); err != nil {
		log.WithError(err).Error("Failed to create role")
		http.Error(w, "Error creating role", http.StatusInternalServerError)
		return
	}

	log.WithField("role_id", role.ID).Info("Role created successfully")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    role,
		"message": "Role created successfully",
	})
}

// UpdateRole maneja PUT /roles/{id}
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid role ID provided")
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateRoleRequest

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
		log.WithError(err).Warn("Validation failed for update role request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener rol existente
	role, err := h.roleRepo.GetRoleByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("role_id", id).Error("Role not found for update")
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	// Actualizar campos
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Active != nil {
		role.Active = *req.Active
	}

	// Guardar cambios
	if err := h.roleRepo.UpdateRole(role); err != nil {
		log.WithError(err).WithField("role_id", id).Error("Failed to update role")
		http.Error(w, "Error updating role", http.StatusInternalServerError)
		return
	}

	log.WithField("role_id", id).Info("Role updated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    role,
		"message": "Role updated successfully",
	})
}

// DeleteRole maneja DELETE /roles/{id}
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid role ID provided")
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	// Verificar que el rol existe
	_, err = h.roleRepo.GetRoleByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("role_id", id).Error("Role not found for deletion")
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	// Eliminar rol
	if err := h.roleRepo.DeleteRole(uint(id)); err != nil {
		log.WithError(err).WithField("role_id", id).Error("Failed to delete role")
		http.Error(w, "Error deleting role", http.StatusInternalServerError)
		return
	}

	log.WithField("role_id", id).Info("Role deleted successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role deleted successfully",
	})
}

// AssignRoleToUser maneja POST /users/{user_id}/roles
func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	
	// Validar que el UserID no esté vacío
	if userID == "" {
		log.Warn("Empty user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.AssignRoleRequest

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
		log.WithError(err).Warn("Validation failed for assign role request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Asignar rol al usuario
	if err := h.roleRepo.AssignRoleToUser(userID, req.RoleName); err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"user_id": userID,
			"role":    req.RoleName,
		}).Error("Failed to assign role to user")
		http.Error(w, "Error assigning role", http.StatusInternalServerError)
		return
	}

	log.WithFields(map[string]interface{}{
		"user_id": userID,
		"role":    req.RoleName,
	}).Info("Role assigned to user successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role assigned successfully",
	})
}

// RemoveRoleFromUser maneja DELETE /users/{user_id}/roles/{role_name}
func (h *RoleHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	roleName := vars["role_name"]
	
	// Validar que el UserID no esté vacío
	if userID == "" {
		log.Warn("Empty user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	if roleName == "" {
		log.Warn("Role name is required")
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	// Remover rol del usuario
	if err := h.roleRepo.RemoveRoleFromUser(userID, roleName); err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"user_id": userID,
			"role":    roleName,
		}).Error("Failed to remove role from user")
		http.Error(w, "Error removing role", http.StatusInternalServerError)
		return
	}

	log.WithFields(map[string]interface{}{
		"user_id": userID,
		"role":    roleName,
	}).Info("Role removed from user successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role removed successfully",
	})
}

// GetUserRoles maneja GET /users/{user_id}/roles
func (h *RoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	
	// Validar que el UserID no esté vacío
	if userID == "" {
		log.Warn("Empty user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	roles, err := h.roleRepo.GetUserRoles(userID)
	if err != nil {
		log.WithError(err).WithField("user_id", userID).Error("Failed to fetch user roles")
		http.Error(w, "Error fetching user roles", http.StatusInternalServerError)
		return
	}

	log.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(roles),
	}).Info("User roles retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    roles,
		"count":   len(roles),
		"message": "User roles retrieved successfully",
	})
}

