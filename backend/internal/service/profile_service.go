package service

import (
	"fmt"
	"html"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type ProfileService struct {
	ProfileStore *store.ProfileStore
}

func NewProfileService(ps *store.ProfileStore) *ProfileService {
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
	user.NumberOfFollowers, user.NumberOfFollowees, err = ps.ProfileStore.GetFollowersStat(userid)
	if err != nil {
		return userDetails, err
	}

	user.NumberOfPosts, err = ps.ProfileStore.GetNumberOfPosts(userid)
	if err != nil {
		return userDetails, err
	}

	userDetails = user
	userDetails.ID = userid
	userDetails.FollowbtnStatus = "hide"
	userDetails.MessageBtnStatus = "hide"
	userDetails.About = html.UnescapeString(userDetails.About)
	userDetails.FirstName = html.UnescapeString(userDetails.FirstName)
	userDetails.LastName = html.UnescapeString(userDetails.LastName)
	userDetails.Nickname = html.UnescapeString(userDetails.Nickname)

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
	user.NumberOfFollowers, user.NumberOfFollowees, err = ps.ProfileStore.GetFollowersStat(userid)
	if err != nil {
		return userDetails, err
	}

	user.NumberOfPosts, err = ps.ProfileStore.GetNumberOfPosts(userid)
	if err != nil {
		return userDetails, err
	}

	userDetails = user
	userDetails.ID = userid
	userDetails.FollowbtnStatus, err = ps.ProfileStore.GetFollowStatus(userid, LoggedInUser)
	if err != nil {
		return userDetails, err
	}
	// The message button should be visible if the logged-in user is following the profile user.
	if userDetails.FollowbtnStatus == "following" {
		userDetails.MessageBtnStatus = "visible"
	} else {
		userDetails.MessageBtnStatus = "hide"
	}
	userDetails.About = html.UnescapeString(userDetails.About)
	fmt.Println("USER ABOUT", userDetails.About)
	userDetails.FirstName = html.UnescapeString(userDetails.FirstName)
	userDetails.LastName = html.UnescapeString(userDetails.LastName)
	userDetails.Nickname = html.UnescapeString(userDetails.Nickname)

	return userDetails, nil
}

func (ps *ProfileService) GetUserPosts(userid int64) ([]models.Post, error) {
	return ps.ProfileStore.GetPostsOfUser(userid)
}

func (ps *ProfileService) GetFollowersList(userid int64) (models.FollowListResponse, error) {
	return ps.ProfileStore.GetUserFollowers(userid)
}

func (ps *ProfileService) GetFolloweesList(userid int64) (models.FollowListResponse, error) {
	return ps.ProfileStore.GetUserFollowees(userid)
}

func (ps *ProfileService) GetUserPhotos(userId int64) ([]models.Photo, error) {
	postphotos, err := ps.ProfileStore.GetUserPostPhotos(userId)
	if err != nil {
		return nil, err
	}
	commentphots, err := ps.ProfileStore.GetUserCommentPhotos(userId)
	if err != nil {
		return nil, err
	}
	photos := append(postphotos, commentphots...)

	var actual []models.Photo
	for i := range photos {
		if photos[i].Image != "" {
			actual = append(actual, photos[i])
		}
	}

	return actual, nil
}
