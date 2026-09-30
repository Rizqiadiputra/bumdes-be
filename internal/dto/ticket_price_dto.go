package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type TicketPriceResponse struct {
	ID        uuid.UUID `json:"id"`
	TypeName  string    `json:"type_name"`
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTicketPriceResponse(t entity.TicketPrice) TicketPriceResponse {
	return TicketPriceResponse{
		ID:        t.ID,
		TypeName:  t.TypeName,
		Price:     t.Price,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func NewTicketPriceResponseList(items []entity.TicketPrice) []TicketPriceResponse {
	result := make([]TicketPriceResponse, 0, len(items))
	for _, t := range items {
		result = append(result, NewTicketPriceResponse(t))
	}
	return result
}

type CreateTicketPriceRequest struct {
	TypeName string `json:"type_name" binding:"required"`
	Price    int64  `json:"price" binding:"required,min=0"`
	Status   string `json:"status" binding:"required"`
}

type UpdateTicketPriceRequest struct {
	TypeName string `json:"type_name" binding:"required"`
	Price    int64  `json:"price" binding:"required,min=0"`
	Status   string `json:"status" binding:"required"`
}
