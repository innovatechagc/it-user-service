package repositories

import "it-user-service/internal/models"

// UserRepositoryInterface define los métodos para el repositorio de usuarios
type UserRepositoryInterface interface {
	// CRUD básico
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByFirebaseID(firebaseID string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll(limit, offset int) ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	
	// Métodos específicos
	UpdateLoginInfo(id uint, loginIP, loginDevice string) error
	GetActiveUsers() ([]models.User, error)
	SearchUsers(query string, limit, offset int) ([]models.User, error)
	CountUsers() (int64, error)
}



// UserProfileRepositoryInterface define los métodos para perfiles de usuario
type UserProfileRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserProfile, error)
	Create(profile *models.UserProfile) error
	Update(profile *models.UserProfile) error
	Delete(userID uint) error
	UpdateAvatar(userID uint, avatarURL string) error
}

// UserSettingsRepositoryInterface define los métodos para configuraciones de usuario
type UserSettingsRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserSettings, error)
	Create(settings *models.UserSettings) error
	Update(settings *models.UserSettings) error
	Delete(userID uint) error
	UpdateLanguage(userID uint, language string) error
	UpdateTheme(userID uint, theme string) error
}

// UserStatsRepositoryInterface define los métodos para estadísticas de usuario
type UserStatsRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserStats, error)
	Create(stats *models.UserStats) error
	Update(stats *models.UserStats) error
	Delete(userID uint) error
	IncrementLoginCount(userID uint) error
	IncrementProfileViews(userID uint) error
	UpdateLastActive(userID uint) error
}