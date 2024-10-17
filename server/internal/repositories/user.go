package repositories

import (
	"database/sql"
	"wip/internal/models"
	"wip/internal/utils/logger"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := repo.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		panic(err)
	}
	return nil
}

func UserRepositoryManager(db *sql.DB) *UserRepository {
	logger.Info("Creating user repository...")
	return &UserRepository{DB: db}
}
