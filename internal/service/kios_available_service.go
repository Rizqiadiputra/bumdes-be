package service

import (
	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type KiosAvailableService interface {
	ListKiosAvailables() ([]dto.KiosAvailableResponse, error)
}

type kiosAvailableService struct {
	kiosLocationRepo repository.KiosLocationRepository
	tenantRepo       repository.TenantRepository
}

func NewKiosAvailableService(kiosLocationRepo repository.KiosLocationRepository, tenantRepo repository.TenantRepository) KiosAvailableService {
	return &kiosAvailableService{kiosLocationRepo: kiosLocationRepo, tenantRepo: tenantRepo}
}

func (s *kiosAvailableService) ListKiosAvailables() ([]dto.KiosAvailableResponse, error) {
	kioses, err := s.kiosLocationRepo.FindAll(repository.KiosLocationFilter{})
	if err != nil {
		return nil, err
	}

	result := make([]dto.KiosAvailableResponse, 0, len(kioses))
	for _, k := range kioses {
		price := k.Price

		if k.TenantID != nil {
			tenant, err := s.tenantRepo.FindByID(*k.TenantID)
			if err == nil {
				price = tenant.Price
			}
		}

		result = append(result, dto.NewKiosAvailableResponse(k, price))
	}

	return result, nil
}
