package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type TenantFilter struct {
	Search string
	Page   int
	Limit  int
}

type TenantRepository interface {
	FindAll(filter TenantFilter) ([]entity.Tenant, int64, error)
	FindByID(id uuid.UUID) (*entity.Tenant, error)
	// FindActiveForBilling returns active tenants whose rental_type is not
	// excludeRentalType and whose end_periode has not passed (or is null).
	FindActiveForBilling(activeStatus, excludeRentalType string) ([]entity.Tenant, error)
	Create(t *entity.Tenant) error
	Update(t *entity.Tenant) error
	Delete(id uuid.UUID) error
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) FindAll(filter TenantFilter) ([]entity.Tenant, int64, error) {
	query := r.db.Model(&entity.Tenant{})

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	var items []entity.Tenant
	err := query.Order("created_at desc").Offset(offset).Limit(filter.Limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *tenantRepository) FindByID(id uuid.UUID) (*entity.Tenant, error) {
	var item entity.Tenant
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *tenantRepository) FindActiveForBilling(activeStatus, excludeRentalType string) ([]entity.Tenant, error) {
	var items []entity.Tenant
	err := r.db.
		Where("status = ?", activeStatus).
		Where("rental_type <> ?", excludeRentalType).
		Where("end_periode IS NULL OR end_periode >= CURRENT_DATE").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *tenantRepository) Create(t *entity.Tenant) error {
	return r.db.Create(t).Error
}

func (r *tenantRepository) Update(t *entity.Tenant) error {
	return r.db.Save(t).Error
}

func (r *tenantRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Tenant{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
