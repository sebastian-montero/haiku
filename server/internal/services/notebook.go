package services

import (
	"haiku/internal/models"
	"haiku/internal/repositories"
	"strconv"
)

type NotebookService struct {
	Repository *repositories.NotebookRepository
	SessionRepository *repositories.SessionRepository
	ContentRepository *repositories.ContentRepository
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
	// Step 1: Retrieve the session associated with the notebook
	session, err := s.SessionRepository.GetSessionByNotebookID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// If no session exists, proceed with deleting the notebook
			return s.Repository.DeleteNotebookByID(id)
		}
		return err
	}

	// Step 2: Delete the content associated with the session
	err = s.ContentRepository.DeleteContentBySessionID(strconv.Itoa(session.ID))
	if err != nil {
		return err
	}

	// Step 3: Delete the session
	err = s.SessionRepository.DeleteSessionByID(strconv.Itoa(session.ID))
	if err != nil {
		return err
	}

	// Step 4: Delete the notebook
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
