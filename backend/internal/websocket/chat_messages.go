package websocket

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	DB        *sql.DB
	Resolver  *DBSessionResolver
	Persister *DBMessagePersister
	Notifier  *NotificationSender
}

func NewChatHandler(db *sql.DB, resolver *DBSessionResolver, persister *DBMessagePersister, notifier *NotificationSender) *ChatHandler {
	return &ChatHandler{
		DB:        db,
		Resolver:  resolver,
		Persister: persister,
		Notifier:  notifier,
	}
}

// GET /api/messages/private?user=123&limit=50&offset=0
func (h *ChatHandler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	userID, err := h.Resolver.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.ParseInt(r.URL.Query().Get("user"), 10, 64)
	if err != nil || targetID <= 0 {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	msgs, err := h.Persister.FetchPrivateMessagesPaginated(userID, targetID, limit, offset)
	if err != nil {
		http.Error(w, "Could not fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(msgs)
}

// GET /api/messages/group?group=123&limit=50&offset=0
func (h *ChatHandler) GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	_, err := h.Resolver.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.ParseInt(r.URL.Query().Get("group"), 10, 64)
	if err != nil || groupID <= 0 {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	msgs, err := h.Persister.FetchGroupMessagesPaginated(groupID, limit, offset)
	if err != nil {
		http.Error(w, "Could not fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(msgs)
}

// POST /api/groups/invite
func (h *ChatHandler) SendGroupInvite(w http.ResponseWriter, r *http.Request) {
	inviterID, err := h.Resolver.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	type InviteRequest struct {
		GroupID   int64  `json:"group_id"`
		UserID    int64  `json:"user_id"`
		GroupName string `json:"group_name"`
	}

	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err = h.DB.Exec(`INSERT INTO Group_Members (group_id, user_id, invited_by, is_accepted) VALUES (?, ?, ?, 0)`, req.GroupID, req.UserID, inviterID)
	if err != nil {
		http.Error(w, "Could not invite user", http.StatusInternalServerError)
		return
	}

	msg := "You were invited to join group '" + req.GroupName + "'"
	_, err = h.DB.Exec(`INSERT INTO Notifications (user_id, type, message) VALUES (?, 'group_invite', ?)`, req.UserID, msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if h.Notifier != nil && h.Notifier.IsOnline(req.UserID) {
		h.Notifier.SendNotification(req.UserID, map[string]interface{}{
			"type":      "notification",
			"subtype":   "group_invite",
			"message":   msg,
			"group_id":  req.GroupID,
			"timestamp": time.Now().Unix(),
		})
	}

	w.WriteHeader(http.StatusCreated)
}

// GET /api/notifications?limit=20&offset=0
func (h *ChatHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := h.Resolver.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	rows, err := h.DB.Query(`
		SELECT id, type, message, is_read, created_at 
		FROM Notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Notification struct {
		ID        int64     `json:"id"`
		Type      string    `json:"type"`
		Message   string    `json:"message"`
		IsRead    bool      `json:"is_read"`
		CreatedAt time.Time `json:"created_at"`
	}

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err := rows.Scan(&n.ID, &n.Type, &n.Message, &n.IsRead, &n.CreatedAt)
		if err == nil {
			notifications = append(notifications, n)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notifications)
}

// POST /api/notifications/read
func (h *ChatHandler) MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := h.Resolver.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = h.DB.Exec(`UPDATE Notifications SET is_read = 1 WHERE user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Failed to mark as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
