package dto

import (
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

const (
	KiosAvailabilityTerisi   = "Terisi"
	KiosAvailabilityTersedia = "Tersedia"
)

type KiosAvailableResponse struct {
	ID       uuid.UUID  `json:"id"`
	Kios     string     `json:"kios"`
	Size     string     `json:"size"`
	Price    int64      `json:"price"`
	TenantID *uuid.UUID `json:"tenant_id"`
	Status   string     `json:"status"`
}

// NewKiosAvailableResponse builds the availability response. price is the
// caller-resolved price: the tenant's price when the kios is occupied
// (status Terisi), or the kios_locations price when it's Tersedia.
func NewKiosAvailableResponse(k entity.KiosLocation, price int64) KiosAvailableResponse {
	status := KiosAvailabilityTersedia
	if k.TenantID != nil {
		status = KiosAvailabilityTerisi
	}

	return KiosAvailableResponse{
		ID:       k.ID,
		Kios:     k.Kios,
		Size:     k.Size,
		Price:    price,
		TenantID: k.TenantID,
		Status:   status,
	}
}
