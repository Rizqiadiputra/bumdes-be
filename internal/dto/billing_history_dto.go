package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

const billingDateFormat = "2006-01-02"

type BillingHistoryResponse struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	TenantName         string     `json:"tenant_name"`
	KiosID             uuid.UUID  `json:"kios_id"`
	KiosName           string     `json:"kios_name"`
	RentalType         string     `json:"rental_type"`
	DateTempo          string     `json:"date_tempo"`
	TotalAmountPayable int64      `json:"total_amount_payable"`
	TotalAmountPaid    *int64     `json:"total_amount_paid"`
	Status             string     `json:"status"`
	UserConfirmID      *uuid.UUID `json:"user_confirm"`
	UserConfirmName    *string    `json:"user_confirm_name"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func NewBillingHistoryResponse(b entity.BillingHistory) BillingHistoryResponse {
	var totalAmountPaid *int64
	if b.TotalAmountPaid != nil {
		totalAmountPaid = b.TotalAmountPaid
	}

	var userConfirmName *string
	if b.UserConfirm != nil {
		name := b.UserConfirm.Name
		userConfirmName = &name
	}

	return BillingHistoryResponse{
		ID:                 b.ID,
		TenantID:           b.TenantID,
		TenantName:         b.Tenant.Name,
		KiosID:             b.KiosID,
		KiosName:           b.KiosLocation.Kios,
		RentalType:         b.RentalType,
		DateTempo:          b.DateTempo.Format(billingDateFormat),
		TotalAmountPayable: b.TotalAmountPayable,
		TotalAmountPaid:    totalAmountPaid,
		Status:             b.Status,
		UserConfirmID:      b.UserConfirmID,
		UserConfirmName:    userConfirmName,
		CreatedAt:          b.CreatedAt,
		UpdatedAt:          b.UpdatedAt,
	}
}

func NewBillingHistoryResponseList(items []entity.BillingHistory) []BillingHistoryResponse {
	result := make([]BillingHistoryResponse, 0, len(items))
	for _, item := range items {
		result = append(result, NewBillingHistoryResponse(item))
	}
	return result
}

type PaginatedBillingHistoryResponse struct {
	BillingHistories []BillingHistoryResponse `json:"billing_histories"`
	Meta             PaginationMeta           `json:"meta"`
}
