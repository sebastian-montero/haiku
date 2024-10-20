package services

import (
	"wip/internal/models"
	"wip/internal/repositories"
	"wip/internal/utils/logger"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
	s.UserRepository.CreateUser(user)
	return nil
}

func (s *UserService) GetUserByID(id string) models.User {
	return s.UserRepository.GetUserByID(id)
}

func (s *UserService) DeleteUserByID(id string) {
	s.UserRepository.DeleteUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User) {
	s.UserRepository.UpdateUser(user)
}

func UserServiceManager(userRepository *repositories.UserRepository) *UserService {
	logger.Info("Creating user service...")
	return &UserService{UserRepository: userRepository}
}
