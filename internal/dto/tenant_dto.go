package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

const tenantDateFormat = "2006-01-02"

type TenantResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	TenantTypeID   uuid.UUID `json:"tenant_type_id"`
	Type           string    `json:"type"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	KiosLocationID uuid.UUID `json:"kios_location_id"`
	KiosName       string    `json:"kios_name"`
	Price          int64     `json:"price"`
	Status         string    `json:"status"`
	RentalType     string    `json:"rental_type"`
	StartPayment   string    `json:"start_payment"`
	EndPeriode     *string   `json:"end_periode"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewTenantResponse(t entity.Tenant) TenantResponse {
	var endPeriode *string
	if t.EndPeriode != nil {
		formatted := t.EndPeriode.Format(tenantDateFormat)
		endPeriode = &formatted
	}

	return TenantResponse{
		ID:             t.ID,
		Name:           t.Name,
		TenantTypeID:   t.TenantTypeID,
		Type:           t.Type,
		Phone:          t.Phone,
		Email:          t.Email,
		KiosLocationID: t.KiosLocationID,
		KiosName:       t.KiosName,
		Price:          t.Price,
		Status:         t.Status,
		RentalType:     t.RentalType,
		StartPayment:   t.StartPayment.Format(tenantDateFormat),
		EndPeriode:     endPeriode,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}

func NewTenantResponseList(items []entity.Tenant) []TenantResponse {
	result := make([]TenantResponse, 0, len(items))
	for _, t := range items {
		result = append(result, NewTenantResponse(t))
	}
	return result
}

type PaginatedTenantResponse struct {
	Tenants []TenantResponse `json:"tenants"`
	Meta    PaginationMeta   `json:"meta"`
}

type CreateTenantRequest struct {
	Name           string `json:"name" binding:"required"`
	TenantTypeID   string `json:"tenant_type_id" binding:"required,uuid"`
	Phone          string `json:"phone" binding:"required"`
	Email          string `json:"email" binding:"omitempty,email"`
	KiosLocationID string `json:"kios_location_id" binding:"required,uuid"`
	Price          int64  `json:"price" binding:"required,min=0"`
	Status         string `json:"status" binding:"required"`
	RentalType     string `json:"rental_type"`
	StartPayment   string `json:"start_payment" binding:"required,datetime=2006-01-02"`
	EndPeriode     string `json:"end_periode" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateTenantRequest struct {
	Name           string `json:"name" binding:"required"`
	TenantTypeID   string `json:"tenant_type_id" binding:"required,uuid"`
	Phone          string `json:"phone" binding:"required"`
	Email          string `json:"email" binding:"omitempty,email"`
	KiosLocationID string `json:"kios_location_id" binding:"required,uuid"`
	Price          int64  `json:"price" binding:"required,min=0"`
	Status         string `json:"status" binding:"required"`
	RentalType     string `json:"rental_type" binding:"required"`
	StartPayment   string `json:"start_payment" binding:"required,datetime=2006-01-02"`
	EndPeriode     string `json:"end_periode" binding:"omitempty,datetime=2006-01-02"`
}
