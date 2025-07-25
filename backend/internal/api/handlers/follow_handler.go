package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/utils"
)

type FollowHandler struct {
	FollowService service.FollowServiceInterface
}

func NewFollowHandler(funf service.FollowServiceInterface) *FollowHandler {
	return &FollowHandler{FollowService: funf}
}

func (follow *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var serverResponse models.Response
	status := http.StatusOK

	// Get current user id from session
	followerId, ok := r.Context().Value(utils.User_id).(int64)
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
		status = http.StatusBadRequest // ✅ Corrected from 500
		serverResponse.Message = "Invalid JSON format"
		models.RespondJSON(w, status, serverResponse)
		return
	}

	isFolloweeAccountPublic, err := follow.FollowService.IsAccountPublic(int64(followee.FolloweeId))
	if err != nil {
		status = http.StatusInternalServerError
		serverResponse.Message = "Failed to check account privacy status"
		models.RespondJSON(w, status, serverResponse)
		return
	}

	if isFolloweeAccountPublic {
		err = follow.FollowService.CreateFollowForPublicAccount(followerId, int64(followee.FolloweeId))
		if err != nil {
			status = http.StatusInternalServerError
			serverResponse.Message = "Failed to create follow connection"
			models.RespondJSON(w, status, serverResponse)
			return
		}
		serverResponse.Message = "Account successfully followed"
		models.RespondJSON(w, status, serverResponse)
		return
	}

	// Optional: Handle private account response
	status = http.StatusForbidden
	serverResponse.Message = "Cannot follow private account without approval"
	models.RespondJSON(w, status, serverResponse)
}
