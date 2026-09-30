package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticketing struct {
	ID                uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	TicketPriceID     uuid.UUID      `gorm:"not null;type:uuid" json:"ticket_price_id"`
	TicketPrice       TicketPrice    `gorm:"foreignKey:TicketPriceID" json:"-"`
	TypeName          string         `gorm:"size:150;not null" json:"type_name"`
	Price             int64          `gorm:"not null" json:"price"`
	PaymentID         uint           `gorm:"not null" json:"payment_id"`
	PaymentMethod     PaymentMethod  `gorm:"foreignKey:PaymentID" json:"-"`
	PaymentMethodName string         `gorm:"size:150;not null" json:"payment_method_name"`
	Count             int            `gorm:"not null" json:"count"`
	UserID            uuid.UUID      `gorm:"not null;type:uuid" json:"user_id"`
	User              User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Date              time.Time      `gorm:"type:date;not null;index" json:"date"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Ticketing) TableName() string {
	return "ticketings"
}

func (t *Ticketing) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
