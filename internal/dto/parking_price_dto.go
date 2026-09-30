package dto

import (
	"time"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type ParkingPriceResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewParkingPriceResponse(p entity.ParkingPrice) ParkingPriceResponse {
	return ParkingPriceResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func NewParkingPriceResponseList(items []entity.ParkingPrice) []ParkingPriceResponse {
	result := make([]ParkingPriceResponse, 0, len(items))
	for _, p := range items {
		result = append(result, NewParkingPriceResponse(p))
	}
	return result
}

type CreateParkingPriceRequest struct {
	Name   string `json:"name" binding:"required"`
	Price  int64  `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required"`
}

type UpdateParkingPriceRequest struct {
	Name   string `json:"name" binding:"required"`
	Price  int64  `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required"`
}
