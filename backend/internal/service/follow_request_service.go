package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

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

func (fr *FollowRequestService) RetrieveUserName(userID int64) (string, string, error) {
	return fr.FollowRequestStore.UserInfo(userID)
}

func (fr *FollowRequestService) GetRequestInfo(requestID int64) (int64, int64, error) {
	return fr.FollowRequestStore.RetrieveRequestInfo(requestID)
}

func (fr *FollowRequestService) AddtoNotification(follower_id int64, message string) error {
	return fr.FollowRequestStore.AddtoNotification(follower_id, message)
}

func (fr *FollowRequestService) CancelFollowRequest(followConnectionID int64) error {
	return fr.FollowRequestStore.FollowRequestCancel(followConnectionID)
}

func (fr *FollowRequestService) GetPendingFollowRequest(userid int64) (models.FollowRequestUserResponse, error) {
	return fr.FollowRequestStore.GetPendingFollowRequest(userid)
}
