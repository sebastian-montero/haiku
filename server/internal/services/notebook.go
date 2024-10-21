package services

import (
	"wout/internal/models"
	"wout/internal/repositories"
)

type NotebookService struct {
	Repository *repositories.NotebookRepository
}

func (s *NotebookService) CreateNotebook(notebook *models.Notebook) error {
	return s.Repository.CreateNotebook(notebook)
}

func (s *NotebookService) GetNotebookByID(id string) (models.Notebook, error) {
	return s.Repository.GetNotebookByID(id)
}

func (s *NotebookService) DeleteNotebookByID(id string) error {
	_, err := s.Repository.GetNotebookByID(id)
	if err != nil {
		return err
	}
	return s.Repository.DeleteNotebookByID(id)
}

// func (s *NotebookService) UpdateUser(user *models.User) {
// 	s.Repository.UpdateUser(user)
// }
