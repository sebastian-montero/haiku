package services

import (
	"strconv"
	"wip/internal/models"
	"wip/internal/repositories"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
	s.Repository.CreateUser(user)
	return nil
}

func (s *UserService) GetUserByID(id string) models.User {
	return s.Repository.GetUserByID(id)
}

func (s *UserService) DeleteUserByID(id string) {
	s.Repository.DeleteUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User) {
	var _ models.User = s.Repository.GetUserByID(strconv.Itoa(user.ID))
	s.Repository.UpdateUser(user)
}
