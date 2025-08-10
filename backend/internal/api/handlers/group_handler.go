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
	groupChatMessageService service.GroupChatMessageService
}

func NewGroupHandler(groupService service.GroupService, groupRequestService service.GroupRequestService, groupChatMessageService service.GroupChatMessageService) *GroupHandler {
	return &GroupHandler{groupService: groupService, groupRequestService: groupRequestService, groupChatMessageService: groupChatMessageService}
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

func (h *GroupHandler) SendGroupChatMessage(w http.ResponseWriter, r *http.Request) {
	groupIDStr := r.PathValue("groupID")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	senderID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Sender ID not found in context", http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.groupChatMessageService.SendGroupChatMessage(groupID, senderID, requestBody.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *GroupHandler) GetGroupChatMessages(w http.ResponseWriter, r *http.Request) {
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

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	messages, err := h.groupChatMessageService.GetGroupChatMessages(groupID, userID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get messages: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
