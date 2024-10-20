package repositories

import (
	"database/sql"
	"wip/internal/models"
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

func (repo *UserRepository) GetUserByID(id string) models.User {
	query := `SELECT * FROM users WHERE id = $1`

	var user models.User
	err := repo.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		panic(err)
	}
	return user
}

func (repo *UserRepository) DeleteUserByID(id string) {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	if err != nil {
		panic(err)
	}
}

func (repo *UserRepository) UpdateUser(user *models.User) {
	query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		panic(err)
	}
}
