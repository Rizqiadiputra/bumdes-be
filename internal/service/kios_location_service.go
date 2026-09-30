package service

import (
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type KiosLocationService interface {
	ListKiosLocations(filter repository.KiosLocationFilter) ([]dto.KiosLocationResponse, error)
	GetKiosLocation(id uuid.UUID) (*dto.KiosLocationResponse, error)
	CreateKiosLocation(req dto.CreateKiosLocationRequest) (*dto.KiosLocationResponse, error)
	UpdateKiosLocation(id uuid.UUID, req dto.UpdateKiosLocationRequest) (old *dto.KiosLocationResponse, updated *dto.KiosLocationResponse, err error)
	DeleteKiosLocation(id uuid.UUID) error
}

type kiosLocationService struct {
	repo repository.KiosLocationRepository
}

func NewKiosLocationService(repo repository.KiosLocationRepository) KiosLocationService {
	return &kiosLocationService{repo: repo}
}

func (s *kiosLocationService) ListKiosLocations(filter repository.KiosLocationFilter) ([]dto.KiosLocationResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewKiosLocationResponseList(items), nil
}

func (s *kiosLocationService) GetKiosLocation(id uuid.UUID) (*dto.KiosLocationResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewKiosLocationResponse(*item)
	return &resp, nil
}

func (s *kiosLocationService) CreateKiosLocation(req dto.CreateKiosLocationRequest) (*dto.KiosLocationResponse, error) {
	item := entity.KiosLocation{
		Kios:   req.Kios,
		Size:   req.Size,
		Price:  req.Price,
		Status: req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewKiosLocationResponse(item)
	return &resp, nil
}

func (s *kiosLocationService) UpdateKiosLocation(id uuid.UUID, req dto.UpdateKiosLocationRequest) (*dto.KiosLocationResponse, *dto.KiosLocationResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewKiosLocationResponse(*existing)

	existing.Kios = req.Kios
	existing.Size = req.Size
	existing.Price = req.Price
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewKiosLocationResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *kiosLocationService) DeleteKiosLocation(id uuid.UUID) error {
	return s.repo.Delete(id)
}
