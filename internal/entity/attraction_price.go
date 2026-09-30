package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttractionPrice struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string         `gorm:"size:150;not null" json:"name"`
	Price     int64          `gorm:"not null" json:"price"`
	Status    string         `gorm:"size:20;not null" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AttractionPrice) TableName() string {
	return "attraction_prices"
}

func (a *AttractionPrice) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
