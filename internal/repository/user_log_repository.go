package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type UserLogFilter struct {
	UserID uuid.UUID
	Action string
	Module string
	Page   int
	Limit  int
}

type UserLogRepository interface {
	Create(log *entity.UserLog) error
	FindAll(filter UserLogFilter) ([]entity.UserLog, int64, error)
}

type userLogRepository struct {
	db *gorm.DB
}

func NewUserLogRepository(db *gorm.DB) UserLogRepository {
	return &userLogRepository{db: db}
}

func (r *userLogRepository) Create(log *entity.UserLog) error {
	return r.db.Create(log).Error
}

func (r *userLogRepository) FindAll(filter UserLogFilter) ([]entity.UserLog, int64, error) {
	query := r.db.Model(&entity.UserLog{})

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []entity.UserLog
	offset := (filter.Page - 1) * filter.Limit
	err := query.Preload("User").
		Order("created_at desc").
		Offset(offset).
		Limit(filter.Limit).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
