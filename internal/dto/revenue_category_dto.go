package dto

import (
	"time"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type RevenueCategoryResponse struct {
	ID           uint      `json:"id"`
	CategoryName string    `json:"category_name"`
	Source       string    `json:"source"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewRevenueCategoryResponse(r entity.RevenueCategory) RevenueCategoryResponse {
	return RevenueCategoryResponse{
		ID:           r.ID,
		CategoryName: r.CategoryName,
		Source:       r.Source,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func NewRevenueCategoryResponseList(items []entity.RevenueCategory) []RevenueCategoryResponse {
	result := make([]RevenueCategoryResponse, 0, len(items))
	for _, r := range items {
		result = append(result, NewRevenueCategoryResponse(r))
	}
	return result
}

type CreateRevenueCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
	Source       string `json:"source" binding:"required"`
	Status       string `json:"status" binding:"required"`
}

type UpdateRevenueCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
	Source       string `json:"source" binding:"required"`
	Status       string `json:"status" binding:"required"`
}
