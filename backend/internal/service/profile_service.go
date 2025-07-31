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
	// Fetch the user's profile details
	// This will include first name, last name, email, nickname, about me, date of birth, profile visibility, and avatar
	// The userid is used to fetch the details of the logged-in user
	user, err := ps.ProfileStore.MyProfileDetails(userid)
	if err != nil {
		return userDetails, err
	}
	// Get the number of followers and followees
	userDetails.NumberOfFollowers, userDetails.NumberOfFollowees, err = ps.ProfileStore.GetFollowersStat(userid)
	if err != nil {
		return userDetails, err
	}

	userDetails.NumberOfPosts, err = ps.ProfileStore.GetNumberOfPosts(userid)
	if err != nil {
		return userDetails, err
	}

	userDetails = user
	userDetails.ID = userid
	userDetails.FollowbtnStatus = "hide"
	userDetails.MessageBtnStatus = "hide"
	return userDetails, nil
}

func (ps *ProfileService) GetUserProfile(userid, LoggedInUser int64) (models.ProfileDetails, error) {
	var userDetails models.ProfileDetails
	// Fetch the user's profile details
	// This will include first name, last name, email, nickname, about me, date of birth, profile visibility, and avatar
	user, err := ps.ProfileStore.MyProfileDetails(userid)
	if err != nil {
		return userDetails, err
	}
	// Get the number of followers and followees
	userDetails.NumberOfFollowers, userDetails.NumberOfFollowees, err = ps.ProfileStore.GetFollowersStat(userid)
	if err != nil {
		return userDetails, err
	}

	userDetails.NumberOfPosts, err = ps.ProfileStore.GetNumberOfPosts(userid)
	if err != nil {
		return userDetails, err
	}

	userDetails = user
	userDetails.ID = userid
	userDetails.FollowbtnStatus, err = ps.ProfileStore.GetFollowStatus(userid, LoggedInUser)
	if err != nil {
		return userDetails, err
	}
	if userDetails.FollowbtnStatus == "follow" {
		userDetails.MessageBtnStatus = "visible"
	} else {
		userDetails.MessageBtnStatus = "hide"
	}

	return userDetails, nil
}
