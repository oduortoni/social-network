package models

import "time"

type Group struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	Privacy     string    `json:"privacy"` // e.g., "public", "private"
	CreatedAt   time.Time `json:"created_at"`
}
