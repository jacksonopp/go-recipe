package services

import (
	"context"
	"errors"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserServiceErrorCode int

const (
	UserNotFound UserServiceErrorCode = iota
)

type UserServiceError struct {
	Code UserServiceErrorCode
	Err  error
}

func NewUserServiceError(code UserServiceErrorCode, err error) UserServiceError {
	return UserServiceError{Code: code, Err: err}
}

func (e UserServiceError) Error() string {
	return e.Err.Error()
}

type UserService interface {
	GetUserById(id uint) (*domain.User, error)
	GetUserByUsername(name string) (*domain.User, error)
}

type userService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewUserService(db *gorm.DB) UserService {
	ctx := context.Background()
	return &userService{db: db, ctx: ctx}
}

type userVal struct {
	user *domain.User
	err  error
}

func (s userService) GetUserById(id uint) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	ch := make(chan userVal)

	go func() {
		defer cancel()
		var user domain.User
		err := s.db.Preload("Recipes.Ingredients").
			Preload("Recipes.Instructions").
			Preload(clause.Associations).
			First(&user, id).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := NewUserServiceError(UserNotFound, err)
				ch <- userVal{user: nil, err: err}
				return
			}
			ch <- userVal{user: nil, err: err}
			return
		}
		ch <- userVal{user: &user, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case val := <-ch:
		return val.user, val.err
	}
}

func (s userService) GetUserByUsername(name string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	ch := make(chan userVal)

	go func() {
		defer cancel()
		var user domain.User

		err := s.db.Preload("Recipes.Ingredients").
			Preload("Recipes.Instructions").
			Preload(clause.Associations).
			Where("username = ?", name).
			First(&user).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := NewUserServiceError(UserNotFound, err)
				ch <- userVal{user: nil, err: err}
				return
			}
			ch <- userVal{user: nil, err: err}
			return
		}
		ch <- userVal{user: &user, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case val := <-ch:
		return val.user, val.err
	}
}
