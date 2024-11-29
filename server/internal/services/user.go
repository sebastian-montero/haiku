package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"haiku/internal/models"
	"haiku/internal/repositories"
	"haiku/internal/utils"
	"strconv"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return errors.New("failed to generate salt")
	}

	hashedPassword, err := utils.HashPassword(user.Password, salt)
	if err != nil {
		return err
	}

	user.Salt = base64.StdEncoding.EncodeToString(salt)
	user.Password = hashedPassword

	return s.Repository.CreateUser(user)
}

func (s *UserService) GetUserByID(id string) (models.User, error) {
	return s.Repository.GetUserByID(id)
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.Repository.GetUserByUsername(username)
}

func (s *UserService) DeleteUserByID(id string) error {
	_, err := s.Repository.GetUserByID(id)
	if err != nil {
		return err
	}
	return s.Repository.DeleteUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	existingUser, err := s.Repository.GetUserByID(strconv.Itoa(user.ID))
	if err != nil {
		return err
	}

	user.Password = existingUser.Password
	user.Salt = existingUser.Salt

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
