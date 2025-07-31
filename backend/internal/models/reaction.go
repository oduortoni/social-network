package models

import "time"

// Reaction represents a like or dislike on a post or comment.	
type Reaction struct {
	UserID       int       `json:"user_id"`
	PostID       *int      `json:"post_id,omitempty"`
	CommentID    *int      `json:"comment_id,omitempty"`
	ReactionType string    `json:"reaction_type"`
	CreatedAt    time.Time `json:"created_at"`
}
