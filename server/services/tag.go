package services

import (
	"context"
	"errors"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
)

type TagService interface {
	GetAllTags() ([]*domain.Tag, error)
	CreateTag(tag string) (*domain.Tag, error)
	DeleteTag(id uint) error
}

type tagService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewTagService(db *gorm.DB) TagService {
	ctx := context.Background()
	return &tagService{db: db, ctx: ctx}
}

type tagVal struct {
	tag *domain.Tag
	err error
}

func (s *tagService) GetAllTags() ([]*domain.Tag, error) {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	type tagsVal struct {
		tags []*domain.Tag
		err  error
	}

	ch := make(chan tagsVal)

	go func() {
		defer cancel()
		log.Println("getting all tags")
		tags := []*domain.Tag{}
		err := s.db.Find(&tags).Error
		ch <- tagsVal{tags: tags, err: err}
	}()

	select {
	case v := <-ch:
		return v.tags, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		return nil, ErrTimeoutNoMessage
	}

}

func (s *tagService) CreateTag(tag string) (*domain.Tag, error) {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	ch := make(chan tagVal)

	go func() {
		tx := s.db.Begin()
		defer recoverTx(tx)
		tag := &domain.Tag{Tag: tag}
		err := tx.Create(tag).Error
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				ch <- tagVal{tag: nil, err: ErrTagConflict}
				return
			}
			ch <- tagVal{tag: nil, err: ErrUnknown}
			return
		}
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			ch <- tagVal{tag: nil, err: ErrUnknown}
			return
		}
		ch <- tagVal{tag: tag, err: nil}
	}()

	select {
	case v := <-ch:
		return v.tag, v.err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		return nil, ErrTimeoutNoMessage
	}
}

func (s *tagService) DeleteTag(id uint) error {
	ctx, cancel := context.WithTimeout(s.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	errCh := make(chan error)

	go func() {
		defer cancel()
		tx := s.db.Begin()
		defer recoverTx(tx)

		// check if tag exists
		tag := &domain.Tag{}
		err := tx.First(tag, id).Error
		if err != nil {
			tx.Rollback()
			errCh <- ErrTagNotFound
			return
		}

		//	remove associations
		err = tx.Model(tag).Association("Recipes").Clear()
		if err != nil {
			tx.Rollback()
			errCh <- ErrUnknown
			return
		}

		// delete tag
		err = tx.Delete(tag).Error
		if err != nil {
			tx.Rollback()
			errCh <- ErrUnknown
			return
		}

		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			errCh <- ErrCommit
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ErrTimeout
		}
		return ErrTimeoutNoMessage
	}
}
