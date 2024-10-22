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

// func (s *NotebookService) GetNotebookByID(id string) (models.Notebook, error) {
// 	return s.Repository.GetNotebookByID(id)
// }

// func (s *NotebookService) DeleteNotebookByID(id string) error {
// 	_, err := s.Repository.GetNotebookByID(id)
// 	if err != nil {
// 		return err
// 	}
// 	return s.Repository.DeleteNotebookByID(id)
// }

// func (s *NotebookService) UpdateNotebook(notebook *models.Notebook) error {
// 	_, err := s.Repository.GetNotebookByID(strconv.Itoa(notebook.ID))
// 	if err != nil {
// 		return err
// 	}

// 	err = s.Repository.UpdateNotebook(notebook)
// 	if err != nil {
// 		return err
// 	}
// 	updatedNotebook, err := s.Repository.GetNotebookByID(strconv.Itoa(notebook.ID))
// 	if err != nil {
// 		return err
// 	}

// 	*notebook = updatedNotebook
// 	return nil
// }
