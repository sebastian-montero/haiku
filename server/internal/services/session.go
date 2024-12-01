package services

import (
	"errors"
	"fmt"
	"haiku/internal/models"
	"haiku/internal/repositories"
	"haiku/internal/utils"
	"haiku/internal/utils/logger"
)

type SessionService struct {
	SessionRepository  *repositories.SessionRepository
	NotebookRepository *repositories.NotebookRepository
	ContentRepository  *repositories.ContentRepository
}

func (s *SessionService) CreateSession(session *models.Session) error {
	fmt.Println(fmt.Sprintf("Creating session with ownerID %v and notebookID %v", session.OwnerID, session.NotebookID))
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

func (s *SessionService) SessionExistsByNotebookID(notebookID string) (models.Session, bool) {
	return s.SessionRepository.SessionExistsByNotebookID(notebookID)
}

func (s *SessionService) ActivateSession(sessionID string) error {
	session, err := s.SessionRepository.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if session.IsActive {
		logger.Info("Session is already active")
	}

	session.IsActive = true
	return s.SessionRepository.UpdateSession(&session)
}

func (s *SessionService) DeleteSessionByID(id string) error {
	_, err := s.SessionRepository.GetSessionByID(id)
	if err != nil {
		return err
	}
	return s.SessionRepository.DeleteSessionByID(id)
}

func (s *SessionService) EndSessionByID(sessionID string, ownerID int) error {
	session, err := s.SessionRepository.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if !session.IsActive {
		return errors.New("session is already inactive")
	}

	if session.OwnerID != ownerID {
		return errors.New("user does not have permission to end this session")
	}

	session.IsActive = false

	now := utils.GetStringTime()
	session.EndedAt = &now

	return s.SessionRepository.UpdateSession(&session)
}

func (s *SessionService) GetNotebookByID(id string) (models.Notebook, error) {
	return s.NotebookRepository.GetNotebookByID(id)
}

func (s *SessionService) UpdateNotebookContent(notebookID string, content string) error {
	notebook, err := s.NotebookRepository.GetNotebookByID(notebookID)
	if err != nil {
		return err
	}

	notebook.LatestContent = content

	return s.NotebookRepository.UpdateNotebook(&notebook)
	// update UpdateNotebookContent

}

func (s *SessionService) CreateContent(notebookID string, msg string) error {
	session, err := s.SessionRepository.GetSessionByNotebookID(notebookID)
	if err != nil {
		return err
	}

	content := models.Content{
		SessionID: session.ID,
		Content:   msg,
	}
	return s.ContentRepository.CreateContent(&content)
}

func (s *SessionService) GetSessionByNotebookId(notebookID string) (models.Session, error) {
	return s.SessionRepository.GetSessionByNotebookID(notebookID)
}

func (s *SessionService) UpdateSession(id string, session *models.Session) error {
	// Fetch the existing session
	existingSession, err := s.SessionRepository.GetSessionByID(id)
	if err != nil {
		return err
	}

	if session.IsActive {
		existingSession.IsActive = session.IsActive
		now := utils.GetStringTime()
		existingSession.StartedAt = &now
	} else {
		existingSession.IsActive = session.IsActive
		existingSession.EndedAt = session.EndedAt
	}	


	if session.StartedAt != nil {
		existingSession.StartedAt = session.StartedAt
	}
	if session.EndedAt != nil {
		existingSession.EndedAt = session.EndedAt
	}
	
	// Update the session in the repository
	return s.SessionRepository.UpdateSession(&existingSession)
}
