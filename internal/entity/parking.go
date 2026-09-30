package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Parking struct {
	ID              uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Type            string         `gorm:"size:150;not null" json:"type"`
	Location        string         `gorm:"size:255;not null" json:"location"`
	Count           int            `gorm:"not null" json:"count"`
	PaymentMethodID uint           `gorm:"not null" json:"payment_method_id"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	Price           int64          `gorm:"not null" json:"price"`
	UserID          uuid.UUID      `gorm:"not null;type:uuid" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status          string         `gorm:"size:20;not null" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Parking) TableName() string {
	return "parkings"
}

func (p *Parking) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
