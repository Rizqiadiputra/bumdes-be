package service

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

const (
	// TenantStatusAktif adalah nilai status tenant yang dianggap masih aktif
	// menyewa, dipakai sebagai filter oleh cron billing.
	TenantStatusAktif = "Aktif"

	BillingStatusUnpaid        = "Unpaid"
	BillingStatusPaid          = "Paid"
	BillingStatusPartiallyPaid = "Partially Paid"
)

type BillingService interface {
	// GenerateDueBillings membuat billing_histories untuk tenant yang
	// tanggal jatuh temponya jatuh persis H-dueReminderDays hari dari
	// referenceDate.
	GenerateDueBillings(referenceDate time.Time) ([]entity.BillingHistory, error)
	ListBillingHistories(filter repository.BillingHistoryFilter) (*dto.PaginatedBillingHistoryResponse, error)
	GetBillingHistory(id uuid.UUID) (*dto.BillingHistoryResponse, error)
}

type billingService struct {
	billingRepo     repository.BillingHistoryRepository
	tenantRepo      repository.TenantRepository
	dueReminderDays int
}

func NewBillingService(
	billingRepo repository.BillingHistoryRepository,
	tenantRepo repository.TenantRepository,
	dueReminderDays int,
) BillingService {
	return &billingService{
		billingRepo:     billingRepo,
		tenantRepo:      tenantRepo,
		dueReminderDays: dueReminderDays,
	}
}

func (s *billingService) GenerateDueBillings(referenceDate time.Time) ([]entity.BillingHistory, error) {
	today := truncateToDate(referenceDate)
	targetDueDate := today.AddDate(0, 0, s.dueReminderDays)

	tenants, err := s.tenantRepo.FindActiveForBilling(TenantStatusAktif, RentalTypeSekaliBayar)
	if err != nil {
		return nil, err
	}

	created := make([]entity.BillingHistory, 0)

	for _, tenant := range tenants {
		dueDate, err := s.nextDueDate(tenant, today)
		if err != nil {
			log.Printf("billing cron: lewati tenant %s: %v", tenant.ID, err)
			continue
		}

		// end_periode masih ada (belum lewat) untuk tanggal jatuh tempo yang dihitung.
		if tenant.EndPeriode != nil && dueDate.After(truncateToDate(*tenant.EndPeriode)) {
			continue
		}

		if !dueDate.Equal(targetDueDate) {
			continue
		}

		exists, err := s.billingRepo.ExistsByTenantAndDateTempo(tenant.ID, dueDate)
		if err != nil {
			return nil, err
		}
		if exists {
			continue
		}

		billing := entity.BillingHistory{
			TenantID:           tenant.ID,
			KiosID:             tenant.KiosLocationID,
			RentalType:         tenant.RentalType,
			DateTempo:          dueDate,
			TotalAmountPayable: tenant.Price,
			Status:             BillingStatusUnpaid,
		}

		if err := s.billingRepo.Create(&billing); err != nil {
			return nil, err
		}

		created = append(created, billing)
	}

	return created, nil
}

func (s *billingService) ListBillingHistories(filter repository.BillingHistoryFilter) (*dto.PaginatedBillingHistoryResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	items, total, err := s.billingRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.PaginatedBillingHistoryResponse{
		BillingHistories: dto.NewBillingHistoryResponseList(items),
		Meta: dto.PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *billingService) GetBillingHistory(id uuid.UUID) (*dto.BillingHistoryResponse, error) {
	item, err := s.billingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewBillingHistoryResponse(*item)
	return &resp, nil
}

// nextDueDate menghitung tanggal jatuh tempo berikutnya untuk tenant:
// mengambil date_tempo billing terakhir (bila ada) sebagai basis, atau
// start_payment bila tenant belum pernah ditagih, lalu maju satu periode
// sewa demi satu periode sampai tanggal tersebut jatuh pada hari ini atau
// setelahnya.
func (s *billingService) nextDueDate(tenant entity.Tenant, today time.Time) (time.Time, error) {
	base := truncateToDate(tenant.StartPayment)

	latest, err := s.billingRepo.FindLatestByTenant(tenant.ID)
	switch {
	case err == nil:
		base = truncateToDate(latest.DateTempo)
	case errors.Is(err, gorm.ErrRecordNotFound):
		// tenant belum pernah ditagih, mulai dari start_payment.
	default:
		return time.Time{}, err
	}

	next := addRentalInterval(base, tenant.RentalType)
	for next.Before(today) {
		next = addRentalInterval(next, tenant.RentalType)
	}

	return next, nil
}

func addRentalInterval(date time.Time, rentalType string) time.Time {
	switch rentalType {
	case "Mingguan":
		return date.AddDate(0, 0, 7)
	case "Tahunan":
		return date.AddDate(1, 0, 0)
	default: // "Bulanan" dan tipe lain yang tidak dikenal dianggap bulanan.
		return date.AddDate(0, 1, 0)
	}
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
