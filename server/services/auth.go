package services

import (
	"errors"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
)

type authService struct {
	db *gorm.DB
}

type AuthService interface {
	CreateUser(user domain.User) error
	GetUserByName(name string) (*domain.User, error)
	LoginUser(name, password string) (*domain.User, error)
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db: db}
}

func (s *authService) CreateUser(user domain.User) error {
	salt, err := genRandStr(32)
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
			return ErrUserAlreadyExists
		}
		return ErrUnknown
	}

	return result.Error
}

func (s *authService) GetUserByName(name string) (*domain.User, error) {
	var user domain.User
	tx := s.db.Where("username = ?", name).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		log.Println("error getting user by name", tx.Error)
		return nil, ErrUnknown
	}
	return &user, nil
}

func (s *authService) LoginUser(name, password string) (*domain.User, error) {
	user, err := s.GetUserByName(name)
	if err != nil {
		log.Println("error getting user by name", err)
		return nil, err
	}

	ok := checkPasswordHash(password, user.Salt, user.Password)
	if !ok {
		log.Println("password mismatch")
		return nil, ErrPasswordMismatch
	}

	user.Salt = ""
	user.Password = ""

	return user, nil
}
