package service

import (
	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

type RoleService interface {
	ListRoles() ([]dto.RoleResponse, error)
	GetRole(id uint) (*dto.RoleResponse, error)
	UpdatePermissions(roleID uint, req dto.UpdateRolePermissionsRequest) (old *dto.RoleResponse, updated *dto.RoleResponse, err error)
}

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{roleRepo: roleRepo, permissionRepo: permissionRepo}
}

func (s *roleService) ListRoles() ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.NewRoleResponseList(roles), nil
}

func (s *roleService) GetRole(id uint) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.NewRoleResponse(*role)
	return &resp, nil
}

func (s *roleService) UpdatePermissions(roleID uint, req dto.UpdateRolePermissionsRequest) (*dto.RoleResponse, *dto.RoleResponse, error) {
	existing, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewRoleResponse(*existing)

	permissions, err := s.permissionRepo.FindByIDs(req.PermissionIDs)
	if err != nil {
		return nil, nil, err
	}

	if err := s.roleRepo.UpdatePermissions(roleID, permissions); err != nil {
		return nil, nil, err
	}

	updated, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, nil, err
	}
	newResp := dto.NewRoleResponse(*updated)

	return &oldResp, &newResp, nil
}
