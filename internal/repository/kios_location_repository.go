package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type KiosLocationFilter struct {
	Search string
	Limit  int
}

type KiosLocationRepository interface {
	FindAll(filter KiosLocationFilter) ([]entity.KiosLocation, error)
	FindByID(id uuid.UUID) (*entity.KiosLocation, error)
	Create(k *entity.KiosLocation) error
	Update(k *entity.KiosLocation) error
	Delete(id uuid.UUID) error
	SetTenant(id uuid.UUID, tenantID *uuid.UUID) error
}

type kiosLocationRepository struct {
	db *gorm.DB
}

func NewKiosLocationRepository(db *gorm.DB) KiosLocationRepository {
	return &kiosLocationRepository{db: db}
}

func (r *kiosLocationRepository) FindAll(filter KiosLocationFilter) ([]entity.KiosLocation, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		query = query.Where("kios ILIKE ?", "%"+filter.Search+"%")
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	var items []entity.KiosLocation
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *kiosLocationRepository) FindByID(id uuid.UUID) (*entity.KiosLocation, error) {
	var item entity.KiosLocation
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *kiosLocationRepository) Create(k *entity.KiosLocation) error {
	return r.db.Create(k).Error
}

func (r *kiosLocationRepository) Update(k *entity.KiosLocation) error {
	return r.db.Save(k).Error
}

func (r *kiosLocationRepository) SetTenant(id uuid.UUID, tenantID *uuid.UUID) error {
	return r.db.Model(&entity.KiosLocation{}).Where("id = ?", id).Update("tenant_id", tenantID).Error
}

func (r *kiosLocationRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.KiosLocation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
