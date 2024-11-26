package services

import (
	"haiku/internal/models"
	"haiku/internal/repositories"
	"strconv"
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
	_, err := s.Repository.GetUserByID(id)
	if err != nil {
		return err
	}
	return s.Repository.DeleteUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	_, err := s.Repository.GetUserByID(strconv.Itoa(user.ID))
	if err != nil {
		return err
	}

	err = s.Repository.UpdateUser(user)
	if err != nil {
		return err
	}

	updatedUser, err := s.Repository.GetUserByID(strconv.Itoa(user.ID))
	if err != nil {
		return err
	}
	*user = updatedUser
	return nil
}
