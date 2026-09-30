package repository

import (
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type PaymentMethodFilter struct {
	Search string
}

type PaymentMethodRepository interface {
	FindAll(filter PaymentMethodFilter) ([]entity.PaymentMethod, error)
	FindByID(id uint) (*entity.PaymentMethod, error)
	Create(p *entity.PaymentMethod) error
	Update(p *entity.PaymentMethod) error
	Delete(id uint) error
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) FindAll(filter PaymentMethodFilter) ([]entity.PaymentMethod, error) {
	query := r.db.Order("id asc")

	if filter.Search != "" {
		query = query.Where("method ILIKE ?", "%"+filter.Search+"%")
	}

	var items []entity.PaymentMethod
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *paymentMethodRepository) FindByID(id uint) (*entity.PaymentMethod, error) {
	var item entity.PaymentMethod
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *paymentMethodRepository) Create(p *entity.PaymentMethod) error {
	return r.db.Create(p).Error
}

func (r *paymentMethodRepository) Update(p *entity.PaymentMethod) error {
	return r.db.Save(p).Error
}

func (r *paymentMethodRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.PaymentMethod{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
