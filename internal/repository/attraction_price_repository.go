package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type AttractionPriceFilter struct {
	Search string
}

type AttractionPriceRepository interface {
	FindAll(filter AttractionPriceFilter) ([]entity.AttractionPrice, error)
	FindByID(id uuid.UUID) (*entity.AttractionPrice, error)
	Create(a *entity.AttractionPrice) error
	Update(a *entity.AttractionPrice) error
	Delete(id uuid.UUID) error
}

type attractionPriceRepository struct {
	db *gorm.DB
}

func NewAttractionPriceRepository(db *gorm.DB) AttractionPriceRepository {
	return &attractionPriceRepository{db: db}
}

func (r *attractionPriceRepository) FindAll(filter AttractionPriceFilter) ([]entity.AttractionPrice, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	var items []entity.AttractionPrice
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *attractionPriceRepository) FindByID(id uuid.UUID) (*entity.AttractionPrice, error) {
	var item entity.AttractionPrice
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *attractionPriceRepository) Create(a *entity.AttractionPrice) error {
	return r.db.Create(a).Error
}

func (r *attractionPriceRepository) Update(a *entity.AttractionPrice) error {
	return r.db.Save(a).Error
}

func (r *attractionPriceRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.AttractionPrice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
