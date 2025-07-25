package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/utils"
)

type UnfollowHandler struct {
	UnfollowService service.UnfollowServiceInterface
}

func NewUnfollowHandler(unf service.UnfollowServiceInterface) *UnfollowHandler {
	return &UnfollowHandler{UnfollowService: unf}
}

func (unf *UnfollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	var serverResponse models.Response
	status := http.StatusOK

	// Get current user id from session
	followerID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		serverResponse.Message = "User not found in context"
		models.RespondJSON(w, http.StatusUnauthorized, serverResponse)
		return
	}

	var followee models.Followee
	body, err := io.ReadAll(r.Body)
	if err != nil {
		status = http.StatusBadRequest
		serverResponse.Message = "Failed to read request body"
		models.RespondJSON(w, status, serverResponse)
		return
	}

	if err = json.Unmarshal(body, &followee); err != nil {
		status = http.StatusBadRequest
		serverResponse.Message = "Invalid JSON format"
		models.RespondJSON(w, status, serverResponse)
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
		models.RespondJSON(w, status, serverResponse)
		return
	}

	// Delete the follow connection
	err = unf.UnfollowService.DeleteFollowConnection(followConnectionID)
	if err != nil {
		status = http.StatusInternalServerError
		serverResponse.Message = "Failed to unfollow user"
		models.RespondJSON(w, status, serverResponse)
		return
	}

	serverResponse.Message = "Successfully unfollowed user"
	models.RespondJSON(w, status, serverResponse)
}
