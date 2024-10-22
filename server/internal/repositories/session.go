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

func (repo *SessionRepository) GetActiveSessions() ([]models.Session, error) {
	query := `SELECT * FROM sessions WHERE is_active = true`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(&session.ID, &session.NotebookID, &session.OwnerID, &session.IsActive, &session.StartedAt, &session.EndedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (repo *SessionRepository) GetSessionByID(id string) (models.Session, error) {
	query := `SELECT * FROM sessions WHERE id = $1`

	var session models.Session
	err := repo.DB.QueryRow(query, id).Scan(&session.ID, &session.NotebookID, &session.OwnerID, &session.IsActive, &session.StartedAt, &session.EndedAt)
	return session, err
}

func (repo *SessionRepository) DeleteSessionByID(id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	return err
}
