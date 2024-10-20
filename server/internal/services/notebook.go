package services

import (
	"wip/internal/models"
	"wip/internal/repositories"
)

type NotebookService struct {
	Repository *repositories.NotebookRepository
}

func (s *NotebookService) CreateNotebook(notebook *models.Notebook) error {
	s.Repository.CreateNotebook(notebook)
	return nil
}

// func (s *NotebookService) GetUserByID(id string) models.User {
// 	return s.Repository.GetUserByID(id)
// }

// func (s *NotebookService) DeleteUserByID(id string) {
// 	s.Repository.DeleteUserByID(id)
// }

// func (s *NotebookService) UpdateUser(user *models.User) {
// 	s.Repository.UpdateUser(user)
// }
