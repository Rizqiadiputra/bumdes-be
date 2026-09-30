package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type TicketingFilter struct {
	Date      time.Time
	UserID    uuid.UUID
	TypeName  string
	PaymentID uint
}

type TicketingRepository interface {
	FindAll(filter TicketingFilter) ([]entity.Ticketing, error)
	FindByID(id uuid.UUID) (*entity.Ticketing, error)
	Create(t *entity.Ticketing) error
	Update(t *entity.Ticketing) error
	Delete(id uuid.UUID) error
}

type ticketingRepository struct {
	db *gorm.DB
}

func NewTicketingRepository(db *gorm.DB) TicketingRepository {
	return &ticketingRepository{db: db}
}

func (r *ticketingRepository) FindAll(filter TicketingFilter) ([]entity.Ticketing, error) {
	query := r.db.Preload("User").Order("created_at desc").
		Where("date = ?", filter.Date.Format("2006-01-02"))

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.TypeName != "" {
		query = query.Where("type_name = ?", filter.TypeName)
	}
	if filter.PaymentID != 0 {
		query = query.Where("payment_id = ?", filter.PaymentID)
	}

	var items []entity.Ticketing
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ticketingRepository) FindByID(id uuid.UUID) (*entity.Ticketing, error) {
	var item entity.Ticketing
	if err := r.db.Preload("User").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ticketingRepository) Create(t *entity.Ticketing) error {
	return r.db.Create(t).Error
}

func (r *ticketingRepository) Update(t *entity.Ticketing) error {
	return r.db.Save(t).Error
}

func (r *ticketingRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Ticketing{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
