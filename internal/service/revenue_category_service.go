package service

import (
	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type RevenueCategoryService interface {
	ListRevenueCategories(filter repository.RevenueCategoryFilter) ([]dto.RevenueCategoryResponse, error)
	GetRevenueCategory(id uint) (*dto.RevenueCategoryResponse, error)
	CreateRevenueCategory(req dto.CreateRevenueCategoryRequest) (*dto.RevenueCategoryResponse, error)
	UpdateRevenueCategory(id uint, req dto.UpdateRevenueCategoryRequest) (old *dto.RevenueCategoryResponse, updated *dto.RevenueCategoryResponse, err error)
	DeleteRevenueCategory(id uint) error
}

type revenueCategoryService struct {
	repo repository.RevenueCategoryRepository
}

func NewRevenueCategoryService(repo repository.RevenueCategoryRepository) RevenueCategoryService {
	return &revenueCategoryService{repo: repo}
}

func (s *revenueCategoryService) ListRevenueCategories(filter repository.RevenueCategoryFilter) ([]dto.RevenueCategoryResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewRevenueCategoryResponseList(items), nil
}

func (s *revenueCategoryService) GetRevenueCategory(id uint) (*dto.RevenueCategoryResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewRevenueCategoryResponse(*item)
	return &resp, nil
}

func (s *revenueCategoryService) CreateRevenueCategory(req dto.CreateRevenueCategoryRequest) (*dto.RevenueCategoryResponse, error) {
	item := entity.RevenueCategory{
		CategoryName: req.CategoryName,
		Source:       req.Source,
		Status:       req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewRevenueCategoryResponse(item)
	return &resp, nil
}

func (s *revenueCategoryService) UpdateRevenueCategory(id uint, req dto.UpdateRevenueCategoryRequest) (*dto.RevenueCategoryResponse, *dto.RevenueCategoryResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewRevenueCategoryResponse(*existing)

	existing.CategoryName = req.CategoryName
	existing.Source = req.Source
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewRevenueCategoryResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *revenueCategoryService) DeleteRevenueCategory(id uint) error {
	return s.repo.Delete(id)
}
