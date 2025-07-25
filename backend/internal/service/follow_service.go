package service

import (
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type FollowService struct {
	FollowStore *store.FollowStore
}

func NewFollowService(ff *store.FollowStore) *FollowService {
	return &FollowService{FollowStore: ff}
}

func (Follow *FollowService) IsAccountPublic(followee int64) (bool, error) {
	return Follow.FollowStore.IsUserAccountPublic(followee)
}

func (Follow *FollowService) CreateFollowForPublicAccount(followerid, followeeid int64) error {
	return Follow.FollowStore.CreatePublicFollowConnection(followerid, followeeid)
}

func (Follow *FollowService) CreateFollowForPrivateAccount(followrid, followeeid int64) (int64, error) {
	return Follow.FollowStore.CreatePrivateFollowConnection(followrid, followeeid)
}
