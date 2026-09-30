package service

import (
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type TicketPriceService interface {
	ListTicketPrices(filter repository.TicketPriceFilter) ([]dto.TicketPriceResponse, error)
	GetTicketPrice(id uuid.UUID) (*dto.TicketPriceResponse, error)
	CreateTicketPrice(req dto.CreateTicketPriceRequest) (*dto.TicketPriceResponse, error)
	UpdateTicketPrice(id uuid.UUID, req dto.UpdateTicketPriceRequest) (old *dto.TicketPriceResponse, updated *dto.TicketPriceResponse, err error)
	DeleteTicketPrice(id uuid.UUID) error
}

type ticketPriceService struct {
	repo repository.TicketPriceRepository
}

func NewTicketPriceService(repo repository.TicketPriceRepository) TicketPriceService {
	return &ticketPriceService{repo: repo}
}

func (s *ticketPriceService) ListTicketPrices(filter repository.TicketPriceFilter) ([]dto.TicketPriceResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewTicketPriceResponseList(items), nil
}

func (s *ticketPriceService) GetTicketPrice(id uuid.UUID) (*dto.TicketPriceResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewTicketPriceResponse(*item)
	return &resp, nil
}

func (s *ticketPriceService) CreateTicketPrice(req dto.CreateTicketPriceRequest) (*dto.TicketPriceResponse, error) {
	item := entity.TicketPrice{
		TypeName: req.TypeName,
		Price:    req.Price,
		Status:   req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewTicketPriceResponse(item)
	return &resp, nil
}

func (s *ticketPriceService) UpdateTicketPrice(id uuid.UUID, req dto.UpdateTicketPriceRequest) (*dto.TicketPriceResponse, *dto.TicketPriceResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewTicketPriceResponse(*existing)

	existing.TypeName = req.TypeName
	existing.Price = req.Price
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewTicketPriceResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *ticketPriceService) DeleteTicketPrice(id uuid.UUID) error {
	return s.repo.Delete(id)
}
