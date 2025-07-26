package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ----------- Interfaces ------------

// SessionResolver resolves the authenticated user ID, nickname, and avatar from the HTTP request.
type SessionResolver interface {
	GetUserFromRequest(r *http.Request) (int64, string, string, error)
}

// GroupMemberFetcher fetches group member IDs from DB
// to allow broadcasting messages to the group.
type GroupMemberFetcher interface {
	GetGroupMemberIDs(groupID string) ([]int64, error)
}

// MessagePersister stores chat messages for retrieval and persistence.
type MessagePersister interface {
	SaveMessage(senderID int64, msg *Message) error
}

// ----------- Client ---------------

type Client struct {
	ID        int64
	Nickname  string
	Avatar    string
	Conn      *websocket.Conn
	Send      chan []byte
	Connected time.Time
}

func NewClient(id int64, nickname string, avatar string, conn *websocket.Conn) *Client {
	return &Client{
		ID:        id,
		Nickname:  nickname,
		Avatar:    avatar,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Connected: time.Now(),
	}
}

// ----------- Manager ---------------

type Manager struct {
	clients    map[int64]*Client
	mu         sync.RWMutex
	resolver   SessionResolver
	groupQuery GroupMemberFetcher
	persister  MessagePersister
}

func NewManager(resolver SessionResolver, groupFetcher GroupMemberFetcher, persister MessagePersister) *Manager {
	return &Manager{
		clients:    make(map[int64]*Client),
		resolver:   resolver,
		groupQuery: groupFetcher,
		persister:  persister,
	}
}

func (m *Manager) Register(c *Client) {
	// Add client to map first
	m.mu.Lock()
	m.clients[c.ID] = c
	fmt.Printf("User %d connected as %s", c.ID, c.Nickname)
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
			"nickname": client.Nickname,
			"avatar":   client.Avatar,
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
			"nickname": client.Nickname,
			"connected": client.Connected.Unix(),
			"avatar":    client.Avatar,
		})
	}

	return users
}

// ----------- Message Format ---------------

type Message struct {
	Type      string `json:"type"`
	To        int64  `json:"to,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func parseMessage(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return &m, err
}

// ----------- Handler Entry Point ---------------

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade", http.StatusBadRequest)
		return
	}

	userID, nickname, avatar, err := m.resolver.GetUserFromRequest(r)
	if err != nil {
		conn.Close()
		return
	}

	client := NewClient(userID, nickname, avatar, conn)
	m.Register(client)
	defer m.Unregister(userID)
	defer conn.Close()

	go writePump(client)
	readPump(m, client)
}

// ----------- read/write loop ---------------

func readPump(m *Manager, c *Client) {
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
		encoded, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		// Save to DB if persister is configured
		if m.persister != nil {
			_ = m.persister.SaveMessage(c.ID, msg)
		}

		switch msg.Type {
		case "private":
			m.SendToUser(msg.To, encoded)
		case "group":
			m.BroadcastToGroup(c.ID, msg.GroupID, encoded)
		case "broadcast":
			m.BroadcastToAll(encoded)
		}
	}
}

func writePump(c *Client) {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
