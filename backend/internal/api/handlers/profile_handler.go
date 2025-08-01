package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type ProfileHandler struct {
	ProfileService service.ProfileServiceInterface
}

func NewProfileHandler(profileService service.ProfileServiceInterface) *ProfileHandler {
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
	fmt.Println(status)
	fmt.Println(IsMyProfile)
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

	posts, err := ps.ProfileService.GetUserPosts(userId)
	if err != nil {
		serverResponse.Message = "Error fetching posts"
		utils.RespondJSON(w, http.StatusInternalServerError, serverResponse)
		return
	}

	utils.RespondJSON(w, status, models.ProfileResponse{
		ProfileDetails: profileDetails,
		UserPosts:      posts,
	})
}
