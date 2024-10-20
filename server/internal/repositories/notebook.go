package repositories

import (
	"database/sql"
	"wip/internal/models"
	"wip/internal/utils"
)

type NotebookRepository struct {
	DB *sql.DB
}

func (repo *NotebookRepository) CreateNotebook(notebook *models.Notebook) error {
	if notebook.CreatedAt == nil {
		now := utils.GetStringTime()
		notebook.CreatedAt = &now
	}
	notebook.LastUpdatedAt = notebook.CreatedAt

	query := `INSERT INTO notebooks (title, owner_id, latest_content, created_at, last_updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := repo.DB.QueryRow(query, notebook.Title, notebook.OwnerID, notebook.LatestContent, notebook.CreatedAt, notebook.LastUpdatedAt).Scan(&notebook.ID)
	if err != nil {
		panic(err)
	}
	return nil
}

// func (repo *NotebookRepository) GetUserByID(id string) models.User {
// 	query := `SELECT * FROM users WHERE id = $1`

// 	var user models.User
// 	err := repo.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return user
// }

// func (repo *NotebookRepository) DeleteUserByID(id string) {
// 	query := `DELETE FROM users WHERE id = $1`
// 	_, err := repo.DB.Exec(query, id)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func (repo *NotebookRepository) UpdateUser(user *models.User) {
// 	query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4`
// 	_, err := repo.DB.Exec(query, user.Username, user.Email, user.Password, user.ID)
// 	if err != nil {
// 		panic(err)
// 	}
// }
