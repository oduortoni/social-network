package service

import "github.com/tajjjjr/social-network/backend/internal/store"

type UnfollowService struct {
	UnfollowStore *store.UnfollowStore
}

func NewUnfollowService(unfollowStore *store.UnfollowStore) *UnfollowService {
	return &UnfollowService{UnfollowStore: unfollowStore}
}

func (unf *UnfollowService) GetFollowConnectionID(followerID, followeeID int64) (int64, error) {
	return unf.UnfollowStore.GetFollowConnectionID(followerID, followeeID)
}

func (unf *UnfollowService) DeleteFollowConnection(followConnectionID int64) error {
	return unf.UnfollowStore.DeleteFollowConnection(followConnectionID)
}
