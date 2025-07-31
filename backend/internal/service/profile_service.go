package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type ProfileService struct {
	ProfileStore store.ProfileStore
}

func NewProfileService(ps store.ProfileStore) *ProfileService {
	return &ProfileService{ProfileStore: ps}
}

func (ps *ProfileService) GetUserOwnProfile(userid int64) (models.ProfileDetails, error) {
	var userDetails models.ProfileDetails
	userDetails.FollowbtnStatus="hide"
	userDetails.MessageBtnStatus="hide"
	
	return userDetails, nil
}
