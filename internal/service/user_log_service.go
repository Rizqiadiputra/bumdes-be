package service

import (
	"encoding/json"
	"math"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type UserLogService interface {
	LogView(meta utils.RequestMeta, userID uuid.UUID, module, description string) error
	LogCreate(meta utils.RequestMeta, userID uuid.UUID, module, description string, input interface{}) error
	LogUpdate(meta utils.RequestMeta, userID uuid.UUID, module, description string, oldData, newData interface{}) error
	LogDelete(meta utils.RequestMeta, userID uuid.UUID, module, description string) error
	ListLogs(filter repository.UserLogFilter) (*dto.PaginatedUserLogResponse, error)
}

type userLogService struct {
	logRepo repository.UserLogRepository
}

func NewUserLogService(logRepo repository.UserLogRepository) UserLogService {
	return &userLogService{logRepo: logRepo}
}

func (s *userLogService) LogView(meta utils.RequestMeta, userID uuid.UUID, module, description string) error {
	return s.record(meta, userID, entity.LogActionView, module, description, nil)
}

func (s *userLogService) LogCreate(meta utils.RequestMeta, userID uuid.UUID, module, description string, input interface{}) error {
	data := map[string]interface{}{"input": input}
	return s.record(meta, userID, entity.LogActionCreate, module, description, data)
}

func (s *userLogService) LogUpdate(meta utils.RequestMeta, userID uuid.UUID, module, description string, oldData, newData interface{}) error {
	data := map[string]interface{}{"old": oldData, "new": newData}
	return s.record(meta, userID, entity.LogActionUpdate, module, description, data)
}

func (s *userLogService) LogDelete(meta utils.RequestMeta, userID uuid.UUID, module, description string) error {
	return s.record(meta, userID, entity.LogActionDelete, module, description, nil)
}

func (s *userLogService) record(meta utils.RequestMeta, userID uuid.UUID, action, module, description string, data interface{}) error {
	var jsonData datatypes.JSON
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		jsonData = datatypes.JSON(b)
	}

	log := entity.UserLog{
		UserID:      userID,
		Action:      action,
		Module:      module,
		Description: description,
		Method:      meta.Method,
		Path:        meta.Path,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
		Data:        jsonData,
	}

	return s.logRepo.Create(&log)
}

func (s *userLogService) ListLogs(filter repository.UserLogFilter) (*dto.PaginatedUserLogResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	logs, total, err := s.logRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.PaginatedUserLogResponse{
		Logs: dto.NewUserLogResponseList(logs),
		Meta: dto.PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
