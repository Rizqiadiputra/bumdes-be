package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type AttractionFilter struct {
	Date              time.Time
	AttractionName    string
	PaymentMethodName string
}

type AttractionRepository interface {
	FindAll(filter AttractionFilter) ([]entity.Attraction, error)
	FindByID(id uuid.UUID) (*entity.Attraction, error)
	Create(a *entity.Attraction) error
	Update(a *entity.Attraction) error
	Delete(id uuid.UUID) error
}

type attractionRepository struct {
	db *gorm.DB
}

func NewAttractionRepository(db *gorm.DB) AttractionRepository {
	return &attractionRepository{db: db}
}

func (r *attractionRepository) FindAll(filter AttractionFilter) ([]entity.Attraction, error) {
	query := r.db.Preload("User").Order("created_at desc").
		Where("date = ?", filter.Date.Format("2006-01-02"))

	if filter.AttractionName != "" {
		query = query.Where("attraction_name = ?", filter.AttractionName)
	}
	if filter.PaymentMethodName != "" {
		query = query.Where("payment_method_name = ?", filter.PaymentMethodName)
	}

	var items []entity.Attraction
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *attractionRepository) FindByID(id uuid.UUID) (*entity.Attraction, error) {
	var item entity.Attraction
	if err := r.db.Preload("User").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *attractionRepository) Create(a *entity.Attraction) error {
	return r.db.Create(a).Error
}

func (r *attractionRepository) Update(a *entity.Attraction) error {
	return r.db.Save(a).Error
}

func (r *attractionRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Attraction{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
