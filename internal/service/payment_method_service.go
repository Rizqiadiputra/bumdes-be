package service

import (
	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type PaymentMethodService interface {
	ListPaymentMethods(filter repository.PaymentMethodFilter) ([]dto.PaymentMethodResponse, error)
	GetPaymentMethod(id uint) (*dto.PaymentMethodResponse, error)
	CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
	UpdatePaymentMethod(id uint, req dto.UpdatePaymentMethodRequest) (old *dto.PaymentMethodResponse, updated *dto.PaymentMethodResponse, err error)
	DeletePaymentMethod(id uint) error
}

type paymentMethodService struct {
	repo repository.PaymentMethodRepository
}

func NewPaymentMethodService(repo repository.PaymentMethodRepository) PaymentMethodService {
	return &paymentMethodService{repo: repo}
}

func (s *paymentMethodService) ListPaymentMethods(filter repository.PaymentMethodFilter) ([]dto.PaymentMethodResponse, error) {
	items, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewPaymentMethodResponseList(items), nil
}

func (s *paymentMethodService) GetPaymentMethod(id uint) (*dto.PaymentMethodResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewPaymentMethodResponse(*item)
	return &resp, nil
}

func (s *paymentMethodService) CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	item := entity.PaymentMethod{
		Method: req.Method,
		Status: req.Status,
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	resp := dto.NewPaymentMethodResponse(item)
	return &resp, nil
}

func (s *paymentMethodService) UpdatePaymentMethod(id uint, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, *dto.PaymentMethodResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewPaymentMethodResponse(*existing)

	existing.Method = req.Method
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, nil, err
	}

	newResp := dto.NewPaymentMethodResponse(*existing)
	return &oldResp, &newResp, nil
}

func (s *paymentMethodService) DeletePaymentMethod(id uint) error {
	return s.repo.Delete(id)
}
