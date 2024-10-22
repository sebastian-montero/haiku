package models

type Content struct {
	ID        int     `json:"id"`
	SessionID int     `json:"session_id"`
	Content   string  `json:"content"`
	CreatedAt *string `json:"created_at,omitempty"`
}
