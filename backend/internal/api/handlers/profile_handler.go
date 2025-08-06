package handlers

import (
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
	}
}

func (ps *ProfileHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	var serverResponse utils.Response
	status := http.StatusOK
	IsMyProfile := false
	// get  LOGGED IN USER
	LoggedInUser, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		serverResponse.Message = "User not found in context"
		utils.RespondJSON(w, http.StatusUnauthorized, serverResponse)
		return
	}

	// Use id
	userIdstr := r.PathValue("userid")
	userId, err := strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid User Id"})
		return
	}

	if userId == LoggedInUser {
		IsMyProfile = true
	}
	var profileDetails models.ProfileDetails

	if IsMyProfile {
		profileDetails, err = ps.ProfileService.GetUserOwnProfile(LoggedInUser)
		if err != nil {
			serverResponse.Message = "Error fetching profile details"
			utils.RespondJSON(w, http.StatusInternalServerError, serverResponse)
			return
		}

	} else {
		profileDetails, err = ps.ProfileService.GetUserProfile(userId, LoggedInUser)
		if err != nil {
			serverResponse.Message = "Error fetching profile details"
			utils.RespondJSON(w, http.StatusInternalServerError, serverResponse)
			return
		}
	}
	if profileDetails.FollowbtnStatus == "follow" && !profileDetails.ProfilePublic {
		utils.RespondJSON(w, status, models.ProfileResponse{
			ProfileDetails: profileDetails,
		})
		return
	}

	posts, err := ps.ProfileService.GetUserPosts(userId)
	if err != nil {
		serverResponse.Message = "Error fetching posts"
		utils.RespondJSON(w, http.StatusInternalServerError, serverResponse)
		return
	}

	photos, err := ps.ProfileService.GetUserPhotos(userId)
	if err != nil {
		serverResponse.Message = "Error fetching photos"
		utils.RespondJSON(w, http.StatusInternalServerError, serverResponse)
		return
	}

	if profileDetails.Avatar != "" {
		photos = append([]models.Photo{{Image: profileDetails.Avatar}}, photos...)
	}

	utils.RespondJSON(w, status, models.ProfileResponse{
		ProfileDetails: profileDetails,
		UserPosts:      posts,
		Photos:         photos,
	})
}

func (f *ProfileHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userIdstr := r.PathValue("userid")
	userId, err := strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid User Id"})
		return
	}
	followers, err := f.ProfileService.GetFollowersList(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	utils.RespondJSON(w, http.StatusOK, followers)
}

func (f *ProfileHandler) GetFollowees(w http.ResponseWriter, r *http.Request) {
	userIdstr := r.PathValue("userid")
	userId, err := strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid User Id"})
		return
	}
	followers, err := f.ProfileService.GetFolloweesList(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	utils.RespondJSON(w, http.StatusOK, followers)
}
