package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type ParkingResponse struct {
	ID                uuid.UUID `json:"id"`
	Type              string    `json:"type"`
	Location          string    `json:"location"`
	Count             int       `json:"count"`
	PaymentMethodID   uint      `json:"payment_method_id"`
	PaymentMethodName string    `json:"payment_method_name,omitempty"`
	Price             int64     `json:"price"`
	UserID            uuid.UUID `json:"user_id"`
	UserName          string    `json:"user_name,omitempty"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewParkingResponse(p entity.Parking) ParkingResponse {
	return ParkingResponse{
		ID:                p.ID,
		Type:              p.Type,
		Location:          p.Location,
		Count:             p.Count,
		PaymentMethodID:   p.PaymentMethodID,
		PaymentMethodName: p.PaymentMethod.Method,
		Price:             p.Price,
		UserID:            p.UserID,
		UserName:          p.User.Name,
		Status:            p.Status,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func NewParkingResponseList(items []entity.Parking) []ParkingResponse {
	result := make([]ParkingResponse, 0, len(items))
	for _, p := range items {
		result = append(result, NewParkingResponse(p))
	}
	return result
}

type CreateParkingRequest struct {
	ParkingPriceID  uint   `json:"parking_price_id" binding:"required"`
	Location        string `json:"location" binding:"required"`
	Count           int    `json:"count" binding:"required,min=1"`
	PaymentMethodID uint   `json:"payment_method_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
}
