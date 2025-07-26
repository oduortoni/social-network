package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type UnfollowHandler struct {
	UnfollowService service.UnfollowServiceInterface
}

func NewUnfollowHandler(unf service.UnfollowServiceInterface) *UnfollowHandler {
	return &UnfollowHandler{UnfollowService: unf}
}

func (unf *UnfollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	var serverResponse utils.Response
	status := http.StatusOK

	// Get current user id from session
	followerID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		serverResponse.Message = "User not found in context"
		utils.RespondJSON(w, http.StatusUnauthorized, serverResponse)
		return
	}

	var followee models.Followee
	body, err := io.ReadAll(r.Body)
	if err != nil {
		status = http.StatusBadRequest
		serverResponse.Message = "Failed to read request body"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	if err = json.Unmarshal(body, &followee); err != nil {
		status = http.StatusBadRequest
		serverResponse.Message = "Invalid JSON format"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Get the follow connection ID
	followConnectionID, err := unf.UnfollowService.GetFollowConnectionID(followerID, int64(followee.FolloweeId))
	if err != nil {
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			serverResponse.Message = "Follow relationship not found"
		} else {
			status = http.StatusInternalServerError
			serverResponse.Message = "Failed to find follow connection"
		}
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Delete the follow connection
	err = unf.UnfollowService.DeleteFollowConnection(followConnectionID)
	if err != nil {
		status = http.StatusInternalServerError
		serverResponse.Message = "Failed to unfollow user"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	serverResponse.Message = "Successfully unfollowed user"
	utils.RespondJSON(w, status, serverResponse)
}
