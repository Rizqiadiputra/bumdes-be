package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

const attractionDateFormat = "2006-01-02"

type AttractionResponse struct {
	ID                uuid.UUID `json:"id"`
	AttractionPriceID uuid.UUID `json:"attraction_price_id"`
	AttractionName    string    `json:"attraction_name"`
	Price             int64     `json:"price"`
	Count             int       `json:"count"`
	PaymentID         uint      `json:"payment_id"`
	PaymentMethodName string    `json:"payment_method_name"`
	UserID            uuid.UUID `json:"user_id"`
	UserName          string    `json:"user_name,omitempty"`
	Date              string    `json:"date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewAttractionResponse(a entity.Attraction) AttractionResponse {
	return AttractionResponse{
		ID:                a.ID,
		AttractionPriceID: a.AttractionPriceID,
		AttractionName:    a.AttractionName,
		Price:             a.Price,
		Count:             a.Count,
		PaymentID:         a.PaymentID,
		PaymentMethodName: a.PaymentMethodName,
		UserID:            a.UserID,
		UserName:          a.User.Name,
		Date:              a.Date.Format(attractionDateFormat),
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

func NewAttractionResponseList(items []entity.Attraction) []AttractionResponse {
	result := make([]AttractionResponse, 0, len(items))
	for _, a := range items {
		result = append(result, NewAttractionResponse(a))
	}
	return result
}

type CreateAttractionRequest struct {
	AttractionPriceID uuid.UUID `json:"attraction_price_id" binding:"required"`
	PaymentID         uint      `json:"payment_id" binding:"required"`
	Count             int       `json:"count" binding:"required,min=1"`
}

type UpdateAttractionRequest struct {
	AttractionPriceID uuid.UUID `json:"attraction_price_id" binding:"required"`
	PaymentID         uint      `json:"payment_id" binding:"required"`
	Count             int       `json:"count" binding:"required,min=1"`
}
