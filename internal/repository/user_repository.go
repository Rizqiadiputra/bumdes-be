package repository

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

type UserFilter struct {
	Search   string
	NameRole string
}

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uuid.UUID) (*entity.User, error)
	FindAll(filter UserFilter) ([]entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uuid.UUID) error
	ExistsByEmail(email string, excludeID uuid.UUID) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role.Permissions").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role.Permissions").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(filter UserFilter) ([]entity.User, error) {
	query := r.db.Preload("Role.Permissions").Order("id asc")

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", like, like)
	}

	if filter.NameRole != "" && !strings.EqualFold(filter.NameRole, "all") {
		query = query.Where("role_id IN (?)", r.db.Model(&entity.Role{}).Select("id").Where("name = ?", filter.NameRole))
	}

	var users []entity.User
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) ExistsByEmail(email string, excludeID uuid.UUID) (bool, error) {
	query := r.db.Model(&entity.User{}).Where("email = ?", email)
	if excludeID != uuid.Nil {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
