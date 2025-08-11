package models

import "time"

type GroupChatMessage struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
