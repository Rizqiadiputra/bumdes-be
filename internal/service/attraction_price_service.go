package service

import (
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type AttractionPriceService interface {
	ListAttractionPrices(filter repository.AttractionPriceFilter) ([]dto.AttractionPriceResponse, error)
	GetAttractionPrice(id uuid.UUID) (*dto.AttractionPriceResponse, error)
	CreateAttractionPrice(req dto.CreateAttractionPriceRequest) (*dto.AttractionPriceResponse, error)
	UpdateAttractionPrice(id uuid.UUID, req dto.UpdateAttractionPriceRequest) (old *dto.AttractionPriceResponse, updated *dto.AttractionPriceResponse, err error)
	DeleteAttractionPrice(id uuid.UUID) error
}

type attractionPriceService struct {
	repo repository.AttractionPriceRepository
}

func NewAttractionPriceService(repo repository.AttractionPriceRepository) AttractionPriceService {
	return &attractionPriceService{repo: repo}
}

func (s *attractionPriceService) ListAttractionPrices(filter repository.AttractionPriceFilter) ([]dto.AttractionPriceResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewAttractionPriceResponseList(items), nil
}

func (s *attractionPriceService) GetAttractionPrice(id uuid.UUID) (*dto.AttractionPriceResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewAttractionPriceResponse(*item)
	return &resp, nil
}

func (s *attractionPriceService) CreateAttractionPrice(req dto.CreateAttractionPriceRequest) (*dto.AttractionPriceResponse, error) {
	item := entity.AttractionPrice{
		Name:   req.Name,
		Price:  req.Price,
		Status: req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewAttractionPriceResponse(item)
	return &resp, nil
}

func (s *attractionPriceService) UpdateAttractionPrice(id uuid.UUID, req dto.UpdateAttractionPriceRequest) (*dto.AttractionPriceResponse, *dto.AttractionPriceResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewAttractionPriceResponse(*existing)

	existing.Name = req.Name
	existing.Price = req.Price
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewAttractionPriceResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *attractionPriceService) DeleteAttractionPrice(id uuid.UUID) error {
	return s.repo.Delete(id)
}
