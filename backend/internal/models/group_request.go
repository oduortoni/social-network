package models

import "time"

type GroupRequest struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"` // e.g., "pending", "approved", "rejected"
	CreatedAt time.Time `json:"created_at"`
}
