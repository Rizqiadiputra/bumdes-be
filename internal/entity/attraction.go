package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attraction struct {
	ID                uuid.UUID       `gorm:"primaryKey;type:uuid" json:"id"`
	AttractionPriceID uuid.UUID       `gorm:"not null;type:uuid" json:"attraction_price_id"`
	AttractionPrice   AttractionPrice `gorm:"foreignKey:AttractionPriceID" json:"-"`
	AttractionName    string          `gorm:"size:150;not null" json:"attraction_name"`
	Price             int64           `gorm:"not null" json:"price"`
	Count             int             `gorm:"not null" json:"count"`
	PaymentID         uint            `gorm:"not null" json:"payment_id"`
	PaymentMethod     PaymentMethod   `gorm:"foreignKey:PaymentID" json:"-"`
	PaymentMethodName string          `gorm:"size:150;not null" json:"payment_method_name"`
	UserID            uuid.UUID       `gorm:"not null;type:uuid" json:"user_id"`
	User              User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Date              time.Time       `gorm:"type:date;not null;index" json:"date"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Attraction) TableName() string {
	return "attractions"
}

func (a *Attraction) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
