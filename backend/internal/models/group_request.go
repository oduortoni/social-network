package models

import "time"

type GroupRequest struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"` // e.g., "pending", "approved", "rejected"
	CreatedAt time.Time `json:"created_at"`
}
