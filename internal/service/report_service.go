package service

import (
	"sort"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

const (
	RevenueUnitParkir = "Parkir"
	RevenueUnitTiket  = "Tiket"
	RevenueUnitWahana = "Wahana"
)

type ReportService interface {
	GetRevenueReport(filter repository.RevenueReportFilter) (*dto.RevenueReportResponse, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GetRevenueReport(filter repository.RevenueReportFilter) (*dto.RevenueReportResponse, error) {
	items := make([]dto.RevenueReportItem, 0)
	var byUnit dto.RevenueByUnit
	var total int64

	includeParkir := filter.Unit == "" || filter.Unit == RevenueUnitParkir
	includeTiket := filter.Unit == "" || filter.Unit == RevenueUnitTiket
	includeWahana := filter.Unit == "" || filter.Unit == RevenueUnitWahana

	if includeParkir {
		parkings, err := s.reportRepo.FindParkingsInRange(filter)
		if err != nil {
			return nil, err
		}
		for _, p := range parkings {
			subtotal := p.Price * int64(p.Count)
			byUnit.Parkir += subtotal
			total += subtotal
			items = append(items, dto.RevenueReportItem{
				ID: p.ID, Unit: RevenueUnitParkir, Jenis: p.Type, Date: p.CreatedAt.Format("2006-01-02"),
				Count: p.Count, Price: p.Price, Subtotal: subtotal,
				PaymentMethodName: p.PaymentMethod.Method, UserName: p.User.Name,
			})
		}
	}

	if includeTiket {
		ticketings, err := s.reportRepo.FindTicketingsInRange(filter)
		if err != nil {
			return nil, err
		}
		for _, t := range ticketings {
			subtotal := t.Price * int64(t.Count)
			byUnit.Tiket += subtotal
			total += subtotal
			items = append(items, dto.RevenueReportItem{
				ID: t.ID, Unit: RevenueUnitTiket, Jenis: t.TypeName, Date: t.Date.Format("2006-01-02"),
				Count: t.Count, Price: t.Price, Subtotal: subtotal,
				PaymentMethodName: t.PaymentMethodName, UserName: t.User.Name,
			})
		}
	}

	if includeWahana {
		attractions, err := s.reportRepo.FindAttractionsInRange(filter)
		if err != nil {
			return nil, err
		}
		for _, a := range attractions {
			subtotal := a.Price * int64(a.Count)
			byUnit.Wahana += subtotal
			total += subtotal
			items = append(items, dto.RevenueReportItem{
				ID: a.ID, Unit: RevenueUnitWahana, Jenis: a.AttractionName, Date: a.Date.Format("2006-01-02"),
				Count: a.Count, Price: a.Price, Subtotal: subtotal,
				PaymentMethodName: a.PaymentMethodName, UserName: a.User.Name,
			})
		}
	}

	sort.SliceStable(items, func(i, j int) bool { return items[i].Date > items[j].Date })

	return &dto.RevenueReportResponse{
		StartDate: filter.StartDate.Format("2006-01-02"),
		EndDate:   filter.EndDate.Format("2006-01-02"),
		Total:     total,
		ByUnit:    byUnit,
		Items:     items,
	}, nil
}
