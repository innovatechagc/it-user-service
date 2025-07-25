package repositories

import "it-user-service/internal/models"

// UserRepositoryInterface define los métodos para el repositorio de usuarios
type UserRepositoryInterface interface {
	// CRUD básico
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByFirebaseID(firebaseID string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll(limit, offset int) ([]*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	
	// Métodos específicos
	UpdateLoginInfo(id uint, loginIP, loginDevice string) error
	GetActiveUsers() ([]*models.User, error)
	SearchUsers(query string, limit, offset int) ([]*models.User, error)
	CountUsers() (int64, error)
}

// ProfileRepositoryInterface define los métodos para perfiles de usuario (placeholder)
type ProfileRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserProfile, error)
	Create(profile *models.UserProfile) error
	Update(profile *models.UserProfile) error
	Delete(userID uint) error
	
	// Settings methods (placeholder)
	CreateSettings(settings *models.UserSettings) error
	GetSettingsByUserID(userID uint) (*models.UserSettings, error)
	UpdateSettings(settings *models.UserSettings) error
	DeleteSettings(userID uint) error
	
	// Stats methods (placeholder)
	CreateStats(stats *models.UserStats) error
	GetStatsByUserID(userID uint) (*models.UserStats, error)
	UpdateStats(stats *models.UserStats) error
	DeleteStats(userID uint) error
	
	// Stats operations (placeholder)
	IncrementLoginCount(userID uint) error
	UpdateLastLogin(userID uint) error
	IncrementProfileViews(userID uint) error
}

// RoleRepositoryInterface define los métodos para roles y permisos (placeholder)
type RoleRepositoryInterface interface {
	GetAllRoles() ([]*models.Role, error)
	GetRoleByID(id uint) (*models.Role, error)
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error
	
	// User-Role relationships (placeholder)
	AssignRoleToUser(userID, roleID uint, expiresAt interface{}) error
	RemoveRoleFromUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]*models.Role, error)
	GetUserPermissions(userID uint) ([]*models.Permission, error)
}