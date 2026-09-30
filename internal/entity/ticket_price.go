package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketPrice struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	TypeName  string         `gorm:"size:150;not null" json:"type_name"`
	Price     int64          `gorm:"not null" json:"price"`
	Status    string         `gorm:"size:20;not null" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TicketPrice) TableName() string {
	return "ticket_prices"
}

func (t *TicketPrice) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
