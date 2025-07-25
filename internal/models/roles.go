package models

import "time"

// Placeholder models para roles - se implementarán más adelante
// Por ahora solo devolvemos datos mock para mantener la API funcionando

type Role struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type Permission struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	IsActive    bool   `json:"is_active"`
}

// Request/Response models
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50,alphanum"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
}

type UpdateRoleRequest struct {
	DisplayName string `json:"display_name,omitempty" validate:"omitempty,min=2,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type AssignRoleRequest struct {
	UserID    uint       `json:"user_id" validate:"required,min=1"`
	RoleID    uint       `json:"role_id" validate:"required,min=1"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}