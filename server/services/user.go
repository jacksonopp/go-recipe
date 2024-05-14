package services

import (
	"errors"
	"fmt"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
)

type UserServiceErrorCode int

const (
	ErrUnknown UserServiceErrorCode = iota
	ErrUserAlreadyExists
	ErrUserNotFound
	ErrPasswordMismatch
)

type UserServiceError struct {
	Code UserServiceErrorCode
	Msg  string
}

func NewUserServiceError(code UserServiceErrorCode, msg string) UserServiceError {
	return UserServiceError{
		Code: code,
		Msg:  msg,
	}
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
			return NewUserServiceError(ErrUserAlreadyExists, "user already exists")
		}
		return NewUserServiceError(ErrUnknown, "unknown error")
	}

	return result.Error
}

func (s *UserService) GetUserByName(name string) (*domain.User, error) {
	var user domain.User
	tx := s.db.Where("name = ?", name).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, NewUserServiceError(ErrUserNotFound, "user not found")
		}
		return nil, NewUserServiceError(ErrUnknown, "unknown error")
	}
	return &user, nil
}

func (s *UserService) LoginUser(name, password string) (*domain.User, error) {
	user, err := s.GetUserByName(name)
	if err != nil {
		return nil, err
	}

	ok := checkPasswordHash(password, user.Password)
	if !ok {
		return nil, NewUserServiceError(ErrPasswordMismatch, "passwords do not match")
	}

	return user, nil
}
