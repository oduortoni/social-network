package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	ws "github.com/tajjjjr/social-network/backend/internal/websocket"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type FollowRequestHandler struct {
	FollowRequestService service.FollowRequestServiceInterface
	Notifier             *ws.NotificationSender
}

func NewFollowRequestHandler(fr service.FollowRequestServiceInterface, notifier *ws.NotificationSender) *FollowRequestHandler {
	return &FollowRequestHandler{
		FollowRequestService: fr,
		Notifier:             notifier,
	}
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
		// Get follow request info before rejecting (only if notifications are enabled)
		var followerID, followeeID int64
		if fr.Notifier != nil {
			followerID, followeeID, err = fr.FollowRequestService.GetRequestInfo(requestID)
			if err != nil {
				status = http.StatusNotFound
				serverResponse.Message = "Follow request not found"
				utils.RespondJSON(w, status, serverResponse)
				return
			}
		}

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

		// Send notification to the requester about rejection
		if fr.Notifier != nil && followerID != 0 {
			followeeName, _, err := fr.FollowRequestService.RetrieveUserName(followeeID)
			if err == nil {
				// Store notification in database
				err = fr.FollowRequestService.AddtoNotification(followerID, followeeName+" rejected your follow request")
				if err != nil {
					status = http.StatusInternalServerError
					serverResponse.Message = "Failed to add notification"
					utils.RespondJSON(w, status, serverResponse)
					return
				}

				// Send real-time notification if user is online
				if fr.Notifier.IsOnline(followerID) {
					fr.Notifier.SendNotification(followerID, map[string]interface{}{
						"type":      "notification",
						"subtype":   "follow_request_rejected",
						"user_id":   followeeID,
						"user_name": followeeName,
						"message":   followeeName + " rejected your follow request",
						"timestamp": time.Now().Unix(),
					})
				}
			}
		}

		serverResponse.Message = "Successfully rejected follow request"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Handle acceptance
	// Get follow request info before accepting (only if notifications are enabled)
	var followerID, followeeID int64
	if fr.Notifier != nil {
		followerID, followeeID, err = fr.FollowRequestService.GetRequestInfo(requestID)
		if err != nil {
			status = http.StatusNotFound
			serverResponse.Message = "Follow request not found"
			utils.RespondJSON(w, status, serverResponse)
			return
		}
	}

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

	// Send notification to the requester about acceptance
	if fr.Notifier != nil && followerID != 0 {
		followeeName, _, err := fr.FollowRequestService.RetrieveUserName(followeeID)
		if err == nil {

			// Store notification in database
			err = fr.FollowRequestService.AddtoNotification(followerID, followeeName+" accepted your follow request")
			if err != nil {
				status = http.StatusInternalServerError
				serverResponse.Message = "Failed to add notification"
				utils.RespondJSON(w, status, serverResponse)
				return
			}

			// Send real-time notification if user is online
			if fr.Notifier.IsOnline(followerID) {
				fr.Notifier.SendNotification(followerID, map[string]interface{}{
					"type":      "notification",
					"subtype":   "follow_request_accepted",
					"user_id":   followeeID,
					"user_name": followeeName,
					"message":   followeeName + " accepted your follow request",
					"timestamp": time.Now().Unix(),
				})
			}
		}
	}

	serverResponse.Message = "Successfully accepted follow request"
	utils.RespondJSON(w, status, serverResponse)
}

func (fr *FollowRequestHandler) CancelFollowRequestRespond(w http.ResponseWriter, r *http.Request) {
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
		// Get follow request info before rejecting (only if notifications are enabled)
		var followerID, followeeID int64
		if fr.Notifier != nil {
			followerID, followeeID, err = fr.FollowRequestService.GetRequestInfo(requestID)
			if err != nil {
				status = http.StatusNotFound
				serverResponse.Message = "Follow request not found"
				utils.RespondJSON(w, status, serverResponse)
				return
			}
		}

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

		// Send notification to the requester about rejection
		if fr.Notifier != nil && followerID != 0 {
			followeeName, _, err := fr.FollowRequestService.RetrieveUserName(followeeID)
			if err == nil {
				// Store notification in database
				err = fr.FollowRequestService.AddtoNotification(followerID, followeeName+" rejected your follow request")
				if err != nil {
					status = http.StatusInternalServerError
					serverResponse.Message = "Failed to add notification"
					utils.RespondJSON(w, status, serverResponse)
					return
				}

				// Send real-time notification if user is online
				if fr.Notifier.IsOnline(followerID) {
					fr.Notifier.SendNotification(followerID, map[string]interface{}{
						"type":      "notification",
						"subtype":   "follow_request_rejected",
						"user_id":   followeeID,
						"user_name": followeeName,
						"message":   followeeName + " rejected your follow request",
						"timestamp": time.Now().Unix(),
					})
				}
			}
		}

		serverResponse.Message = "Successfully rejected follow request"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Handle acceptance
	// Get follow request info before accepting (only if notifications are enabled)
	var followerID, followeeID int64
	if fr.Notifier != nil {
		followerID, followeeID, err = fr.FollowRequestService.GetRequestInfo(requestID)
		if err != nil {
			status = http.StatusNotFound
			serverResponse.Message = "Follow request not found"
			utils.RespondJSON(w, status, serverResponse)
			return
		}
	}

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

	// Send notification to the requester about acceptance
	if fr.Notifier != nil && followerID != 0 {
		followeeName, _, err := fr.FollowRequestService.RetrieveUserName(followeeID)
		if err == nil {

			// Store notification in database
			err = fr.FollowRequestService.AddtoNotification(followerID, followeeName+" accepted your follow request")
			if err != nil {
				status = http.StatusInternalServerError
				serverResponse.Message = "Failed to add notification"
				utils.RespondJSON(w, status, serverResponse)
				return
			}

			// Send real-time notification if user is online
			if fr.Notifier.IsOnline(followerID) {
				fr.Notifier.SendNotification(followerID, map[string]interface{}{
					"type":      "notification",
					"subtype":   "follow_request_accepted",
					"user_id":   followeeID,
					"user_name": followeeName,
					"message":   followeeName + " accepted your follow request",
					"timestamp": time.Now().Unix(),
				})
			}
		}
	}

	serverResponse.Message = "Successfully accepted follow request"
	utils.RespondJSON(w, status, serverResponse)
}

func (fr *FollowRequestHandler) CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	var serverResponse utils.Response
	status := http.StatusOK

	// Get current user id from session (to verify they can cancel this request)
	userID, ok := r.Context().Value(utils.User_id).(int64)
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

	followerID, _, err := fr.FollowRequestService.GetRequestInfo(requestID)
	if err != nil {
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			serverResponse.Message = "Follow request not found"
		} else {
			status = http.StatusInternalServerError
			serverResponse.Message = "Failed to retrieve follow request info"
		}
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Check if the current user is the one who made the follow request
	if followerID != userID {
		serverResponse.Message = "You can only cancel your own follow requests"
		utils.RespondJSON(w, http.StatusForbidden, serverResponse)
		return
	}

	// Cancel the follow request
	err = fr.FollowRequestService.CancelFollowRequest(requestID)
	if err != nil {
		status = http.StatusInternalServerError
		serverResponse.Message = "Failed to cancel follow request"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	serverResponse.Message = "Successfully cancelled follow request"
	utils.RespondJSON(w, status, serverResponse)
}
