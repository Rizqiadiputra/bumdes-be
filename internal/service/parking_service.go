package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

var (
	ErrParkingPriceNotFound  = errors.New("tarif parkir tidak ditemukan")
	ErrPaymentMethodNotFound = errors.New("metode pembayaran tidak ditemukan")
)

type ParkingService interface {
	ListParkings(filter repository.ParkingFilter) ([]dto.ParkingResponse, error)
	GetParking(id uuid.UUID) (*dto.ParkingResponse, error)
	CreateParking(req dto.CreateParkingRequest, actorUserID uuid.UUID) (*dto.ParkingResponse, error)
	DeleteParking(id uuid.UUID) error
}

type parkingService struct {
	parkingRepo       repository.ParkingRepository
	parkingPriceRepo  repository.ParkingPriceRepository
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewParkingService(
	parkingRepo repository.ParkingRepository,
	parkingPriceRepo repository.ParkingPriceRepository,
	paymentMethodRepo repository.PaymentMethodRepository,
) ParkingService {
	return &parkingService{
		parkingRepo:       parkingRepo,
		parkingPriceRepo:  parkingPriceRepo,
		paymentMethodRepo: paymentMethodRepo,
	}
}

func (s *parkingService) ListParkings(filter repository.ParkingFilter) ([]dto.ParkingResponse, error) {
	items, err := s.parkingRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewParkingResponseList(items), nil
}

func (s *parkingService) GetParking(id uuid.UUID) (*dto.ParkingResponse, error) {
	item, err := s.parkingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewParkingResponse(*item)
	return &resp, nil
}

func (s *parkingService) CreateParking(req dto.CreateParkingRequest, actorUserID uuid.UUID) (*dto.ParkingResponse, error) {
	parkingPrice, err := s.parkingPriceRepo.FindByID(req.ParkingPriceID)
	if err != nil {
		return nil, ErrParkingPriceNotFound
	}

	if _, err := s.paymentMethodRepo.FindByID(req.PaymentMethodID); err != nil {
		return nil, ErrPaymentMethodNotFound
	}

	parking := entity.Parking{
		Type:            parkingPrice.Name,
		Location:        req.Location,
		Count:           req.Count,
		PaymentMethodID: req.PaymentMethodID,
		Price:           parkingPrice.Price * int64(req.Count),
		UserID:          actorUserID,
		Status:          req.Status,
	}

	if err := s.parkingRepo.Create(&parking); err != nil {
		return nil, err
	}

	created, err := s.parkingRepo.FindByID(parking.ID)
	if err != nil {
		return nil, err
	}

	resp := dto.NewParkingResponse(*created)
	return &resp, nil
}

func (s *parkingService) DeleteParking(id uuid.UUID) error {
	return s.parkingRepo.Delete(id)
}
