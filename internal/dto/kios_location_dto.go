package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type KiosLocationResponse struct {
	ID        uuid.UUID  `json:"id"`
	Kios      string     `json:"kios"`
	Size      string     `json:"size"`
	Price     int64      `json:"price"`
	Status    string     `json:"status"`
	TenantID  *uuid.UUID `json:"tenant_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewKiosLocationResponse(k entity.KiosLocation) KiosLocationResponse {
	return KiosLocationResponse{
		ID:        k.ID,
		Kios:      k.Kios,
		Size:      k.Size,
		Price:     k.Price,
		Status:    k.Status,
		TenantID:  k.TenantID,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func NewKiosLocationResponseList(items []entity.KiosLocation) []KiosLocationResponse {
	result := make([]KiosLocationResponse, 0, len(items))
	for _, k := range items {
		result = append(result, NewKiosLocationResponse(k))
	}
	return result
}

type CreateKiosLocationRequest struct {
	Kios   string `json:"kios" binding:"required"`
	Size   string `json:"size" binding:"required"`
	Price  int64  `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required"`
}

type UpdateKiosLocationRequest struct {
	Kios   string `json:"kios" binding:"required"`
	Size   string `json:"size" binding:"required"`
	Price  int64  `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required"`
}
