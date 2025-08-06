package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// ReactionHandler handles reaction-related requests.
type ReactionHandler struct {
	service *service.ReactionService
}

// NewReactionHandler creates a new ReactionHandler.
func NewReactionHandler(service *service.ReactionService) *ReactionHandler {
	return &ReactionHandler{service}
}

func (h *ReactionHandler) ReactToPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}
	fmt.Println("got post if from url: ", postID)
	userIDValue := r.Context().Value(utils.User_id)
	if userIDValue == nil {
		fmt.Println("unable to get user from context: ", r.Context())
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		if userID64, ok := userIDValue.(int64); ok {
			userID = int(userID64)
		} else {
			fmt.Println("unable to cast user ID from context, type:", fmt.Sprintf("%T", userIDValue))
			utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
			return
		}
	}

	var reaction models.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid request body"})
		return
	}

	reaction.UserID = userID
	reaction.PostID = &postID

	if err := h.service.ReactToPost(&reaction); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to react to post"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Successfully reacted to post"})
}

func (h *ReactionHandler) UnreactToPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}
	userIDValue := r.Context().Value(utils.User_id)
	if userIDValue == nil {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		if userID64, ok := userIDValue.(int64); ok {
			userID = int(userID64)
		} else {
			utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
			return
		}
	}

	if err := h.service.UnreactToPost(userID, postID); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to unreact to post"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Successfully unreacted to post"})
}

func (h *ReactionHandler) ReactToComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.PathValue("commentId")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid comment ID"})
		return
	}
	userIDValue := r.Context().Value(utils.User_id)
	if userIDValue == nil {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		if userID64, ok := userIDValue.(int64); ok {
			userID = int(userID64)
		} else {
			utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
			return
		}
	}

	var reaction models.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid request body"})
		return
	}

	reaction.UserID = userID
	reaction.CommentID = &commentID

	if err := h.service.ReactToComment(&reaction); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to react to comment"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Successfully reacted to comment"})
}

func (h *ReactionHandler) UnreactToComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.PathValue("commentId")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid comment ID"})
		return
	}
	userIDValue := r.Context().Value(utils.User_id)
	if userIDValue == nil {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		if userID64, ok := userIDValue.(int64); ok {
			userID = int(userID64)
		} else {
			utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
			return
		}
	}

	if err := h.service.UnreactToComment(userID, commentID); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to unreact to comment"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Successfully unreacted to comment"})
}
