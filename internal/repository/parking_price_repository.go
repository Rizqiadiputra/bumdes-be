package repository

import (
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type ParkingPriceFilter struct {
	Search string
}

type ParkingPriceRepository interface {
	FindAll(filter ParkingPriceFilter) ([]entity.ParkingPrice, error)
	FindByID(id uint) (*entity.ParkingPrice, error)
	Create(p *entity.ParkingPrice) error
	Update(p *entity.ParkingPrice) error
	Delete(id uint) error
}

type parkingPriceRepository struct {
	db *gorm.DB
}

func NewParkingPriceRepository(db *gorm.DB) ParkingPriceRepository {
	return &parkingPriceRepository{db: db}
}

func (r *parkingPriceRepository) FindAll(filter ParkingPriceFilter) ([]entity.ParkingPrice, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	var items []entity.ParkingPrice
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *parkingPriceRepository) FindByID(id uint) (*entity.ParkingPrice, error) {
	var item entity.ParkingPrice
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *parkingPriceRepository) Create(p *entity.ParkingPrice) error {
	return r.db.Create(p).Error
}

func (r *parkingPriceRepository) Update(p *entity.ParkingPrice) error {
	return r.db.Save(p).Error
}

func (r *parkingPriceRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.ParkingPrice{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
