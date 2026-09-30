package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

var (
	ErrTicketPriceNotFound = errors.New("tarif tiket tidak ditemukan")
)

type TicketingService interface {
	ListTicketings(filter repository.TicketingFilter) ([]dto.TicketingResponse, error)
	GetTicketing(id uuid.UUID) (*dto.TicketingResponse, error)
	CreateTicketing(req dto.CreateTicketingRequest, actorUserID uuid.UUID) (*dto.TicketingResponse, error)
	UpdateTicketing(id uuid.UUID, req dto.UpdateTicketingRequest) (old *dto.TicketingResponse, updated *dto.TicketingResponse, err error)
	DeleteTicketing(id uuid.UUID) error
}

type ticketingService struct {
	ticketingRepo     repository.TicketingRepository
	ticketPriceRepo   repository.TicketPriceRepository
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewTicketingService(
	ticketingRepo repository.TicketingRepository,
	ticketPriceRepo repository.TicketPriceRepository,
	paymentMethodRepo repository.PaymentMethodRepository,
) TicketingService {
	return &ticketingService{
		ticketingRepo:     ticketingRepo,
		ticketPriceRepo:   ticketPriceRepo,
		paymentMethodRepo: paymentMethodRepo,
	}
}

func (s *ticketingService) ListTicketings(filter repository.TicketingFilter) ([]dto.TicketingResponse, error) {
	items, err := s.ticketingRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewTicketingResponseList(items), nil
}

func (s *ticketingService) GetTicketing(id uuid.UUID) (*dto.TicketingResponse, error) {
	item, err := s.ticketingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewTicketingResponse(*item)
	return &resp, nil
}

func (s *ticketingService) CreateTicketing(req dto.CreateTicketingRequest, actorUserID uuid.UUID) (*dto.TicketingResponse, error) {
	ticketPrice, err := s.ticketPriceRepo.FindByID(req.TicketPriceID)
	if err != nil {
		return nil, ErrTicketPriceNotFound
	}

	paymentMethod, err := s.paymentMethodRepo.FindByID(req.PaymentID)
	if err != nil {
		return nil, ErrPaymentMethodNotFound
	}

	ticketing := entity.Ticketing{
		TicketPriceID:     req.TicketPriceID,
		TypeName:          ticketPrice.TypeName,
		Price:             ticketPrice.Price,
		PaymentID:         req.PaymentID,
		PaymentMethodName: paymentMethod.Method,
		Count:             req.Count,
		UserID:            actorUserID,
		Date:              time.Now(),
	}

	if err := s.ticketingRepo.Create(&ticketing); err != nil {
		return nil, err
	}

	created, err := s.ticketingRepo.FindByID(ticketing.ID)
	if err != nil {
		return nil, err
	}

	resp := dto.NewTicketingResponse(*created)
	return &resp, nil
}

func (s *ticketingService) UpdateTicketing(id uuid.UUID, req dto.UpdateTicketingRequest) (*dto.TicketingResponse, *dto.TicketingResponse, error) {
	existing, err := s.ticketingRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewTicketingResponse(*existing)

	ticketPrice, err := s.ticketPriceRepo.FindByID(req.TicketPriceID)
	if err != nil {
		return nil, nil, ErrTicketPriceNotFound
	}

	paymentMethod, err := s.paymentMethodRepo.FindByID(req.PaymentID)
	if err != nil {
		return nil, nil, ErrPaymentMethodNotFound
	}

	existing.TicketPriceID = req.TicketPriceID
	existing.TypeName = ticketPrice.TypeName
	existing.Price = ticketPrice.Price
	existing.PaymentID = req.PaymentID
	existing.PaymentMethodName = paymentMethod.Method
	existing.Count = req.Count

	if err := s.ticketingRepo.Update(existing); err != nil {
		return nil, nil, err
	}

	updated, err := s.ticketingRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	newResp := dto.NewTicketingResponse(*updated)

	return &oldResp, &newResp, nil
}

func (s *ticketingService) DeleteTicketing(id uuid.UUID) error {
	return s.ticketingRepo.Delete(id)
}
