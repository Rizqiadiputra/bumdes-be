package service

import (
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type TenantTypeService interface {
	ListTenantTypes(filter repository.TenantTypeFilter) ([]dto.TenantTypeResponse, error)
	GetTenantType(id uuid.UUID) (*dto.TenantTypeResponse, error)
	CreateTenantType(req dto.CreateTenantTypeRequest) (*dto.TenantTypeResponse, error)
	UpdateTenantType(id uuid.UUID, req dto.UpdateTenantTypeRequest) (old *dto.TenantTypeResponse, updated *dto.TenantTypeResponse, err error)
	DeleteTenantType(id uuid.UUID) error
}

type tenantTypeService struct {
	repo repository.TenantTypeRepository
}

func NewTenantTypeService(repo repository.TenantTypeRepository) TenantTypeService {
	return &tenantTypeService{repo: repo}
}

func (s *tenantTypeService) ListTenantTypes(filter repository.TenantTypeFilter) ([]dto.TenantTypeResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewTenantTypeResponseList(items), nil
}

func (s *tenantTypeService) GetTenantType(id uuid.UUID) (*dto.TenantTypeResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewTenantTypeResponse(*item)
	return &resp, nil
}

func (s *tenantTypeService) CreateTenantType(req dto.CreateTenantTypeRequest) (*dto.TenantTypeResponse, error) {
	item := entity.TenantType{
		Type:        req.Type,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewTenantTypeResponse(item)
	return &resp, nil
}

func (s *tenantTypeService) UpdateTenantType(id uuid.UUID, req dto.UpdateTenantTypeRequest) (*dto.TenantTypeResponse, *dto.TenantTypeResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewTenantTypeResponse(*existing)

	existing.Type = req.Type
	existing.Description = req.Description
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewTenantTypeResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *tenantTypeService) DeleteTenantType(id uuid.UUID) error {
	return s.repo.Delete(id)
}
