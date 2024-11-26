package services

import (
	"haiku/internal/models"
	"haiku/internal/repositories"
	"strconv"
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

func (s *NotebookService) GetNotebooksByOwnerId(ownerId string) ([]models.Notebook, error) {
	return s.Repository.GetNotebooksByOwnerId(ownerId)
}

func (s *NotebookService) DeleteNotebookByID(id string) error {
	_, err := s.Repository.GetNotebookByID(id)
	if err != nil {
		return err
	}
	return s.Repository.DeleteNotebookByID(id)
}

func (s *NotebookService) UpdateNotebook(notebook *models.Notebook) error {
	_, err := s.Repository.GetNotebookByID(strconv.Itoa(notebook.ID))
	if err != nil {
		return err
	}

	err = s.Repository.UpdateNotebook(notebook)
	if err != nil {
		return err
	}
	updatedNotebook, err := s.Repository.GetNotebookByID(strconv.Itoa(notebook.ID))
	if err != nil {
		return err
	}

	*notebook = updatedNotebook
	return nil
}
