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

func (repo *NotebookRepository) GetNotebooksByOwnerId(owner_id string) ([]models.Notebook, error) {
	query := `SELECT * FROM notebooks WHERE owner_id = $1`

	rows, err := repo.DB.Query(query, owner_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notebooks []models.Notebook
	for rows.Next() {
		var notebook models.Notebook
		err = rows.Scan(&notebook.ID, &notebook.Title, &notebook.OwnerID, &notebook.LatestContent, &notebook.CreatedAt, &notebook.LastUpdatedAt)
		if err != nil {
			return nil, err
		}
		notebooks = append(notebooks, notebook)
	}
	return notebooks, nil
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

func (repo *NotebookRepository) NotebookHasOwnerId(owner_id int, notebook_id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM notebooks WHERE owner_id = $1 AND id = $2)`
	var exists bool
	err := repo.DB.QueryRow(query, owner_id, notebook_id).Scan(&exists)
	return exists, err
}
