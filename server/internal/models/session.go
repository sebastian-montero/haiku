package models

type Session struct {
	ID         int `json:"id"`
	NotebookID int `json:"notebook_id"`
	OwnerID    int `json:"owner_id"`

	IsActive  bool    `json:"is_active,omitempty"`
	StartedAt *string `json:"started_at,omitempty"`
	EndedAt   *string `json:"ended_at,omitempty"`
}
