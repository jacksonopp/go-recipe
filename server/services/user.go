package services

import (
	"errors"
	"fmt"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
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

type userService struct {
	db *gorm.DB
}

type UserService interface {
	CreateUser(user domain.User) error
	GetUserByName(name string) (*domain.User, error)
	LoginUser(name, password string) (*domain.User, error)
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) CreateUser(user domain.User) error {
	salt, err := generateSalt(32)
	if err != nil {
		return err
	}
	user.Salt = salt

	pass, err := hashPassword(user.Password, user.Salt)
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

func (s *userService) GetUserByName(name string) (*domain.User, error) {
	var user domain.User
	tx := s.db.Where("username = ?", name).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, NewUserServiceError(ErrUserNotFound, "user not found")
		}
		log.Println("error getting user by name", tx.Error)
		return nil, NewUserServiceError(ErrUnknown, "unknown error")
	}
	return &user, nil
}

func (s *userService) LoginUser(name, password string) (*domain.User, error) {
	user, err := s.GetUserByName(name)
	if err != nil {
		log.Println("error getting user by name", err)
		return nil, err
	}

	ok := checkPasswordHash(password, user.Salt, user.Password)
	if !ok {
		log.Println("password mismatch")
		return nil, NewUserServiceError(ErrPasswordMismatch, "passwords do not match")
	}

	user.Salt = ""
	user.Password = ""

	return user, nil
}
