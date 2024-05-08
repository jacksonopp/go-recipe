package services

import (
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
)

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

	return result.Error
}
