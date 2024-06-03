package services

import (
	"errors"
	"github.com/jacksonopp/go-recipe/domain"
	"gorm.io/gorm"
	"log"
	"time"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) SessionService {
	return SessionService{db: db}
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
			return ErrSessionNotFound
		}
		return ErrUnknown
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := s.DeleteSessionByToken(token)
		if err != nil {
			return err
		}
		return ErrSessionExpired
	}

	return nil
}

func (s *SessionService) DeleteSessionByToken(token string) error {
	res := s.db.Delete(&domain.Session{}, "token = ?", token)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrSessionNotFound
		}
		return ErrUnknown
	}
	return nil
}

func (s *SessionService) PruneSessions() error {
	res := s.db.Delete(&domain.Session{}, "expires_at < ?", time.Now())
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *SessionService) PruneOnSchedule(t time.Duration) (chan<- bool, error) {
	done := make(chan bool)
	ticker := time.NewTicker(t)

	go func() {
		for {
			select {
			case <-done:
				log.Printf("stopping prune job\n")
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("Pruning sessions")
				if err := s.PruneSessions(); err != nil {
					log.Printf("error pruning session: %v\n", err)
					done <- true
				}
			}
		}
	}()

	return done, nil
}
