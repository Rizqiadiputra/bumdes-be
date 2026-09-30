package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type UserLogResponse struct {
	ID          uint            `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	UserName    string          `json:"user_name,omitempty"`
	Action      string          `json:"action"`
	Module      string          `json:"module"`
	Description string          `json:"description"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	IPAddress   string          `json:"ip_address"`
	UserAgent   string          `json:"user_agent"`
	Data        json.RawMessage `json:"data" swaggertype:"object"`
	CreatedAt   time.Time       `json:"created_at"`
}

func NewUserLogResponse(l entity.UserLog) UserLogResponse {
	return UserLogResponse{
		ID:          l.ID,
		UserID:      l.UserID,
		UserName:    l.User.Name,
		Action:      l.Action,
		Module:      l.Module,
		Description: l.Description,
		Method:      l.Method,
		Path:        l.Path,
		IPAddress:   l.IPAddress,
		UserAgent:   l.UserAgent,
		Data:        json.RawMessage(l.Data),
		CreatedAt:   l.CreatedAt,
	}
}

func NewUserLogResponseList(logs []entity.UserLog) []UserLogResponse {
	result := make([]UserLogResponse, 0, len(logs))
	for _, l := range logs {
		result = append(result, NewUserLogResponse(l))
	}
	return result
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedUserLogResponse struct {
	Logs []UserLogResponse `json:"logs"`
	Meta PaginationMeta    `json:"meta"`
}
