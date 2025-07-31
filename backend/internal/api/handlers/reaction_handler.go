package handlers

import (
	"encoding/json"
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
	userID, ok := r.Context().Value(utils.User_id).(int)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
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
	userID, ok := r.Context().Value(utils.User_id).(int)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	if err := h.service.UnreactToPost(userID, postID); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to unreact to post"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Successfully unreacted to post"})
}

func (h *ReactionHandler) ReactToComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.User_id).(int)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	commentIDStr := r.PathValue("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid comment ID"})
		return
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
	userID, ok := r.Context().Value(utils.User_id).(int)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	commentIDStr := r.PathValue("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid comment ID"})
		return
	}

	if err := h.service.UnreactToComment(userID, commentID); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to unreact to comment"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Successfully unreacted to comment"})
}
