package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingHistory struct {
	ID                 uuid.UUID    `gorm:"primaryKey;type:uuid" json:"id"`
	TenantID           uuid.UUID    `gorm:"not null;type:uuid" json:"tenant_id"`
	Tenant             Tenant       `gorm:"foreignKey:TenantID" json:"-"`
	KiosID             uuid.UUID    `gorm:"not null;type:uuid" json:"kios_id"`
	KiosLocation       KiosLocation `gorm:"foreignKey:KiosID" json:"-"`
	RentalType         string       `gorm:"size:50;not null" json:"rental_type"`
	DateTempo          time.Time    `gorm:"type:date;not null" json:"date_tempo"`
	TotalAmountPayable int64        `gorm:"not null" json:"total_amount_payable"`
	TotalAmountPaid    *int64       `json:"total_amount_paid"`
	Status             string       `gorm:"size:30;not null" json:"status"`
	UserConfirmID      *uuid.UUID   `gorm:"type:uuid" json:"user_confirm"`
	UserConfirm        *User        `gorm:"foreignKey:UserConfirmID" json:"-"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

func (BillingHistory) TableName() string {
	return "billing_histories"
}

func (b *BillingHistory) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
