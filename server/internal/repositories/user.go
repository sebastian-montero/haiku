package repositories

import (
	"database/sql"
	"haiku/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password, salt) VALUES ($1, $2, $3, $4) RETURNING id`
	err := repo.DB.QueryRow(query, user.Username, user.Email, user.Password, user.Salt).Scan(&user.ID)
	return err
}

func (repo *UserRepository) GetUserByID(id string) (models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user models.User
	err := repo.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Salt)
	return user, err
}

func (repo *UserRepository) DeleteUserByID(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	return err
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password, salt FROM users WHERE username = $1"
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Salt)
	if err != nil {
		return user, err
	}

	return user, nil
}
