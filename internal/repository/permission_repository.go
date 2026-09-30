package repository

import (
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type PermissionRepository interface {
	FindAll() ([]entity.Permission, error)
	FindByIDs(ids []uint) ([]entity.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll() ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Order("id asc").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepository) FindByIDs(ids []uint) ([]entity.Permission, error) {
	var permissions []entity.Permission
	if len(ids) == 0 {
		return permissions, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
