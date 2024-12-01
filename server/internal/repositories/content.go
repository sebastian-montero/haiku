package repositories

import (
	"database/sql"
	"haiku/internal/models"
	"haiku/internal/utils"
)

type ContentRepository struct {
	DB *sql.DB
}

func (repo *ContentRepository) CreateContent(content *models.Content) error {
	if content.CreatedAt == nil {
		now := utils.GetStringTime()
		content.CreatedAt = &now
	}
	query := "INSERT INTO content (session_id, content, created_at) VALUES ($1, $2, $3)"
	_, err := repo.DB.Exec(query, content.SessionID, content.Content, content.CreatedAt)
	return err
}

func (repo *ContentRepository) GetLatestContentBySessionId(sessionId string) (*models.Content, error) {
	var content models.Content
	query := "SELECT * FROM content WHERE session_id = $1 ORDER BY created_at DESC LIMIT 1"
	row := repo.DB.QueryRow(query, sessionId)
	err := row.Scan(&content.ID, &content.SessionID, &content.Content, &content.CreatedAt)
	return &content, err
}

func (repo *ContentRepository) DeleteContentBySessionID(sessionId string) error {
	query := "DELETE FROM content WHERE session_id = $1"
	_, err := repo.DB.Exec(query, sessionId)
	return err
}
