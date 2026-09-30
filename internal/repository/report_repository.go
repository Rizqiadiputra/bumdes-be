package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

const reportDateLayout = "2006-01-02"

// RevenueReportFilter membatasi laporan pendapatan pada rentang tanggal, dan
// opsional pada satu unit usaha (Parkir, Tiket, Wahana) beserta jenis
// pendapatannya (mis. tipe kendaraan, jenis tiket, atau nama wahana).
type RevenueReportFilter struct {
	StartDate time.Time
	EndDate   time.Time
	Unit      string
	Jenis     string
}

type ReportRepository interface {
	FindParkingsInRange(filter RevenueReportFilter) ([]entity.Parking, error)
	FindTicketingsInRange(filter RevenueReportFilter) ([]entity.Ticketing, error)
	FindAttractionsInRange(filter RevenueReportFilter) ([]entity.Attraction, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) FindParkingsInRange(filter RevenueReportFilter) ([]entity.Parking, error) {
	dayStart := truncateToDate(filter.StartDate)
	dayEndExclusive := truncateToDate(filter.EndDate).AddDate(0, 0, 1)

	query := r.db.Preload("PaymentMethod").Preload("User").Order("created_at desc").
		Where("created_at >= ? AND created_at < ?", dayStart, dayEndExclusive)
	if filter.Jenis != "" {
		query = query.Where("type = ?", filter.Jenis)
	}

	var items []entity.Parking
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *reportRepository) FindTicketingsInRange(filter RevenueReportFilter) ([]entity.Ticketing, error) {
	query := r.db.Preload("User").Order("date desc").
		Where("date >= ? AND date <= ?", filter.StartDate.Format(reportDateLayout), filter.EndDate.Format(reportDateLayout))
	if filter.Jenis != "" {
		query = query.Where("type_name = ?", filter.Jenis)
	}

	var items []entity.Ticketing
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *reportRepository) FindAttractionsInRange(filter RevenueReportFilter) ([]entity.Attraction, error) {
	query := r.db.Preload("User").Order("date desc").
		Where("date >= ? AND date <= ?", filter.StartDate.Format(reportDateLayout), filter.EndDate.Format(reportDateLayout))
	if filter.Jenis != "" {
		query = query.Where("attraction_name = ?", filter.Jenis)
	}

	var items []entity.Attraction
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
