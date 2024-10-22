package repositories

import (
	"database/sql"
	"wout/internal/models"
	"wout/internal/utils"
)

type SessionRepository struct {
	DB *sql.DB
}

// type Session struct {
// 	ID 	  int    `json:"id"`
// 	NotebookID int `json:"notebook_id"`
// 	OwnerID    int `json:"owner_id"`

// 	IsActive   *bool `json:"is_active,omitempty"`
// 	StartedAt  *string `json:"started_at,omitempty"`
// 	EndedAt    *string `json:"ended_at,omitempty"`
// }

func (repo *SessionRepository) CreateSession(session *models.Session) error {
	now := utils.GetStringTime()
	session.StartedAt = &now
	session.IsActive = true

	query := `INSERT INTO sessions (notebook_id, owner_id, is_active, started_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := repo.DB.QueryRow(query, session.NotebookID, session.OwnerID, session.IsActive, session.StartedAt).Scan(&session.ID)
	return err
}

// func (repo *NotebookRepository) GetNotebookByID(id string) (models.Notebook, error) {
// 	query := `SELECT * FROM notebooks WHERE id = $1`

// 	var notebook models.Notebook
// 	err := repo.DB.QueryRow(query, id).Scan(&notebook.ID, &notebook.Title, &notebook.OwnerID, &notebook.LatestContent, &notebook.CreatedAt, &notebook.LastUpdatedAt)
// 	return notebook, err
// }

// func (repo *NotebookRepository) DeleteNotebookByID(id string) error {
// 	query := `DELETE FROM notebooks WHERE id = $1`
// 	_, err := repo.DB.Exec(query, id)
// 	return err
// }

// func (repo *NotebookRepository) UpdateNotebook(notebook *models.Notebook) error {
// 	now := utils.GetStringTime()
// 	query := `UPDATE notebooks SET title = $1, latest_content = $2, last_updated_at = $3 WHERE id = $4`
// 	_, err := repo.DB.Exec(query, notebook.Title, notebook.LatestContent, now, notebook.ID)
// 	return err
// }
