package entity

import (
	"time"

	"gorm.io/gorm"
)

type RevenueCategory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CategoryName string         `gorm:"size:150;not null" json:"category_name"`
	Source       string         `gorm:"size:150;not null" json:"source"`
	Status       string         `gorm:"size:20;not null" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RevenueCategory) TableName() string {
	return "revenue_categories"
}
