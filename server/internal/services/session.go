package services

import (
	"errors"
	"wout/internal/models"
	"wout/internal/repositories"
)

type SessionService struct {
	SessionRepository  *repositories.SessionRepository
	NotebookRepository *repositories.NotebookRepository
}

func (s *SessionService) CreateSession(session *models.Session) error {
	permission, err := s.NotebookRepository.NotebookHasOwnerId(session.OwnerID, session.NotebookID)
	if err != nil {
		return err
	}
	if !permission {
		return errors.New("user does not have permission to create a session with this notebook")
	}
	return s.SessionRepository.CreateSession(session)
}

func (s *SessionService) GetActiveSessions() ([]models.Session, error) {
	return s.SessionRepository.GetActiveSessions()
}

func (s *SessionService) GetSessionByID(id string) (models.Session, error) {
	return s.SessionRepository.GetSessionByID(id)
}

func (s *SessionService) DeleteSessionByID(id string) error {
	_, err := s.SessionRepository.GetSessionByID(id)
	if err != nil {
		return err
	}
	return s.SessionRepository.DeleteSessionByID(id)
}
