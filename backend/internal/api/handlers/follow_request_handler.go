package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type FollowRequestHandler struct {
	FollowRequestService service.FollowRequestServiceInterface
}

func NewFollowRequestHandler(fr service.FollowRequestServiceInterface) *FollowRequestHandler {
	return &FollowRequestHandler{FollowRequestService: fr}
}

func (fr *FollowRequestHandler) FollowRequestRespond(w http.ResponseWriter, r *http.Request) {
	var serverResponse utils.Response
	status := http.StatusOK

	// Get current user id from session (to verify they can respond to this request)
	_, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		serverResponse.Message = "User not found in context"
		utils.RespondJSON(w, http.StatusUnauthorized, serverResponse)
		return
	}

	// Parse request ID from URL path
	requestIDStr := r.PathValue("requestId")
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		serverResponse.Message = "Invalid request ID"
		utils.RespondJSON(w, http.StatusBadRequest, serverResponse)
		return
	}

	// Parse request body
	var requestStatus models.FollowRequestResponseStatus
	body, err := io.ReadAll(r.Body)
	if err != nil {
		serverResponse.Message = "Failed to read request body"
		utils.RespondJSON(w, http.StatusBadRequest, serverResponse)
		return
	}

	if err = json.Unmarshal(body, &requestStatus); err != nil {
		serverResponse.Message = "Invalid JSON format"
		utils.RespondJSON(w, http.StatusBadRequest, serverResponse)
		return
	}

	// Validate status value
	if requestStatus.Status != "accepted" && requestStatus.Status != "rejected" {
		serverResponse.Message = "Invalid status. Must be 'accepted' or 'rejected'"
		utils.RespondJSON(w, http.StatusBadRequest, serverResponse)
		return
	}

	// Handle rejection
	if requestStatus.Status == "rejected" {
		err = fr.FollowRequestService.RejectedFollowConnection(requestID)
		if err != nil {
			if err == sql.ErrNoRows {
				status = http.StatusNotFound
				serverResponse.Message = "Follow request not found"
			} else {
				status = http.StatusInternalServerError
				serverResponse.Message = "Failed to reject follow request"
			}
			utils.RespondJSON(w, status, serverResponse)
			return
		}
		serverResponse.Message = "Successfully rejected follow request"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Handle acceptance
	err = fr.FollowRequestService.AcceptedFollowConnection(requestID)
	if err != nil {
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			serverResponse.Message = "Follow request not found"
		} else {
			status = http.StatusInternalServerError
			serverResponse.Message = "Failed to accept follow request"
		}
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	serverResponse.Message = "Successfully accepted follow request"
	utils.RespondJSON(w, status, serverResponse)
}
