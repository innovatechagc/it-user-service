package repositories

import "it-user-service/internal/models"

// UserRepositoryInterface define los métodos para el repositorio de usuarios
type UserRepositoryInterface interface {
	// CRUD básico
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByFirebaseID(firebaseID string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll(limit, offset int) ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
	
	// Métodos específicos
	UpdateLoginInfo(id string, loginIP, loginDevice string) error
	GetActiveUsers() ([]models.User, error)
	SearchUsers(query string, limit, offset int) ([]models.User, error)
	CountUsers() (int64, error)
}



// UserProfileRepositoryInterface define los métodos para perfiles de usuario
type UserProfileRepositoryInterface interface {
	GetByUserID(userID string) (*models.UserProfile, error)
	Create(profile *models.UserProfile) error
	Update(profile *models.UserProfile) error
	Delete(userID string) error
	UpdateAvatar(userID string, avatarURL string) error
}

// UserSettingsRepositoryInterface define los métodos para configuraciones de usuario
type UserSettingsRepositoryInterface interface {
	GetByUserID(userID string) (*models.UserSettings, error)
	Create(settings *models.UserSettings) error
	Update(settings *models.UserSettings) error
	Delete(userID string) error
	UpdateLanguage(userID string, language string) error
	UpdateTheme(userID string, theme string) error
}

// UserStatsRepositoryInterface define los métodos para estadísticas de usuario
type UserStatsRepositoryInterface interface {
	GetByUserID(userID string) (*models.UserStats, error)
	Create(stats *models.UserStats) error
	Update(stats *models.UserStats) error
	Delete(userID string) error
	IncrementLoginCount(userID string) error
	IncrementProfileViews(userID string) error
	UpdateLastActive(userID string) error
}

// ProfileRepositoryInterface define los métodos para el repositorio de perfiles
type ProfileRepositoryInterface interface {
	GetByUserID(userID string) (*models.UserProfile, error)
	Create(profile *models.UserProfile) error
	Update(profile *models.UserProfile) error
	Delete(userID string) error
	CreateSettings(settings *models.UserSettings) error
	GetSettingsByUserID(userID string) (*models.UserSettings, error)
	UpdateSettings(settings *models.UserSettings) error
	DeleteSettings(userID string) error
	CreateStats(stats *models.UserStats) error
	GetStatsByUserID(userID string) (*models.UserStats, error)
	UpdateStats(stats *models.UserStats) error
	DeleteStats(userID string) error
	IncrementLoginCount(userID string) error
	UpdateLastLogin(userID string) error
	IncrementProfileViews(userID string) error
	UpdateLastActivity(userID string) error
}

// RoleRepositoryInterface define los métodos para el repositorio de roles
type RoleRepositoryInterface interface {
	GetAllRoles() ([]*models.Role, error)
	GetRoleByID(id uint) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error
	GetActiveRoles() ([]*models.Role, error)
	AssignRoleToUser(userID string, roleName string) error
	RemoveRoleFromUser(userID string, roleName string) error
	GetUserRoles(userID string) ([]*models.UserRole, error)
	UserHasRole(userID string, roleName string) (bool, error)
	UserHasAnyRole(userID string, roleNames []string) (bool, error)
}