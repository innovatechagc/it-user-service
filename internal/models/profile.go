package models

import "time"

// Placeholder models - estas funcionalidades se implementarán más adelante
// Por ahora solo devolvemos datos mock para mantener la API funcionando

type UserProfile struct {
	UserID   uint   `json:"user_id"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Website  string `json:"website,omitempty"`
	Location string `json:"location,omitempty"`
}

type UserSettings struct {
	UserID        uint   `json:"user_id"`
	Language      string `json:"language"`
	Timezone      string `json:"timezone"`
	Theme         string `json:"theme"`
	Notifications string `json:"notifications,omitempty"`
}

type UserStats struct {
	UserID       uint       `json:"user_id"`
	LoginCount   int        `json:"login_count"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	ProfileViews int        `json:"profile_views"`
	IsActive     bool       `json:"is_active"`
}

// Request/Response models for profiles
type UpdateProfileRequest struct {
	Avatar      string     `json:"avatar,omitempty" validate:"omitempty,url"`
	Bio         string     `json:"bio,omitempty" validate:"max=500"`
	Website     string     `json:"website,omitempty" validate:"omitempty,url"`
	Location    string     `json:"location,omitempty" validate:"max=100"`
	Birthday    *time.Time `json:"birthday,omitempty"`
	Gender      string     `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	Phone       string     `json:"phone,omitempty" validate:"omitempty,e164"`
	Preferences string     `json:"preferences,omitempty"`
	Privacy     string     `json:"privacy,omitempty"`
}

type UpdateSettingsRequest struct {
	Language      string `json:"language,omitempty" validate:"omitempty,oneof=en es fr de it pt"`
	Timezone      string `json:"timezone,omitempty" validate:"max=50"`
	Theme         string `json:"theme,omitempty" validate:"omitempty,oneof=light dark auto"`
	Notifications string `json:"notifications,omitempty"`
	Privacy       string `json:"privacy,omitempty"`
	Security      string `json:"security,omitempty"`
}