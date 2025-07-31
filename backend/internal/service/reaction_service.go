package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type ReactionService struct {
	store *store.ReactionStore
}

func NewReactionService(store *store.ReactionStore) *ReactionService {
	return &ReactionService{store}
}

func (s *ReactionService) ReactToPost(reaction *models.Reaction) error {
	return s.store.AddPostReaction(reaction)
}

func (s *ReactionService) UnreactToPost(userID, postID int) error {
	return s.store.RemovePostReaction(userID, postID)
}

func (s *ReactionService) ReactToComment(reaction *models.Reaction) error {
	return s.store.AddCommentReaction(reaction)
}

func (s *ReactionService) UnreactToComment(userID, commentID int) error {
	return s.store.RemoveCommentReaction(userID, commentID)
}
