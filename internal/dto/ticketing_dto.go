package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

const ticketingDateFormat = "2006-01-02"

type TicketingResponse struct {
	ID                uuid.UUID `json:"id"`
	TicketPriceID     uuid.UUID `json:"ticket_price_id"`
	TypeName          string    `json:"type_name"`
	Price             int64     `json:"price"`
	PaymentID         uint      `json:"payment_id"`
	PaymentMethodName string    `json:"payment_method_name"`
	Count             int       `json:"count"`
	UserID            uuid.UUID `json:"user_id"`
	UserName          string    `json:"user_name,omitempty"`
	Date              string    `json:"date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewTicketingResponse(t entity.Ticketing) TicketingResponse {
	return TicketingResponse{
		ID:                t.ID,
		TicketPriceID:     t.TicketPriceID,
		TypeName:          t.TypeName,
		Price:             t.Price,
		PaymentID:         t.PaymentID,
		PaymentMethodName: t.PaymentMethodName,
		Count:             t.Count,
		UserID:            t.UserID,
		UserName:          t.User.Name,
		Date:              t.Date.Format(ticketingDateFormat),
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
	}
}

func NewTicketingResponseList(items []entity.Ticketing) []TicketingResponse {
	result := make([]TicketingResponse, 0, len(items))
	for _, t := range items {
		result = append(result, NewTicketingResponse(t))
	}
	return result
}

type CreateTicketingRequest struct {
	TicketPriceID uuid.UUID `json:"ticket_price_id" binding:"required"`
	PaymentID     uint      `json:"payment_id" binding:"required"`
	Count         int       `json:"count" binding:"required,min=1"`
}

type UpdateTicketingRequest struct {
	TicketPriceID uuid.UUID `json:"ticket_price_id" binding:"required"`
	PaymentID     uint      `json:"payment_id" binding:"required"`
	Count         int       `json:"count" binding:"required,min=1"`
}
