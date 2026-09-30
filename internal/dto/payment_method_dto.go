package dto

import (
	"time"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type PaymentMethodResponse struct {
	ID        uint      `json:"id"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPaymentMethodResponse(p entity.PaymentMethod) PaymentMethodResponse {
	return PaymentMethodResponse{
		ID:        p.ID,
		Method:    p.Method,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func NewPaymentMethodResponseList(items []entity.PaymentMethod) []PaymentMethodResponse {
	result := make([]PaymentMethodResponse, 0, len(items))
	for _, p := range items {
		result = append(result, NewPaymentMethodResponse(p))
	}
	return result
}

type CreatePaymentMethodRequest struct {
	Method string `json:"method" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type UpdatePaymentMethodRequest struct {
	Method string `json:"method" binding:"required"`
	Status string `json:"status" binding:"required"`
}
