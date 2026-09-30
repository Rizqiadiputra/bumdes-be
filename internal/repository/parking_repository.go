package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type ParkingFilter struct {
	Date     time.Time
	Location string
}

type ParkingRepository interface {
	FindAll(filter ParkingFilter) ([]entity.Parking, error)
	FindByID(id uuid.UUID) (*entity.Parking, error)
	Create(p *entity.Parking) error
	Delete(id uuid.UUID) error
}

type parkingRepository struct {
	db *gorm.DB
}

func NewParkingRepository(db *gorm.DB) ParkingRepository {
	return &parkingRepository{db: db}
}

func (r *parkingRepository) FindAll(filter ParkingFilter) ([]entity.Parking, error) {
	dayStart := time.Date(filter.Date.Year(), filter.Date.Month(), filter.Date.Day(), 0, 0, 0, 0, filter.Date.Location())
	query := r.db.Preload("PaymentMethod").Preload("User").Order("created_at desc").
		Where("created_at >= ? AND created_at < ?", dayStart, dayStart.AddDate(0, 0, 1))

	if filter.Location != "" {
		query = query.Where("location = ?", filter.Location)
	}

	var items []entity.Parking
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *parkingRepository) FindByID(id uuid.UUID) (*entity.Parking, error) {
	var item entity.Parking
	err := r.db.Preload("PaymentMethod").Preload("User").First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *parkingRepository) Create(p *entity.Parking) error {
	return r.db.Create(p).Error
}

func (r *parkingRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Parking{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
