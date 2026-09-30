package entity

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Method    string         `gorm:"size:150;not null" json:"method"`
	Status    string         `gorm:"size:20;not null" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}
