package services

import (
	"haiku/internal/models"
	"haiku/internal/repositories"
	"strconv"
)

type ContentService struct {
	ContentRepository  *repositories.ContentRepository
	SessionRepository  *repositories.SessionRepository
	NotebookRepository *repositories.NotebookRepository
}

func (s *ContentService) CreateContent(content *models.Content) error {
	sessionID := strconv.Itoa(content.SessionID)
	session, err := s.SessionRepository.GetSessionByID(sessionID)
	if err != nil {
		return err
	}

	notebookID := strconv.Itoa(session.NotebookID)
	notebook, err := s.NotebookRepository.GetNotebookByID(notebookID)
	if err != nil {
		return err
	}
	notebook.LatestContent = content.Content

	err = s.NotebookRepository.UpdateNotebook(&notebook)
	if err != nil {
		return err
	}

	return s.ContentRepository.CreateContent(content)
}

func (s *ContentService) GetLatestContentBySessionId(sessionId string) (*models.Content, error) {
	return s.ContentRepository.GetLatestContentBySessionId(sessionId)
}
