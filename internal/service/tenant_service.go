package service

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

const (
	tenantDateLayout        = "2006-01-02"
	defaultTenantRentalType = "Bulanan"

	// RentalTypeSekaliBayar menandakan tenant hanya menyewa untuk satu periode
	// tertentu (sesuai end_periode), bukan sewa berkelanjutan.
	RentalTypeSekaliBayar = "Sekali Bayar"
)

var (
	ErrTenantTypeNotFound      = errors.New("tipe tenant tidak ditemukan")
	ErrKiosLocationNotFound    = errors.New("lokasi kios tidak ditemukan")
	ErrInvalidStartPaymentDate = errors.New("format start_payment tidak valid, gunakan YYYY-MM-DD")
	ErrInvalidEndPeriodeDate   = errors.New("format end_periode tidak valid, gunakan YYYY-MM-DD")
)

type TenantService interface {
	ListTenants(filter repository.TenantFilter) (*dto.PaginatedTenantResponse, error)
	GetTenant(id uuid.UUID) (*dto.TenantResponse, error)
	CreateTenant(req dto.CreateTenantRequest) (*dto.TenantResponse, error)
	UpdateTenant(id uuid.UUID, req dto.UpdateTenantRequest) (old *dto.TenantResponse, updated *dto.TenantResponse, err error)
	DeleteTenant(id uuid.UUID) error
}

type tenantService struct {
	tenantRepo     repository.TenantRepository
	tenantTypeRepo repository.TenantTypeRepository
	kiosLocRepo    repository.KiosLocationRepository
}

func NewTenantService(
	tenantRepo repository.TenantRepository,
	tenantTypeRepo repository.TenantTypeRepository,
	kiosLocRepo repository.KiosLocationRepository,
) TenantService {
	return &tenantService{
		tenantRepo:     tenantRepo,
		tenantTypeRepo: tenantTypeRepo,
		kiosLocRepo:    kiosLocRepo,
	}
}

func (s *tenantService) ListTenants(filter repository.TenantFilter) (*dto.PaginatedTenantResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	items, total, err := s.tenantRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.PaginatedTenantResponse{
		Tenants: dto.NewTenantResponseList(items),
		Meta: dto.PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *tenantService) GetTenant(id uuid.UUID) (*dto.TenantResponse, error) {
	item, err := s.tenantRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewTenantResponse(*item)
	return &resp, nil
}

func (s *tenantService) CreateTenant(req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	tenantTypeID, err := uuid.Parse(req.TenantTypeID)
	if err != nil {
		return nil, ErrTenantTypeNotFound
	}
	kiosLocationID, err := uuid.Parse(req.KiosLocationID)
	if err != nil {
		return nil, ErrKiosLocationNotFound
	}

	tenantType, err := s.tenantTypeRepo.FindByID(tenantTypeID)
	if err != nil {
		return nil, ErrTenantTypeNotFound
	}

	kiosLocation, err := s.kiosLocRepo.FindByID(kiosLocationID)
	if err != nil {
		return nil, ErrKiosLocationNotFound
	}

	startPayment, err := time.Parse(tenantDateLayout, req.StartPayment)
	if err != nil {
		return nil, ErrInvalidStartPaymentDate
	}

	rentalType := req.RentalType
	if rentalType == "" {
		rentalType = defaultTenantRentalType
	}

	endPeriode, err := parseTenantEndPeriode(req.EndPeriode)
	if err != nil {
		return nil, err
	}

	tenant := entity.Tenant{
		Name:           req.Name,
		TenantTypeID:   tenantTypeID,
		Type:           tenantType.Type,
		Phone:          req.Phone,
		Email:          req.Email,
		KiosLocationID: kiosLocationID,
		KiosName:       kiosLocation.Kios,
		Price:          req.Price,
		Status:         req.Status,
		RentalType:     rentalType,
		StartPayment:   startPayment,
		EndPeriode:     endPeriode,
	}

	if err := s.tenantRepo.Create(&tenant); err != nil {
		return nil, err
	}

	if err := s.kiosLocRepo.SetTenant(kiosLocationID, &tenant.ID); err != nil {
		return nil, err
	}

	created, err := s.tenantRepo.FindByID(tenant.ID)
	if err != nil {
		return nil, err
	}

	resp := dto.NewTenantResponse(*created)
	return &resp, nil
}

func (s *tenantService) UpdateTenant(id uuid.UUID, req dto.UpdateTenantRequest) (*dto.TenantResponse, *dto.TenantResponse, error) {
	existing, err := s.tenantRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewTenantResponse(*existing)

	tenantTypeID, err := uuid.Parse(req.TenantTypeID)
	if err != nil {
		return nil, nil, ErrTenantTypeNotFound
	}
	kiosLocationID, err := uuid.Parse(req.KiosLocationID)
	if err != nil {
		return nil, nil, ErrKiosLocationNotFound
	}

	tenantType, err := s.tenantTypeRepo.FindByID(tenantTypeID)
	if err != nil {
		return nil, nil, ErrTenantTypeNotFound
	}

	kiosLocation, err := s.kiosLocRepo.FindByID(kiosLocationID)
	if err != nil {
		return nil, nil, ErrKiosLocationNotFound
	}

	startPayment, err := time.Parse(tenantDateLayout, req.StartPayment)
	if err != nil {
		return nil, nil, ErrInvalidStartPaymentDate
	}

	endPeriode, err := parseTenantEndPeriode(req.EndPeriode)
	if err != nil {
		return nil, nil, err
	}

	previousKiosLocationID := existing.KiosLocationID

	existing.Name = req.Name
	existing.TenantTypeID = tenantTypeID
	existing.Type = tenantType.Type
	existing.Phone = req.Phone
	existing.Email = req.Email
	existing.KiosLocationID = kiosLocationID
	existing.KiosName = kiosLocation.Kios
	existing.Price = req.Price
	existing.Status = req.Status
	existing.RentalType = req.RentalType
	existing.StartPayment = startPayment
	existing.EndPeriode = endPeriode

	if err := s.tenantRepo.Update(existing); err != nil {
		return nil, nil, err
	}

	if previousKiosLocationID != kiosLocationID {
		if err := s.kiosLocRepo.SetTenant(previousKiosLocationID, nil); err != nil {
			return nil, nil, err
		}
		if err := s.kiosLocRepo.SetTenant(kiosLocationID, &existing.ID); err != nil {
			return nil, nil, err
		}
	}

	updated, err := s.tenantRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	newResp := dto.NewTenantResponse(*updated)

	return &oldResp, &newResp, nil
}

func parseTenantEndPeriode(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(tenantDateLayout, raw)
	if err != nil {
		return nil, ErrInvalidEndPeriodeDate
	}

	return &parsed, nil
}

func (s *tenantService) DeleteTenant(id uuid.UUID) error {
	existing, err := s.tenantRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.tenantRepo.Delete(id); err != nil {
		return err
	}

	return s.kiosLocRepo.SetTenant(existing.KiosLocationID, nil)
}
