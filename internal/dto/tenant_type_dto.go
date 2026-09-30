package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type TenantTypeResponse struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTenantTypeResponse(t entity.TenantType) TenantTypeResponse {
	return TenantTypeResponse{
		ID:          t.ID,
		Type:        t.Type,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func NewTenantTypeResponseList(items []entity.TenantType) []TenantTypeResponse {
	result := make([]TenantTypeResponse, 0, len(items))
	for _, t := range items {
		result = append(result, NewTenantTypeResponse(t))
	}
	return result
}

type CreateTenantTypeRequest struct {
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required"`
}

type UpdateTenantTypeRequest struct {
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required"`
}
