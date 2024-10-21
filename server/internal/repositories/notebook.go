package repositories

import (
	"database/sql"
	"wout/internal/models"
	"wout/internal/utils"
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
	return err
}

func (repo *NotebookRepository) GetNotebookByID(id string) (models.Notebook, error) {
	query := `SELECT * FROM notebooks WHERE id = $1`

	var notebook models.Notebook
	err := repo.DB.QueryRow(query, id).Scan(&notebook.ID, &notebook.Title, &notebook.OwnerID, &notebook.LatestContent, &notebook.CreatedAt, &notebook.LastUpdatedAt)
	return notebook, err
}

func (repo *NotebookRepository) DeleteNotebookByID(id string) error {
	query := `DELETE FROM notebooks WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	return err
}

func (repo *NotebookRepository) UpdateNotebook(notebook *models.Notebook) error {
	now := utils.GetStringTime()
	query := `UPDATE notebooks SET title = $1, latest_content = $2, last_updated_at = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, notebook.Title, notebook.LatestContent, now, notebook.ID)
	return err
}
