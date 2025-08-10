package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
)

type GroupHandler struct {
	groupService service.GroupService
	groupRequestService service.GroupRequestService
}

func NewGroupHandler(groupService service.GroupService, groupRequestService service.GroupRequestService) *GroupHandler {
	return &GroupHandler{groupService: groupService, groupRequestService: groupRequestService}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// The user ID should be extracted from the request context,
	// which is populated by your authentication middleware.
	creatorID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	group.CreatorID = creatorID

	// Set default privacy if not provided
	if group.Privacy == "" {
		group.Privacy = "public"
	}

	newGroup, err := h.groupService.CreateGroup(&group)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGroup)
}

func (h *GroupHandler) SendJoinRequest(w http.ResponseWriter, r *http.Request) {
	groupIDStr := r.PathValue("groupID")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	request, err := h.groupRequestService.SendJoinRequest(groupID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send join request: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

func (h *GroupHandler) ApproveJoinRequest(w http.ResponseWriter, r *http.Request) {
	requestIDStr := r.PathValue("requestID")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	approverID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Approver ID not found in context", http.StatusUnauthorized)
		return
	}

	err = h.groupRequestService.ApproveJoinRequest(requestID, approverID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to approve join request: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Join request approved"})
}

func (h *GroupHandler) RejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	requestIDStr := r.PathValue("requestID")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	rejecterID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Rejecter ID not found in context", http.StatusUnauthorized)
		return
	}

	err = h.groupRequestService.RejectJoinRequest(requestID, rejecterID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reject join request: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Join request rejected"})
}
