package service

import (
	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type PermissionService interface {
	ListPermissions() ([]dto.PermissionResponse, error)
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{permissionRepo: permissionRepo}
}

func (s *permissionService) ListPermissions() ([]dto.PermissionResponse, error) {
	permissions, err := s.permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.NewPermissionResponseList(permissions), nil
}
