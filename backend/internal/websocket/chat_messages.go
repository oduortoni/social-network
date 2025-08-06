package websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	DB                *sql.DB
	Resolver          *DBSessionResolver
	Persister         *DBMessagePersister
	Notifier          *NotificationSender
	WSManager         *Manager
	PermissionChecker PermissionChecker
}

func NewChatHandler(db *sql.DB, resolver *DBSessionResolver, persister *DBMessagePersister, notifier *NotificationSender, wsManager *Manager, permissionChecker PermissionChecker) *ChatHandler {
	return &ChatHandler{
		DB:                db,
		Resolver:          resolver,
		Persister:         persister,
		Notifier:          notifier,
		WSManager:         wsManager,
		PermissionChecker: permissionChecker,
	}
}

// GET /api/messages/private?user=123&limit=50&offset=0
func (h *ChatHandler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	userID, _, _, err := h.Resolver.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.ParseInt(r.URL.Query().Get("user"), 10, 64)
	if err != nil || targetID <= 0 {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	// Requirement #2 & #4: Validate that the users are allowed to chat
	allowed, err := h.PermissionChecker.CanUsersChat(userID, targetID)
	if err != nil {
		http.Error(w, "Could not verify relationship", http.StatusInternalServerError)
		return
	}
	if !allowed {
		// Instead of blocking access, return empty history to allow new conversations
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]interface{}{})
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
	_, _, _, err := h.Resolver.GetUserFromRequest(r)
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
	inviterID, _, _, err := h.Resolver.GetUserFromRequest(r)
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
	userID, _, _, err := h.Resolver.GetUserFromRequest(r)
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
	userID, _, _, err := h.Resolver.GetUserFromRequest(r)
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

// GET /api/users/online
func (h *ChatHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetOnlineUsers API called")
	onlineUsers := h.WSManager.GetOnlineUsers()
	fmt.Printf("Found %d online users: %+v\n", len(onlineUsers), onlineUsers)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"online_users": onlineUsers,
		"count":        len(onlineUsers),
	})
}

// GET /api/users/messageable

/*
/* GetMessageableUsers: finds all users that the current user can message
/* It combines two groups:
/* 1. Mutual followers.
/* 2. Users with public profiles.
/* The UNION operator automatically handles duplicates.
*/
func (h *ChatHandler) GetMessageableUsers(w http.ResponseWriter, r *http.Request) {
	userID, _, _, err := h.Resolver.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// [ RESTRICTIVE: IF FOLLOW AND PUBLIC PROFILE WORK AND MUTUAL FOLLOWERS ]
	// ===================================================================================
	// rows, err := h.DB.Query(`
	// 	SELECT u.id, u.nickname, u.avatar
	// 	FROM Users u
	// 	INNER JOIN Followers f1 ON u.id = f1.followee_id AND f1.follower_id = ? AND f1.status = 'accepted'
	// 	INNER JOIN Followers f2 ON u.id = f2.follower_id AND f2.followee_id = ? AND f2.status = 'accepted'
	// 	WHERE u.id != ?
	// 	UNION
	// 	SELECT u.id, u.nickname, u.avatar
	// 	FROM Users u
	// 	WHERE u.is_profile_public = 1 AND u.id != ?
	// `, userID, userID, userID, userID)
	// if err != nil {
	// 	http.Error(w, "Failed to fetch messageable users", http.StatusInternalServerError)
	// 	return
	// }

	// [ PERMISSIBLE PERMISSIONS ]
	rows, err := h.DB.Query(`
		SELECT u.id, u.nickname, u.avatar
		FROM Users u
		WHERE u.id != ?
	`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch messageable users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MessageableUser struct {
		ID       int64  `json:"id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	var users []MessageableUser
	for rows.Next() {
		fmt.Println("Scanning row")
		var u MessageableUser
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Avatar); err == nil {
			users = append(users, u)
			fmt.Printf("Found messageable user: %+v\n", u)
		}
	}

	fmt.Printf("Found %d messageable users: %+v\n", len(users), users)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}
