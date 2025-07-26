package repositories

import (
	"time"
	"gorm.io/gorm"
	"it-user-service/internal/models"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepositoryInterface {
	return &ProfileRepository{db: db}
}

// Profile CRUD operations

// GetByUserID obtiene el perfil de un usuario
func (r *ProfileRepository) GetByUserID(userID uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Create crea un nuevo perfil de usuario
func (r *ProfileRepository) Create(profile *models.UserProfile) error {
	return r.db.Create(profile).Error
}

// Update actualiza un perfil de usuario
func (r *ProfileRepository) Update(profile *models.UserProfile) error {
	return r.db.Save(profile).Error
}

// Delete elimina un perfil de usuario
func (r *ProfileRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserProfile{}).Error
}

// Settings CRUD operations

// CreateSettings crea nuevas configuraciones de usuario
func (r *ProfileRepository) CreateSettings(settings *models.UserSettings) error {
	return r.db.Create(settings).Error
}

// GetSettingsByUserID obtiene las configuraciones de un usuario
func (r *ProfileRepository) GetSettingsByUserID(userID uint) (*models.UserSettings, error) {
	var settings models.UserSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// UpdateSettings actualiza las configuraciones de usuario
func (r *ProfileRepository) UpdateSettings(settings *models.UserSettings) error {
	return r.db.Save(settings).Error
}

// DeleteSettings elimina las configuraciones de usuario
func (r *ProfileRepository) DeleteSettings(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSettings{}).Error
}

// Stats CRUD operations

// CreateStats crea nuevas estadísticas de usuario
func (r *ProfileRepository) CreateStats(stats *models.UserStats) error {
	return r.db.Create(stats).Error
}

// GetStatsByUserID obtiene las estadísticas de un usuario
func (r *ProfileRepository) GetStatsByUserID(userID uint) (*models.UserStats, error) {
	var stats models.UserStats
	err := r.db.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpdateStats actualiza las estadísticas de usuario
func (r *ProfileRepository) UpdateStats(stats *models.UserStats) error {
	return r.db.Save(stats).Error
}

// DeleteStats elimina las estadísticas de usuario
func (r *ProfileRepository) DeleteStats(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserStats{}).Error
}

// Stats operations

// IncrementLoginCount incrementa el contador de logins en la tabla users
func (r *ProfileRepository) IncrementLoginCount(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("login_count", gorm.Expr("login_count + 1")).Error
}

// UpdateLastLogin actualiza la fecha del último login en la tabla users
func (r *ProfileRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
		}).Error
}

// IncrementProfileViews incrementa el contador de vistas del perfil
func (r *ProfileRepository) IncrementProfileViews(userID uint) error {
	// Primero intentar actualizar si existe
	result := r.db.Model(&models.UserStats{}).Where("user_id = ?", userID).
		UpdateColumn("profile_views", gorm.Expr("profile_views + 1"))
	
	if result.Error != nil {
		return result.Error
	}
	
	// Si no se actualizó ninguna fila, crear el registro
	if result.RowsAffected == 0 {
		stats := &models.UserStats{
			UserID:       userID,
			ProfileViews: 1,
			IsActive:     true,
		}
		return r.db.Create(stats).Error
	}
	
	return nil
}

// UpdateLastActivity actualiza la última actividad del usuario
func (r *ProfileRepository) UpdateLastActivity(userID uint) error {
	now := time.Now()
	
	// Primero intentar actualizar si existe
	result := r.db.Model(&models.UserStats{}).Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"last_active_at": now,
			"is_active":      true,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	// Si no se actualizó ninguna fila, crear el registro
	if result.RowsAffected == 0 {
		stats := &models.UserStats{
			UserID:       userID,
			LastActiveAt: &now,
			IsActive:     true,
		}
		return r.db.Create(stats).Error
	}
	
	return nil
}

// GetCompleteProfile obtiene el perfil completo con usuario, perfil, configuraciones y estadísticas
func (r *ProfileRepository) GetCompleteProfile(userID uint) (*models.ProfileResponse, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	response := &models.ProfileResponse{
		User: user,
	}
	
	// Obtener perfil (opcional)
	var profile models.UserProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err == nil {
		response.Profile = &profile
	}
	
	// Obtener configuraciones (opcional)
	var settings models.UserSettings
	if err := r.db.Where("user_id = ?", userID).First(&settings).Error; err == nil {
		response.Settings = &settings
	}
	
	// Obtener estadísticas (opcional)
	var stats models.UserStats
	if err := r.db.Where("user_id = ?", userID).First(&stats).Error; err == nil {
		response.Stats = &stats
	}
	
	return response, nil
}

// CreateCompleteProfile crea un perfil completo con configuraciones iniciales
func (r *ProfileRepository) CreateCompleteProfile(userID uint, profileReq *models.CreateProfileRequest, settingsReq *models.CreateSettingsRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Crear perfil si se proporciona
		if profileReq != nil {
			profile := &models.UserProfile{
				UserID:      userID,
				Avatar:      profileReq.Avatar,
				Bio:         profileReq.Bio,
				Website:     profileReq.Website,
				Location:    profileReq.Location,
				Birthday:    profileReq.Birthday,
				Gender:      profileReq.Gender,
				Phone:       profileReq.Phone,
				Preferences: profileReq.Preferences,
				Privacy:     profileReq.Privacy,
			}
			if err := tx.Create(profile).Error; err != nil {
				return err
			}
		}
		
		// Crear configuraciones si se proporciona
		if settingsReq != nil {
			settings := &models.UserSettings{
				UserID:        userID,
				Language:      settingsReq.Language,
				Timezone:      settingsReq.Timezone,
				Theme:         settingsReq.Theme,
				Notifications: settingsReq.Notifications,
				Privacy:       settingsReq.Privacy,
				Security:      settingsReq.Security,
			}
			if err := tx.Create(settings).Error; err != nil {
				return err
			}
		}
		
		// Crear estadísticas iniciales
		stats := &models.UserStats{
			UserID:   userID,
			IsActive: true,
		}
		if err := tx.Create(stats).Error; err != nil {
			return err
		}
		
		return nil
	})
}