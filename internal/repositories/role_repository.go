package repositories

import (
	"gorm.io/gorm"
	"it-user-service/internal/models"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &RoleRepository{db: db}
}

// Role CRUD operations

// GetAllRoles obtiene todos los roles
func (r *RoleRepository) GetAllRoles() ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// GetRoleByID obtiene un rol por ID
func (r *RoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName obtiene un rol por nombre
func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole crea un nuevo rol
func (r *RoleRepository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

// UpdateRole actualiza un rol
func (r *RoleRepository) UpdateRole(role *models.Role) error {
	return r.db.Save(role).Error
}

// DeleteRole elimina un rol
func (r *RoleRepository) DeleteRole(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// GetActiveRoles obtiene todos los roles activos
func (r *RoleRepository) GetActiveRoles() ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.Where("active = ?", true).Find(&roles).Error
	return roles, err
}

// User-Role relationships (usando string role según tu SQL)

// AssignRoleToUser asigna un rol a un usuario
func (r *RoleRepository) AssignRoleToUser(userID string, roleName string) error {
	userRole := &models.UserRole{
		UserID: userID,
		Role:   roleName,
	}
	return r.db.Create(userRole).Error
}

// RemoveRoleFromUser remueve un rol de un usuario
func (r *RoleRepository) RemoveRoleFromUser(userID string, roleName string) error {
	return r.db.Where("user_id = ? AND role = ?", userID, roleName).
		Delete(&models.UserRole{}).Error
}

// GetUserRoles obtiene todos los roles de un usuario
func (r *RoleRepository) GetUserRoles(userID string) ([]*models.UserRole, error) {
	var userRoles []*models.UserRole
	err := r.db.Where("user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

// GetUserWithRoles obtiene un usuario con sus roles
func (r *RoleRepository) GetUserWithRoles(userID string) (*models.UserWithRoles, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	
	roles, err := r.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}
	
	// Convertir []*models.UserRole a []models.UserRole
	userRoles := make([]models.UserRole, len(roles))
	for i, role := range roles {
		userRoles[i] = *role
	}
	
	return &models.UserWithRoles{
		User:  user,
		Roles: userRoles,
	}, nil
}

// GetUsersWithRole obtiene todos los usuarios que tienen un rol específico
func (r *RoleRepository) GetUsersWithRole(roleName string) ([]*models.User, error) {
	var users []*models.User
	err := r.db.Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role = ?", roleName).
		Find(&users).Error
	return users, err
}

// Role checking

// UserHasRole verifica si un usuario tiene un rol específico
func (r *RoleRepository) UserHasRole(userID string, roleName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role = ?", userID, roleName).
		Count(&count).Error
	return count > 0, err
}

// UserHasAnyRole verifica si un usuario tiene alguno de los roles especificados
func (r *RoleRepository) UserHasAnyRole(userID string, roleNames []string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role IN ?", userID, roleNames).
		Count(&count).Error
	return count > 0, err
}

// Bulk operations

// AssignMultipleRolesToUser asigna múltiples roles a un usuario
func (r *RoleRepository) AssignMultipleRolesToUser(userID string, roleNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, roleName := range roleNames {
			userRole := &models.UserRole{
				UserID: userID,
				Role:   roleName,
			}
			if err := tx.Create(userRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveMultipleRolesFromUser remueve múltiples roles de un usuario
func (r *RoleRepository) RemoveMultipleRolesFromUser(userID string, roleNames []string) error {
	return r.db.Where("user_id = ? AND role IN ?", userID, roleNames).
		Delete(&models.UserRole{}).Error
}

// RemoveAllUserRoles remueve todos los roles de un usuario
func (r *RoleRepository) RemoveAllUserRoles(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error
}