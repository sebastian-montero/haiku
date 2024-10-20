package models

type Notebook struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	OwnerID   int     `json:"owner_id"`
	LatestContent string `json:"latest_content,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	LastUpdatedAt *string `json:"updated_at,omitempty"`
}
