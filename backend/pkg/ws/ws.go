package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ----------- Interfaces ------------

// SessionResolver resolves the authenticated user ID from the HTTP request.
type SessionResolver interface {
	GetUserIDFromRequest(r *http.Request) (int64, error)
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
	Conn      *websocket.Conn
	Send      chan []byte
	Connected time.Time
}

func NewClient(id int64, conn *websocket.Conn) *Client {
	return &Client{
		ID:        id,
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
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.ID] = c
	log.Printf("User %d connected", c.ID)
}

func (m *Manager) Unregister(id int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[id]; ok {
		delete(m.clients, id)
		log.Printf("User %d disconnected", id)
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

	userID, err := m.resolver.GetUserIDFromRequest(r)
	if err != nil {
		conn.Close()
		return
	}

	client := NewClient(userID, conn)
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
