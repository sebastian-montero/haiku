package repositories

import (
	"database/sql"
	"haiku/internal/models"
	"haiku/internal/utils"
)

type SessionRepository struct {
	DB *sql.DB
}

func (repo *SessionRepository) CreateSession(session *models.Session) error {
	now := utils.GetStringTime()
	session.StartedAt = &now

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

func (repo *SessionRepository) UpdateSession(session *models.Session) error {
	query := `UPDATE sessions SET is_active = $1, started_at = $2, ended_at = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, session.IsActive, session.StartedAt, session.EndedAt, session.ID)
	return err
}

func (repo *SessionRepository) SessionExistsByNotebookID(notebookID string) (models.Session, bool) {
	query := `SELECT * FROM sessions WHERE notebook_id = $1`

	var session models.Session
	err := repo.DB.QueryRow(query, notebookID).Scan(&session.ID, &session.NotebookID, &session.OwnerID, &session.IsActive, &session.StartedAt, &session.EndedAt)
	if err != nil {
		return session, false
	}
	return session, true
}

func (repo *SessionRepository) GetSessionByNotebookID(notebookID string) (models.Session, error) {
	query := `SELECT * FROM sessions WHERE notebook_id = $1`

	var session models.Session
	err := repo.DB.QueryRow(query, notebookID).Scan(&session.ID, &session.NotebookID, &session.OwnerID, &session.IsActive, &session.StartedAt, &session.EndedAt)
	return session, err
}
