package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type AttractionPriceResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAttractionPriceResponse(a entity.AttractionPrice) AttractionPriceResponse {
	return AttractionPriceResponse{
		ID:        a.ID,
		Name:      a.Name,
		Price:     a.Price,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func NewAttractionPriceResponseList(items []entity.AttractionPrice) []AttractionPriceResponse {
	result := make([]AttractionPriceResponse, 0, len(items))
	for _, a := range items {
		result = append(result, NewAttractionPriceResponse(a))
	}
	return result
}

type CreateAttractionPriceRequest struct {
	Name   string `json:"name" binding:"required"`
	Price  int64  `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required"`
}

type UpdateAttractionPriceRequest struct {
	Name   string `json:"name" binding:"required"`
	Price  int64  `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required"`
}
