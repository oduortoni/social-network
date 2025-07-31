package store

import (
	"database/sql"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type ReactionStore struct {
	*sql.DB
}

func NewReactionStore(db *sql.DB) *ReactionStore {
	return &ReactionStore{db}
}

// AddPostReaction adds a reaction to a post.
func (s *ReactionStore) AddPostReaction(reaction *models.Reaction) error {
	_, err := s.Exec(`
		INSERT INTO post_reactions (user_id, post_id, reaction_type)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) DO UPDATE SET reaction_type = excluded.reaction_type;
	`, reaction.UserID, reaction.PostID, reaction.ReactionType)
	return err
}

// RemovePostReaction removes a reaction from a post.
func (s *ReactionStore) RemovePostReaction(userID, postID int) error {
	_, err := s.Exec(`
		DELETE FROM post_reactions
		WHERE user_id = ? AND post_id = ?;
	`, userID, postID)
	return err
}

// AddCommentReaction adds a reaction to a comment.
func (s *ReactionStore) AddCommentReaction(reaction *models.Reaction) error {
	_, err := s.Exec(`
		INSERT INTO comment_reactions (user_id, comment_id, reaction_type)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, comment_id) DO UPDATE SET reaction_type = excluded.reaction_type;
	`, reaction.UserID, reaction.CommentID, reaction.ReactionType)
	return err
}

// RemoveCommentReaction removes a reaction from a comment.
func (s *ReactionStore) RemoveCommentReaction(userID, commentID int) error {
	_, err := s.Exec(`
		DELETE FROM comment_reactions
		WHERE user_id = ? AND comment_id = ?;
	`, userID, commentID)
	return err
}