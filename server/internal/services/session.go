package services

import (
	"wout/internal/models"
	"wout/internal/repositories"
)

type SessionService struct {
	Repository *repositories.SessionRepository
}

func (s *SessionService) CreateSession(session *models.Session) error {
	return s.Repository.CreateSession(session)
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
