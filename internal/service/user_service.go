package service

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/entity"
	"github.com/liyansasongko/bumdes-be/internal/repository"
)

var ErrEmailAlreadyExists = errors.New("email sudah digunakan")

type UserService interface {
	GetMe(userID uuid.UUID) (*dto.UserResponse, error)
	ListAccounts(filter repository.UserFilter) ([]dto.UserResponse, error)
	CreateAccount(req dto.CreateAccountRequest) (*dto.UserResponse, error)
	UpdateAccount(id uuid.UUID, req dto.UpdateAccountRequest) (old *dto.UserResponse, updated *dto.UserResponse, err error)
	DeleteAccount(id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetMe(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	resp := dto.NewUserResponse(*user)
	return &resp, nil
}

func (s *userService) ListAccounts(filter repository.UserFilter) ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponseList(users), nil
}

func (s *userService) CreateAccount(req dto.CreateAccountRequest) (*dto.UserResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(req.Email, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		RoleID:   req.RoleID,
		IsActive: true,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	created, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	resp := dto.NewUserResponse(*created)
	return &resp, nil
}

func (s *userService) UpdateAccount(id uuid.UUID, req dto.UpdateAccountRequest) (*dto.UserResponse, *dto.UserResponse, error) {
	existing, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	oldResp := dto.NewUserResponse(*existing)

	exists, err := s.userRepo.ExistsByEmail(req.Email, id)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrEmailAlreadyExists
	}

	existing.Name = req.Name
	existing.Email = req.Email
	existing.RoleID = req.RoleID
	existing.IsActive = req.IsActive

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}
		existing.Password = string(hashed)
	}

	if err := s.userRepo.Update(existing); err != nil {
		return nil, nil, err
	}

	updated, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	newResp := dto.NewUserResponse(*updated)

	return &oldResp, &newResp, nil
}

func (s *userService) DeleteAccount(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}
