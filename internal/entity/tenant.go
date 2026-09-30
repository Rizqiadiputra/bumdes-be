package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Name           string         `gorm:"size:150;not null" json:"name"`
	TenantTypeID   uuid.UUID      `gorm:"not null;type:uuid" json:"tenant_type_id"`
	TenantType     TenantType     `gorm:"foreignKey:TenantTypeID" json:"-"`
	Type           string         `gorm:"size:150;not null" json:"type"`
	Phone          string         `gorm:"size:30;not null" json:"phone"`
	Email          string         `gorm:"size:150" json:"email"`
	KiosLocationID uuid.UUID      `gorm:"not null;type:uuid" json:"kios_location_id"`
	KiosLocation   KiosLocation   `gorm:"foreignKey:KiosLocationID" json:"-"`
	KiosName       string         `gorm:"size:150;not null" json:"kios_name"`
	Price          int64          `gorm:"not null" json:"price"`
	Status         string         `gorm:"size:20;not null" json:"status"`
	RentalType     string         `gorm:"size:50;not null" json:"rental_type"`
	StartPayment   time.Time      `gorm:"type:date;not null" json:"start_payment"`
	EndPeriode     *time.Time     `gorm:"type:date" json:"end_periode"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tenant) TableName() string {
	return "tenants"
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
