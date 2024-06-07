package services

import (
	"context"
	"errors"
	"github.com/jacksonopp/go-recipe/db"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type UserService interface {
	GetUserById(id uint) (*domain.User, error)
	GetUserByUsername(name string) (*domain.User, error)
	GetUsersRecipes(name string, page, limit int) ([]domain.Recipe, error)
	GetUserFiles(name string, _, _ int) ([]domain.FileDto, error)
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
				err := ErrUserNotFound
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
				err := ErrUserNotFound
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

func (s userService) GetUsersRecipes(name string, page, limit int) ([]domain.Recipe, error) {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	type recipesVal struct {
		recipes []domain.Recipe
		err     error
	}

	ch := make(chan recipesVal)

	go func() {
		defer cancel()

		log.Printf("page %d, limit %d", page, limit)

		var recipes []domain.Recipe
		err := s.db.
			Scopes(db.Paginate(page, limit)).
			Preload("Ingredients").
			Preload("Instructions").
			Preload(clause.Associations).
			Joins("JOIN users ON users.id = recipes.user_id").
			Where("users.username = ?", name).
			Find(&recipes).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := ErrUserNotFound
				ch <- recipesVal{recipes: nil, err: err}
				return
			}
			ch <- recipesVal{recipes: nil, err: err}
			return
		}
		ch <- recipesVal{recipes: recipes, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case val := <-ch:
		return val.recipes, val.err
	}
}

func (s userService) GetUserFiles(name string, _, _ int) ([]domain.FileDto, error) {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	type filesVal struct {
		files []domain.FileDto
		err   error
	}

	ch := make(chan filesVal)

	go func() {
		defer cancel()

		var user domain.User

		err := s.db.Where("username = ?", name).
			Preload("Files").
			First(&user).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := ErrUserNotFound
				ch <- filesVal{files: nil, err: err}
				return
			}
			ch <- filesVal{files: nil, err: err}
			return
		}

		ch <- filesVal{files: user.GetFiles(), err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case val := <-ch:
		return val.files, val.err
	}
}
