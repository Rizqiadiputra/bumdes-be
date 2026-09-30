package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KiosLocation struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Kios      string         `gorm:"size:150;not null" json:"kios"`
	Size      string         `gorm:"size:50;not null" json:"size"`
	Price     int64          `gorm:"not null" json:"price"`
	Status    string         `gorm:"size:20;not null" json:"status"`
	TenantID  *uuid.UUID     `gorm:"type:uuid" json:"tenant_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (KiosLocation) TableName() string {
	return "kios_locations"
}

func (k *KiosLocation) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}
