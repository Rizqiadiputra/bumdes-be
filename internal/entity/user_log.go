package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	LogActionView   = "view"
	LogActionCreate = "create"
	LogActionUpdate = "update"
	LogActionDelete = "delete"
)

type UserLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"not null;index;type:uuid" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action      string         `gorm:"size:20;not null;index" json:"action"`
	Module      string         `gorm:"size:100;not null;index" json:"module"`
	Description string         `gorm:"size:255" json:"description"`
	Method      string         `gorm:"size:10;not null" json:"method"`
	Path        string         `gorm:"size:255;not null" json:"path"`
	IPAddress   string         `gorm:"size:64" json:"ip_address"`
	UserAgent   string         `gorm:"size:255" json:"user_agent"`
	Data        datatypes.JSON `gorm:"type:jsonb" json:"data"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (UserLog) TableName() string {
	return "user_logs"
}
