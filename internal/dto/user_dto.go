package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type UserResponse struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	IsActive  bool         `json:"is_active"`
	Role      RoleResponse `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func NewUserResponse(u entity.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		IsActive:  u.IsActive,
		Role:      NewRoleResponse(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewUserResponseList(users []entity.User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, NewUserResponse(u))
	}
	return result
}

type CreateAccountRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type UpdateAccountRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	RoleID   uint   `json:"role_id" binding:"required"`
	IsActive bool   `json:"is_active"`
	Password string `json:"password" binding:"omitempty,min=6"`
}
