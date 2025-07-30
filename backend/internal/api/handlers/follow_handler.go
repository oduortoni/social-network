package handlers

import (
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

type FollowHandler struct {
	FollowService service.FollowServiceInterface
	Notifier      *ws.NotificationSender
}

func NewFollowHandler(followService service.FollowServiceInterface, notifier *ws.NotificationSender) *FollowHandler {
	return &FollowHandler{
		FollowService: followService,
		Notifier:      notifier,
	}
}

func (follow *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var serverResponse utils.Response
	status := http.StatusOK

	// Get current user id from session
	followerId, ok := r.Context().Value(utils.User_id).(int64)
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
		status = http.StatusBadRequest // ✅ Corrected from 500
		serverResponse.Message = "Invalid JSON format"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	isFolloweeAccountPublic, err := follow.FollowService.IsAccountPublic(int64(followee.FolloweeId))
	if err != nil {
		status = http.StatusInternalServerError
		serverResponse.Message = "Failed to check account privacy status"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	if isFolloweeAccountPublic {
		err = follow.FollowService.CreateFollowForPublicAccount(followerId, int64(followee.FolloweeId))
		if err != nil {
			status = http.StatusInternalServerError
			serverResponse.Message = "Failed to create follow connection"
			utils.RespondJSON(w, status, serverResponse)
			return
		}

		// Send notification for public account follow
		if follow.Notifier != nil {
			followerName, _, err := follow.FollowService.GetUserInfo(followerId)
			if err == nil {
				// Store notification in database
				err = follow.FollowService.AddtoNotification(int64(followee.FolloweeId), followerName+" started following you")
				if err != nil {
					status = http.StatusInternalServerError
					serverResponse.Message = "Failed to add notification"
					utils.RespondJSON(w, status, serverResponse)
					return
				}

				// Send real-time notification if user is online
				if follow.Notifier.IsOnline(int64(followee.FolloweeId)) {
					follow.Notifier.SendNotification(int64(followee.FolloweeId), map[string]interface{}{
						"type":      "notification",
						"subtype":   "follow",
						"user_id":   followerId,
						"user_name": followerName,
						"message":   followerName + " started following you",
						"timestamp": time.Now().Unix(),
					})
				}
			}
		}

		serverResponse.Message = "Account successfully followed"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Handle private account response
	followID, err := follow.FollowService.CreateFollowForPrivateAccount(followerId, int64(followee.FolloweeId))
	if err != nil {
		status = http.StatusInternalServerError
		serverResponse.Message = "Failed to create follow connection"
		utils.RespondJSON(w, status, serverResponse)
		return
	}

	// Send notification for follow request
	if follow.Notifier != nil {
		followerName, _, err := follow.FollowService.GetUserInfo(followerId)
		if err == nil {
			// Store notification in database
			err = follow.FollowService.AddtoNotification(int64(followee.FolloweeId), followerName+" sent you a follow request")
			if err != nil {
				status = http.StatusInternalServerError
				serverResponse.Message = "Failed to add notification"
				utils.RespondJSON(w, status, serverResponse)
				return
			}

			// Send real-time notification if user is online
			if follow.Notifier.IsOnline(int64(followee.FolloweeId)) {
				follow.Notifier.SendNotification(int64(followee.FolloweeId), map[string]interface{}{
					"type":      "notification",
					"subtype":   "follow_request",
					"user_id":   followerId,
					"user_name": followerName,
					"message":   followerName + " sent you a follow request",
					"timestamp": time.Now().Unix(),
					"additional_data": map[string]interface{}{
						"request_id": followID,
						"actions":    []string{"accept", "reject"},
					},
				})
			}
		}
	}

	serverResponse.Message = "Follow request sent. You will be able to follow once approved."
	utils.RespondJSON(w, status, serverResponse)
}

func (f *FollowHandler) FollowStats(w http.ResponseWriter, r *http.Request) {
	userIdstr := r.PathValue("userid")
	userId, err := strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid User Id"})
		return
	}
	followers, following, err := f.FollowService.GetFollowFollowingStat(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	var followfollowingstat models.FollowFollowingStat
	followfollowingstat.NumberOfFollowers = followers
	followfollowingstat.NumberOfFollowing = following

	utils.RespondJSON(w, http.StatusOK, followfollowingstat)
}

func (f *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userIdstr := r.PathValue("userid")
	userId, err := strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid User Id"})
		return
	}
	followers, err := f.FollowService.GetFollowersList(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	utils.RespondJSON(w, http.StatusOK, followers)
}

func (f *FollowHandler) GetFollowees(w http.ResponseWriter, r *http.Request) {
	userIdstr := r.PathValue("userid")
	userId, err := strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid User Id"})
		return
	}
	followers, err := f.FollowService.GetFolloweesList(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	utils.RespondJSON(w, http.StatusOK, followers)
}
