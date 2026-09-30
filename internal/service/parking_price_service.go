package service

import (
	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type ParkingPriceService interface {
	ListParkingPrices(filter repository.ParkingPriceFilter) ([]dto.ParkingPriceResponse, error)
	GetParkingPrice(id uint) (*dto.ParkingPriceResponse, error)
	CreateParkingPrice(req dto.CreateParkingPriceRequest) (*dto.ParkingPriceResponse, error)
	UpdateParkingPrice(id uint, req dto.UpdateParkingPriceRequest) (old *dto.ParkingPriceResponse, updated *dto.ParkingPriceResponse, err error)
	DeleteParkingPrice(id uint) error
}

type parkingPriceService struct {
	repo repository.ParkingPriceRepository
}

func NewParkingPriceService(repo repository.ParkingPriceRepository) ParkingPriceService {
	return &parkingPriceService{repo: repo}
}

func (s *parkingPriceService) ListParkingPrices(filter repository.ParkingPriceFilter) ([]dto.ParkingPriceResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewParkingPriceResponseList(items), nil
}

func (s *parkingPriceService) GetParkingPrice(id uint) (*dto.ParkingPriceResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewParkingPriceResponse(*item)
	return &resp, nil
}

func (s *parkingPriceService) CreateParkingPrice(req dto.CreateParkingPriceRequest) (*dto.ParkingPriceResponse, error) {
	item := entity.ParkingPrice{
		Name:   req.Name,
		Price:  req.Price,
		Status: req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewParkingPriceResponse(item)
	return &resp, nil
}

func (s *parkingPriceService) UpdateParkingPrice(id uint, req dto.UpdateParkingPriceRequest) (*dto.ParkingPriceResponse, *dto.ParkingPriceResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewParkingPriceResponse(*existing)

	existing.Name = req.Name
	existing.Price = req.Price
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewParkingPriceResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *parkingPriceService) DeleteParkingPrice(id uint) error {
	return s.repo.Delete(id)
}
