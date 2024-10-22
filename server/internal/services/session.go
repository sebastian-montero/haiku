package services

import (
	"errors"
	"wout/internal/models"
	"wout/internal/repositories"
	"wout/internal/utils"
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

func (s *SessionService) EndSessionByID(sessionId string, ownerId int) error {
	session, err := s.SessionRepository.GetSessionByID(sessionId)
	if err != nil {
		return err
	}
	if !session.IsActive {
		return errors.New("session is already inactive")
	}

	if session.OwnerID != ownerId {
		return errors.New("user does not have permission to end this session")
	}

	session.IsActive = false

	now := utils.GetStringTime()
	session.EndedAt = &now

	return s.SessionRepository.UpdateSession(&session)
}
