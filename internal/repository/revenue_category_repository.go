package repository

import (
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type RevenueCategoryFilter struct {
	Search string
}

type RevenueCategoryRepository interface {
	FindAll(filter RevenueCategoryFilter) ([]entity.RevenueCategory, error)
	FindByID(id uint) (*entity.RevenueCategory, error)
	Create(r *entity.RevenueCategory) error
	Update(r *entity.RevenueCategory) error
	Delete(id uint) error
}

type revenueCategoryRepository struct {
	db *gorm.DB
}

func NewRevenueCategoryRepository(db *gorm.DB) RevenueCategoryRepository {
	return &revenueCategoryRepository{db: db}
}

func (r *revenueCategoryRepository) FindAll(filter RevenueCategoryFilter) ([]entity.RevenueCategory, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("category_name ILIKE ? OR source ILIKE ?", like, like)
	}

	var items []entity.RevenueCategory
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *revenueCategoryRepository) FindByID(id uint) (*entity.RevenueCategory, error) {
	var item entity.RevenueCategory
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *revenueCategoryRepository) Create(item *entity.RevenueCategory) error {
	return r.db.Create(item).Error
}

func (r *revenueCategoryRepository) Update(item *entity.RevenueCategory) error {
	return r.db.Save(item).Error
}

func (r *revenueCategoryRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.RevenueCategory{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
