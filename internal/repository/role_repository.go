package repository

import (
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type RoleRepository interface {
	FindAll() ([]entity.Role, error)
	FindByID(id uint) (*entity.Role, error)
	UpdatePermissions(roleID uint, permissions []entity.Permission) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Preload("Permissions").Order("id asc").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindByID(id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) UpdatePermissions(roleID uint, permissions []entity.Permission) error {
	role := entity.Role{ID: roleID}
	return r.db.Model(&role).Association("Permissions").Replace(permissions)
}
