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

func UserServiceManager(userRepository *repositories.UserRepository) *UserService {
	logger.Info("Creating user service...")
	return &UserService{UserRepository: userRepository}
}
