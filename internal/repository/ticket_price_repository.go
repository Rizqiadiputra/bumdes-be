package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type TicketPriceFilter struct {
	Search string
}

type TicketPriceRepository interface {
	FindAll(filter TicketPriceFilter) ([]entity.TicketPrice, error)
	FindByID(id uuid.UUID) (*entity.TicketPrice, error)
	Create(t *entity.TicketPrice) error
	Update(t *entity.TicketPrice) error
	Delete(id uuid.UUID) error
}

type ticketPriceRepository struct {
	db *gorm.DB
}

func NewTicketPriceRepository(db *gorm.DB) TicketPriceRepository {
	return &ticketPriceRepository{db: db}
}

func (r *ticketPriceRepository) FindAll(filter TicketPriceFilter) ([]entity.TicketPrice, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		query = query.Where("type_name ILIKE ?", "%"+filter.Search+"%")
	}

	var items []entity.TicketPrice
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ticketPriceRepository) FindByID(id uuid.UUID) (*entity.TicketPrice, error) {
	var item entity.TicketPrice
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ticketPriceRepository) Create(t *entity.TicketPrice) error {
	return r.db.Create(t).Error
}

func (r *ticketPriceRepository) Update(t *entity.TicketPrice) error {
	return r.db.Save(t).Error
}

func (r *ticketPriceRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.TicketPrice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
