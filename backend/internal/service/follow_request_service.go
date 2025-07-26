package service

import "github.com/tajjjjr/social-network/backend/internal/store"

type FollowRequestService struct {
	FollowRequestStore *store.FollowRequestStore
}

func NewFollowRequestService(fr *store.FollowRequestStore) *FollowRequestService {
	return &FollowRequestService{FollowRequestStore: fr}
}

func (fr *FollowRequestService) AcceptedFollowConnection(followConnectionID int64) error {
	return fr.FollowRequestStore.AcceptFollowConnection(followConnectionID)
}

func (fr *FollowRequestService) RejectedFollowConnection(followConnectionID int64) error {
	return fr.FollowRequestStore.RejectFollowConnection(followConnectionID)
}
