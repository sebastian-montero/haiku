package services

import (
	"strconv"
	"wout/internal/models"
	"wout/internal/repositories"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.Repository.CreateUser(user)
}

func (s *UserService) GetUserByID(id string) (models.User, error) {
	return s.Repository.GetUserByID(id)
}

func (s *UserService) DeleteUserByID(id string) error {
	return s.Repository.DeleteUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	_, err := s.Repository.GetUserByID(strconv.Itoa(user.ID))
	if err != nil {
		return err
	}
	return s.Repository.UpdateUser(user)
}
