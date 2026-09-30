package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantType struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Type        string         `gorm:"size:150;not null" json:"type"`
	Description string         `gorm:"size:255" json:"description"`
	Status      string         `gorm:"size:20;not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TenantType) TableName() string {
	return "tenant_types"
}

func (t *TenantType) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
