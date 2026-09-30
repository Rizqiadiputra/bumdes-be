package dto

import "github.com/google/uuid"

type RevenueReportItem struct {
	ID                uuid.UUID `json:"id"`
	Unit              string    `json:"unit"`
	Jenis             string    `json:"jenis"`
	Date              string    `json:"date"`
	Count             int       `json:"count"`
	Price             int64     `json:"price"`
	Subtotal          int64     `json:"subtotal"`
	PaymentMethodName string    `json:"payment_method_name"`
	UserName          string    `json:"user_name"`
}

type RevenueByUnit struct {
	Parkir int64 `json:"parkir"`
	Tiket  int64 `json:"tiket"`
	Wahana int64 `json:"wahana"`
}

type RevenueReportResponse struct {
	StartDate string              `json:"start_date"`
	EndDate   string              `json:"end_date"`
	Total     int64               `json:"total"`
	ByUnit    RevenueByUnit       `json:"by_unit"`
	Items     []RevenueReportItem `json:"items"`
}
