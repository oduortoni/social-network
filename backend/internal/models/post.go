package models

import "time"

type Post struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	Content      string     `json:"content"`
	Image        string     `json:"image,omitempty"`
	Privacy      string     `json:"privacy"` // "public", "private", "followers"
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	IsEdited     bool       `json:"is_edited"`
	Author       User       `json:"author"`
	LikesCount   int        `json:"likes_count"`
	DislikesCount int       `json:"dislikes_count"`
	UserReaction *string    `json:"user_reaction,omitempty"`
}
