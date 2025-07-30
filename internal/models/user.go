package models

import (
	"time"
)

type User struct {
	ID              string     `json:"id" gorm:"primaryKey;type:uuid"`
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

// User Profile models - Modelos relacionados con el perfil del usuario
type UserProfile struct {
	ID          uint                   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      string                 `json:"user_id" gorm:"not null;uniqueIndex;type:uuid"`
	Avatar      string                 `json:"avatar,omitempty" gorm:"size:500"`
	Bio         string                 `json:"bio,omitempty" gorm:"type:text"`
	Website     string                 `json:"website,omitempty" gorm:"size:255"`
	Location    string                 `json:"location,omitempty" gorm:"size:100"`
	Birthday    *time.Time             `json:"birthday,omitempty"`
	Gender      string                 `json:"gender,omitempty" gorm:"size:20"`
	Phone       string                 `json:"phone,omitempty" gorm:"size:20"`
	Preferences string                 `json:"preferences,omitempty" gorm:"type:jsonb"` // JSON en PostgreSQL
	Privacy     string                 `json:"privacy,omitempty" gorm:"type:jsonb"`     // JSON en PostgreSQL
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// User Settings models - Modelos relacionados con configuraciones del usuario
type UserSettings struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        string    `json:"user_id" gorm:"not null;uniqueIndex;type:uuid"`
	Language      string    `json:"language" gorm:"size:10;default:'en'"`
	Timezone      string    `json:"timezone" gorm:"size:50;default:'UTC'"`
	Theme         string    `json:"theme" gorm:"size:20;default:'light'"`
	Notifications string    `json:"notifications" gorm:"type:jsonb"` // JSON en PostgreSQL
	Privacy       string    `json:"privacy" gorm:"type:jsonb"`       // JSON en PostgreSQL
	Security      string    `json:"security" gorm:"type:jsonb"`      // JSON en PostgreSQL
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// User Statistics models - Modelos relacionados con estadísticas del usuario
type UserStats struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       string     `json:"user_id" gorm:"not null;uniqueIndex;type:uuid"`
	LoginCount   int        `json:"login_count" gorm:"default:0"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	ProfileViews int        `json:"profile_views" gorm:"default:0"`
	AccountAge   int        `json:"account_age_days" gorm:"default:0"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
// Role models - Modelos relacionados con roles
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description string    `json:"description" gorm:"size:255"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserRole models - Modelos relacionados con roles de usuario
type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"not null;type:uuid"`
	Role      string    `json:"role" gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Response models - Modelos para respuestas
type ProfileResponse struct {
	User     User          `json:"user"`
	Profile  *UserProfile  `json:"profile,omitempty"`
	Settings *UserSettings `json:"settings,omitempty"`
	Stats    *UserStats    `json:"stats,omitempty"`
}

type UserWithRoles struct {
	User  User       `json:"user"`
	Roles []UserRole `json:"roles"`
}

// Request models - Modelos para requests
type CreateProfileRequest struct {
	Avatar      string     `json:"avatar,omitempty" validate:"omitempty,url,max=500"`
	Bio         string     `json:"bio,omitempty" validate:"max=1000"`
	Website     string     `json:"website,omitempty" validate:"omitempty,url,max=255"`
	Location    string     `json:"location,omitempty" validate:"max=100"`
	Birthday    *time.Time `json:"birthday,omitempty"`
	Gender      string     `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	Phone       string     `json:"phone,omitempty" validate:"omitempty,max=20"`
	Preferences string     `json:"preferences,omitempty"`
	Privacy     string     `json:"privacy,omitempty"`
}

type CreateSettingsRequest struct {
	Language      string `json:"language,omitempty" validate:"omitempty,min=2,max=10"`
	Timezone      string `json:"timezone,omitempty" validate:"max=50"`
	Theme         string `json:"theme,omitempty" validate:"omitempty,oneof=light dark auto"`
	Notifications string `json:"notifications,omitempty"`
	Privacy       string `json:"privacy,omitempty"`
	Security      string `json:"security,omitempty"`
}

type UpdateProfileRequest struct {
	Avatar      string     `json:"avatar,omitempty" validate:"omitempty,url,max=500"`
	Bio         string     `json:"bio,omitempty" validate:"max=1000"`
	Website     string     `json:"website,omitempty" validate:"omitempty,url,max=255"`
	Location    string     `json:"location,omitempty" validate:"max=100"`
	Birthday    *time.Time `json:"birthday,omitempty"`
	Gender      string     `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	Phone       string     `json:"phone,omitempty" validate:"omitempty,max=20"`
	Preferences string     `json:"preferences,omitempty"`
	Privacy     string     `json:"privacy,omitempty"`
}

type UpdateSettingsRequest struct {
	Language      string `json:"language,omitempty" validate:"omitempty,min=2,max=10"`
	Timezone      string `json:"timezone,omitempty" validate:"max=50"`
	Theme         string `json:"theme,omitempty" validate:"omitempty,oneof=light dark auto"`
	Notifications string `json:"notifications,omitempty"`
	Privacy       string `json:"privacy,omitempty"`
	Security      string `json:"security,omitempty"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
	Active      bool   `json:"active"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
	Active      *bool  `json:"active,omitempty"`
}

type AssignRoleRequest struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	RoleName string `json:"role_name" validate:"required,min=2,max=50"`
}