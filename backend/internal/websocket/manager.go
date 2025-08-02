package websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Manager struct {
	clients    map[int64]*Client
	mu         sync.RWMutex
	Resolver   SessionResolver
	groupQuery GroupMemberFetcher
	persister  MessagePersister
	DB         *sql.DB
}

func NewManager(resolver SessionResolver, groupFetcher GroupMemberFetcher, persister MessagePersister, db *sql.DB) *Manager {
	return &Manager{
		clients:    make(map[int64]*Client),
		Resolver:   resolver,
		groupQuery: groupFetcher,
		persister:  persister,
		DB:         db,
	}
}

func (m *Manager) Register(c *Client) {
	// Add client to map first
	m.mu.Lock()
	m.clients[c.ID] = c
	fmt.Printf("User %d connected as %s\n", c.ID, c.Nickname)
	m.mu.Unlock() // Release lock before broadcasting

	// Broadcast user connection notification to all other users
	connectionNotification := map[string]interface{}{
		"type":      "notification",
		"subtype":   "user_connected",
		"user_id":   c.ID,
		"nickname":  c.Nickname,
		"avatar":    c.Avatar,
		"message":   c.Nickname + " is now online",
		"timestamp": time.Now().Unix(),
	}

	m.broadcastNotificationToAll(connectionNotification, c.ID)
}

func (m *Manager) Unregister(id int64) {
	var disconnectionNotification map[string]interface{}

	// Remove client and prepare notification
	m.mu.Lock()
	if client, ok := m.clients[id]; ok {
		// Prepare notification before removing client
		disconnectionNotification = map[string]interface{}{
			"type":      "notification",
			"subtype":   "user_disconnected",
			"user_id":   id,
			"nickname":  client.Nickname,
			"avatar":    client.Avatar,
			"message":   client.Nickname + " went offline",
			"timestamp": time.Now().Unix(),
		}

		delete(m.clients, id)
		fmt.Printf("User %d disconnected", id)
	}
	m.mu.Unlock() // Release lock before broadcasting

	// Broadcast after releasing the lock
	if disconnectionNotification != nil {
		m.broadcastNotificationToAll(disconnectionNotification, id)
	}
}

// ----------- read/write loop ---------------

func (m *Manager) ReadPump(c *Client) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		msg, err := parseMessage(data)
		if err != nil {
			continue
		}

		msg.Timestamp = time.Now().Unix()
		msg.From = c.ID // Add sender's ID to the message
		encoded, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		switch msg.Type {
		case "private":
			// Requirement #2 & #4: Validate if users are allowed to chat
			allowed, err := m.canUsersChat(c.ID, msg.To)
			if err != nil {
				log.Printf("Error checking chat permissions for user %d to %d: %v", c.ID, msg.To, err)
				continue // Silently drop on error
			}

			if allowed {
				// Requirement #3: Save to DB if persister is configured
				if m.persister != nil {
					_ = m.persister.SaveMessage(c.ID, msg)
				}
				// Requirement #3: Forward to recipient and sender
				m.SendToUser(msg.To, encoded)
				m.SendToUser(c.ID, encoded)
			} else {
				// Requirement #4: Return error response for unauthorized message
				errorMsg, _ := json.Marshal(Message{
					Type:      "error",
					Content:   "You are not permitted to message this user.",
					Timestamp: time.Now().Unix(),
				})
				m.SendToUser(c.ID, errorMsg)
			}
		case "group":
			m.BroadcastToGroup(c.ID, msg.GroupID, encoded)
		case "broadcast":
			m.BroadcastToAll(encoded)
		}
	}
}

func (m *Manager) canUsersChat(userA, userB int64) (bool, error) {
	// 1. check that two users have a mutual following
	var count int
	err := m.DB.QueryRow(`
		SELECT COUNT(*) FROM Followers
		WHERE (follower_id = ? AND followee_id = ? AND status = 'accepted')
		   OR (follower_id = ? AND followee_id = ? AND status = 'accepted')
	`, userA, userB, userB, userA).Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if count == 2 { // mutually follow each other
		return true, nil
	}

	// 2. does the receiver has a public profile
	var isPublic bool
	err = m.DB.QueryRow(`SELECT isprofilepublic FROM Users WHERE id = ?`, userB).Scan(&isPublic)
	if err != nil {
		return false, err
	}

	return isPublic, nil
}

func (m *Manager) WritePump(c *Client) {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

// ----------- Message Broadcasting Methods ---------------

func (m *Manager) SendToUser(id int64, msg []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.clients[id]; ok {
		c.Send <- msg
	}
}

func (m *Manager) BroadcastToGroup(sender int64, groupID string, msg []byte) {
	ids, err := m.groupQuery.GetGroupMemberIDs(groupID)
	if err != nil {
		return
	}
	for _, id := range ids {
		if id == sender {
			continue
		}
		m.SendToUser(id, msg)
	}
}

func (m *Manager) BroadcastToAll(msg []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, client := range m.clients {
		client.Send <- msg
	}
}

func (m *Manager) IsOnline(userID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.clients[userID]
	return exists
}

func (m *Manager) OnlineUserIDs() []int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]int64, 0, len(m.clients))
	for id := range m.clients {
		ids = append(ids, id)
	}
	return ids
}

// Notification Broadcasting Methods

func (m *Manager) broadcastNotificationToAll(notification map[string]interface{}, excludeUserID int64) {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	for id, client := range m.clients {
		if id != excludeUserID {
			select {
			case client.Send <- notificationBytes:
			default:
				// Channel is full, skip this client
				log.Printf("Failed to send notification to user %d: channel full", id)
			}
		}
	}
}

func (m *Manager) SendNotificationToUser(userID int64, notification map[string]interface{}) {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	if client, ok := m.clients[userID]; ok {
		select {
		case client.Send <- notificationBytes:
		default:
			log.Printf("Failed to send notification to user %d: channel full", userID)
		}
	}
}

func (m *Manager) SendNotificationToGroup(groupID string, notification map[string]interface{}, excludeUserID int64) {
	memberIDs, err := m.groupQuery.GetGroupMemberIDs(groupID)
	if err != nil {
		log.Printf("Failed to get group members: %v", err)
		return
	}

	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, memberID := range memberIDs {
		if memberID != excludeUserID {
			if client, ok := m.clients[memberID]; ok {
				select {
				case client.Send <- notificationBytes:
				default:
					log.Printf("Failed to send notification to user %d: channel full", memberID)
				}
			}
		}
	}
}

func (m *Manager) GetOnlineUsers() []map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]map[string]interface{}, 0, len(m.clients))
	for _, client := range m.clients {
		users = append(users, map[string]interface{}{
			"user_id":   client.ID,
			"nickname":  client.Nickname,
			"connected": client.Connected.Unix(),
			"avatar":    client.Avatar,
		})
	}

	return users
}
