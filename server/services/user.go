package services

import (
	"errors"
	"fmt"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
)

type UserServiceErrorCode int

const (
	ErrUserAlreadyExists UserServiceErrorCode = iota
)

type UserServiceError struct {
	Code UserServiceErrorCode
	Msg  string
}

func (e UserServiceError) Error() string {
	return fmt.Sprintf("user service error: %v", e.Msg)
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user domain.User) error {
	pass, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = pass

	result := s.db.Create(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return UserServiceError{
				Code: ErrUserAlreadyExists,
				Msg:  "user already exists",
			}
		}
		return result.Error
	}

	return result.Error
}

func (s *UserService) GetUserByName(name string) (*domain.User, error) {
	var user domain.User
	tx := s.db.Where("name = ?", name).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}
