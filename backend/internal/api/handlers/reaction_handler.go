package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
)

// ReactionHandler handles reaction-related requests.	ype ReactionHandler struct {
	service *service.ReactionService
}

// NewReactionHandler creates a new ReactionHandler.
func NewReactionHandler(service *service.ReactionService) *ReactionHandler {
	return &ReactionHandler{service}
}

// RegisterReactionRoutes registers the reaction routes to the router.
func (h *ReactionHandler) RegisterReactionRoutes(router *mux.Router) {
	router.HandleFunc("/api/posts/{id}/reaction", h.ReactToPost).Methods("POST")
	router.HandleFunc("/api/posts/{id}/reaction", h.UnreactToPost).Methods("DELETE")
	router.HandleFunc("/api/comments/{id}/reaction", h.ReactToComment).Methods("POST")
	router.HandleFunc("/api/comments/{id}/reaction", h.UnreactToComment).Methods("DELETE")
}

func (h *ReactionHandler) ReactToPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var reaction models.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from session
	reaction.UserID = 1
	reaction.PostID = &postID

	if err := h.service.ReactToPost(&reaction); err != nil {
		http.Error(w, "Failed to react to post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ReactionHandler) UnreactToPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from session
	userID := 1

	if err := h.service.UnreactToPost(userID, postID); err != nil {
		http.Error(w, "Failed to unreact to post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ReactionHandler) ReactToComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var reaction models.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from session
	reaction.UserID = 1
	reaction.CommentID = &commentID

	if err := h.service.ReactToComment(&reaction); err != nil {
		http.Error(w, "Failed to react to comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ReactionHandler) UnreactToComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from session
	userID := 1

	if err := h.service.UnreactToComment(userID, commentID); err != nil {
		http.Error(w, "Failed to unreact to comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
