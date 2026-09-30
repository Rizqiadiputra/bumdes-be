package dto

import (
	"time"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type PermissionResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func NewRoleResponse(r entity.Role) RoleResponse {
	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: NewPermissionResponseList(r.Permissions),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func NewRoleResponseList(roles []entity.Role) []RoleResponse {
	result := make([]RoleResponse, 0, len(roles))
	for _, r := range roles {
		result = append(result, NewRoleResponse(r))
	}
	return result
}

func NewPermissionResponse(p entity.Permission) PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func NewPermissionResponseList(permissions []entity.Permission) []PermissionResponse {
	result := make([]PermissionResponse, 0, len(permissions))
	for _, p := range permissions {
		result = append(result, NewPermissionResponse(p))
	}
	return result
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}
