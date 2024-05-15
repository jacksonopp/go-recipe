package services

import (
	"errors"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"time"
)

type SessionServiceErrorCode int

const (
	ErrUnknownSession SessionServiceErrorCode = iota
	ErrSessionNotFound
	ErrSessionExpired
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) SessionService {
	return SessionService{db: db}
}

type SessionServiceError struct {
	Code SessionServiceErrorCode
	Msg  string
}

func NewSessionServiceError(code SessionServiceErrorCode, msg string) SessionServiceError {
	return SessionServiceError{
		Code: code,
		Msg:  msg,
	}
}

func (e SessionServiceError) Error() string {
	return e.Msg
}

func (s *SessionService) CreateSession(userID uint) (string, error) {
	token, err := genRandStr(32)
	if err != nil {
		return "", err
	}

	session := domain.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	res := s.db.Create(&session)
	if res.Error != nil {
		return "", res.Error
	}

	return token, nil
}

func (s *SessionService) CheckSession(token string) error {
	var session domain.Session
	res := s.db.First(&session, "token = ?", token)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return NewSessionServiceError(ErrSessionNotFound, "session not found")
		}
		return NewSessionServiceError(ErrUnknownSession, res.Error.Error())
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := s.DeleteSessionByToken(token)
		if err != nil {
			return err
		}
		return NewSessionServiceError(ErrSessionExpired, "session expired")
	}

	return nil
}

func (s *SessionService) DeleteSessionByToken(token string) error {
	res := s.db.Delete(&domain.Session{}, "token = ?", token)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return NewSessionServiceError(ErrSessionNotFound, "session not found")
		}
		return NewSessionServiceError(ErrUnknownSession, res.Error.Error())
	}
	return nil
}
