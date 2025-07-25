package models

import (
	"time"
)

// User model - usando la tabla existente de it-app_user
type User struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	FirebaseID      string     `json:"firebase_id" gorm:"uniqueIndex;size:128;not null"`
	Email           string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	Username        string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
	FirstName       string     `json:"first_name" gorm:"size:100"`
	LastName        string     `json:"last_name" gorm:"size:100"`
	Provider        string     `json:"provider" gorm:"size:50"`
	ProviderID      string     `json:"provider_id" gorm:"size:128"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	LoginCount      int        `json:"login_count" gorm:"default:0"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP     *string    `json:"last_login_ip" gorm:"size:45"`
	LastLoginDevice *string    `json:"last_login_device" gorm:"size:255"`
	Disabled        bool       `json:"disabled" gorm:"default:false"`
	Status          string     `json:"status" gorm:"size:20;default:'active';check:status IN ('active','inactive','pending')"`
	PkAutUseID      string     `json:"pk_aut_use_id" gorm:"size:128"`
}

type CreateUserRequest struct {
	FirebaseID    string `json:"firebase_id" validate:"required,min=1,max=128"`
	Email         string `json:"email" validate:"required,email,max=255"`
	EmailVerified bool   `json:"email_verified"`
	Username      string `json:"username" validate:"required,min=3,max=50,alphanum"`
	FirstName     string `json:"first_name" validate:"max=100"`
	LastName      string `json:"last_name" validate:"max=100"`
	Provider      string `json:"provider" validate:"max=50"`
	ProviderID    string `json:"provider_id" validate:"max=128"`
	Status        string `json:"status" validate:"oneof=active inactive pending"`
}

type UpdateUserRequest struct {
	EmailVerified   *bool   `json:"email_verified"`
	Username        string  `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	FirstName       string  `json:"first_name" validate:"max=100"`
	LastName        string  `json:"last_name" validate:"max=100"`
	Provider        string  `json:"provider" validate:"max=50"`
	ProviderID      string  `json:"provider_id" validate:"max=128"`
	LastLoginIP     string  `json:"last_login_ip" validate:"omitempty,ip"`
	LastLoginDevice string  `json:"last_login_device" validate:"max=255"`
	Disabled        *bool   `json:"disabled"`
	Status          string  `json:"status" validate:"omitempty,oneof=active inactive pending"`
}