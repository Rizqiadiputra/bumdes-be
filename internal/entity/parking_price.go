package entity

import (
	"time"

	"gorm.io/gorm"
)

type ParkingPrice struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:150;not null" json:"name"`
	Price     int64          `gorm:"not null" json:"price"`
	Status    string         `gorm:"size:20;not null" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ParkingPrice) TableName() string {
	return "parking_prices"
}
