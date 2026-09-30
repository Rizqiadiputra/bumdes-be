package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type BillingHistoryFilter struct {
	TenantID *uuid.UUID
	Status   string
	Page     int
	Limit    int
}

type BillingHistoryRepository interface {
	FindAll(filter BillingHistoryFilter) ([]entity.BillingHistory, int64, error)
	FindByID(id uuid.UUID) (*entity.BillingHistory, error)
	FindLatestByTenant(tenantID uuid.UUID) (*entity.BillingHistory, error)
	ExistsByTenantAndDateTempo(tenantID uuid.UUID, dateTempo time.Time) (bool, error)
	Create(b *entity.BillingHistory) error
	Update(b *entity.BillingHistory) error
}

type billingHistoryRepository struct {
	db *gorm.DB
}

func NewBillingHistoryRepository(db *gorm.DB) BillingHistoryRepository {
	return &billingHistoryRepository{db: db}
}

func (r *billingHistoryRepository) FindAll(filter BillingHistoryFilter) ([]entity.BillingHistory, int64, error) {
	query := r.db.Model(&entity.BillingHistory{}).
		Preload("Tenant").
		Preload("KiosLocation").
		Preload("UserConfirm")

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	var items []entity.BillingHistory
	err := query.Order("date_tempo desc").Offset(offset).Limit(filter.Limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *billingHistoryRepository) FindByID(id uuid.UUID) (*entity.BillingHistory, error) {
	var item entity.BillingHistory
	err := r.db.
		Preload("Tenant").
		Preload("KiosLocation").
		Preload("UserConfirm").
		First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *billingHistoryRepository) FindLatestByTenant(tenantID uuid.UUID) (*entity.BillingHistory, error) {
	var item entity.BillingHistory
	err := r.db.Where("tenant_id = ?", tenantID).Order("date_tempo desc").First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *billingHistoryRepository) ExistsByTenantAndDateTempo(tenantID uuid.UUID, dateTempo time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&entity.BillingHistory{}).
		Where("tenant_id = ? AND date_tempo = ?", tenantID, dateTempo.Format("2006-01-02")).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *billingHistoryRepository) Create(b *entity.BillingHistory) error {
	return r.db.Create(b).Error
}

func (r *billingHistoryRepository) Update(b *entity.BillingHistory) error {
	return r.db.Save(b).Error
}
