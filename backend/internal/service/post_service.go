package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type PostService struct {
	PostStore *store.PostStore
}

func NewPostService(ps *store.PostStore) *PostService {
	return &PostService{PostStore: ps}
}

func (s *PostService) CreatePost(post *models.Post) (int64, error) {
	return s.PostStore.CreatePost(post)
}
