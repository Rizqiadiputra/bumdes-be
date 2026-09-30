package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type TenantTypeFilter struct {
	Search string
	Limit  int
}

type TenantTypeRepository interface {
	FindAll(filter TenantTypeFilter) ([]entity.TenantType, error)
	FindByID(id uuid.UUID) (*entity.TenantType, error)
	Create(t *entity.TenantType) error
	Update(t *entity.TenantType) error
	Delete(id uuid.UUID) error
}

type tenantTypeRepository struct {
	db *gorm.DB
}

func NewTenantTypeRepository(db *gorm.DB) TenantTypeRepository {
	return &tenantTypeRepository{db: db}
}

func (r *tenantTypeRepository) FindAll(filter TenantTypeFilter) ([]entity.TenantType, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		query = query.Where("type ILIKE ?", "%"+filter.Search+"%")
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	var items []entity.TenantType
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *tenantTypeRepository) FindByID(id uuid.UUID) (*entity.TenantType, error) {
	var item entity.TenantType
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *tenantTypeRepository) Create(t *entity.TenantType) error {
	return r.db.Create(t).Error
}

func (r *tenantTypeRepository) Update(t *entity.TenantType) error {
	return r.db.Save(t).Error
}

func (r *tenantTypeRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.TenantType{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
