package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

var ErrAttractionPriceNotFound = errors.New("tarif wahana tidak ditemukan")

type AttractionService interface {
	ListAttractions(filter repository.AttractionFilter) ([]dto.AttractionResponse, error)
	GetAttraction(id uuid.UUID) (*dto.AttractionResponse, error)
	CreateAttraction(req dto.CreateAttractionRequest, actorUserID uuid.UUID) (*dto.AttractionResponse, error)
	UpdateAttraction(id uuid.UUID, req dto.UpdateAttractionRequest) (old *dto.AttractionResponse, updated *dto.AttractionResponse, err error)
	DeleteAttraction(id uuid.UUID) error
}

type attractionService struct {
	attractionRepo      repository.AttractionRepository
	attractionPriceRepo repository.AttractionPriceRepository
	paymentMethodRepo   repository.PaymentMethodRepository
}

func NewAttractionService(
	attractionRepo repository.AttractionRepository,
	attractionPriceRepo repository.AttractionPriceRepository,
	paymentMethodRepo repository.PaymentMethodRepository,
) AttractionService {
	return &attractionService{
		attractionRepo:      attractionRepo,
		attractionPriceRepo: attractionPriceRepo,
		paymentMethodRepo:   paymentMethodRepo,
	}
}

func (s *attractionService) ListAttractions(filter repository.AttractionFilter) ([]dto.AttractionResponse, error) {
	items, err := s.attractionRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewAttractionResponseList(items), nil
}

func (s *attractionService) GetAttraction(id uuid.UUID) (*dto.AttractionResponse, error) {
	item, err := s.attractionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewAttractionResponse(*item)
	return &resp, nil
}

func (s *attractionService) CreateAttraction(req dto.CreateAttractionRequest, actorUserID uuid.UUID) (*dto.AttractionResponse, error) {
	attractionPrice, err := s.attractionPriceRepo.FindByID(req.AttractionPriceID)
	if err != nil {
		return nil, ErrAttractionPriceNotFound
	}

	paymentMethod, err := s.paymentMethodRepo.FindByID(req.PaymentID)
	if err != nil {
		return nil, ErrPaymentMethodNotFound
	}

	attraction := entity.Attraction{
		AttractionPriceID: req.AttractionPriceID,
		AttractionName:    attractionPrice.Name,
		Price:             attractionPrice.Price,
		Count:             req.Count,
		PaymentID:         req.PaymentID,
		PaymentMethodName: paymentMethod.Method,
		UserID:            actorUserID,
		Date:              time.Now(),
	}

	if err := s.attractionRepo.Create(&attraction); err != nil {
		return nil, err
	}

	created, err := s.attractionRepo.FindByID(attraction.ID)
	if err != nil {
		return nil, err
	}

	resp := dto.NewAttractionResponse(*created)
	return &resp, nil
}

func (s *attractionService) UpdateAttraction(id uuid.UUID, req dto.UpdateAttractionRequest) (*dto.AttractionResponse, *dto.AttractionResponse, error) {
	existing, err := s.attractionRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewAttractionResponse(*existing)

	attractionPrice, err := s.attractionPriceRepo.FindByID(req.AttractionPriceID)
	if err != nil {
		return nil, nil, ErrAttractionPriceNotFound
	}

	paymentMethod, err := s.paymentMethodRepo.FindByID(req.PaymentID)
	if err != nil {
		return nil, nil, ErrPaymentMethodNotFound
	}

	existing.AttractionPriceID = req.AttractionPriceID
	existing.AttractionName = attractionPrice.Name
	existing.Price = attractionPrice.Price
	existing.PaymentID = req.PaymentID
	existing.PaymentMethodName = paymentMethod.Method
	existing.Count = req.Count

	if err := s.attractionRepo.Update(existing); err != nil {
		return nil, nil, err
	}

	updated, err := s.attractionRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	newResp := dto.NewAttractionResponse(*updated)

	return &oldResp, &newResp, nil
}

func (s *attractionService) DeleteAttraction(id uuid.UUID) error {
	return s.attractionRepo.Delete(id)
}
