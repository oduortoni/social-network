package websocket

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestServer(t *testing.T) (*httptest.Server, *sql.DB, *Manager) {
	// Setup test database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create test tables
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			nickname TEXT,
			first_name TEXT,
			last_name TEXT,
			avatar TEXT,
			date_of_birth DATE,
			about_me TEXT,
			is_profile_public INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE Messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender_id INTEGER,
			receiver_id INTEGER,
			group_id INTEGER,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE Group_Members (
			user_id INTEGER,
			group_id INTEGER,
			is_accepted INTEGER DEFAULT 0,
			invited_by INTEGER
		);
		CREATE TABLE Notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			type TEXT,
			message TEXT,
			is_read INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE Groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT,
			creator_id INTEGER
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO users (id, email, nickname, first_name, last_name, avatar) VALUES
		(1, 'user1@test.com', 'testuser1', 'Test', 'User1', 'avatar1.jpg'),
		(2, 'user2@test.com', 'testuser2', 'Test', 'User2', 'avatar2.jpg'),
		(3, 'user3@test.com', 'testuser3', 'Test', 'User3', 'avatar3.jpg');

		INSERT INTO sessions (id, user_id, expires_at) VALUES
		('test-session-1', 1, datetime('now', '+1 day')),
		('test-session-2', 2, datetime('now', '+1 day'));

		INSERT INTO Groups (id, title, description, creator_id) VALUES (1, 'Test Group', 'A test group', 1);

		INSERT INTO Group_Members (user_id, group_id, is_accepted) VALUES
		(1, 1, 1), (2, 1, 1);
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Create WebSocket manager
	manager := NewManager(
		NewDBSessionResolver(db),
		NewDBGroupMemberFetcher(db),
		NewDBMessagePersister(db),
	)

	// Create chat handler
	notifier := NewDBNotificationSender(manager)
	chatHandler := NewChatHandler(db, NewDBSessionResolver(db), NewDBMessagePersister(db), notifier, manager)

	// Create test server with routes
	mux := http.NewServeMux()
	
	// WebSocket endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Add test session cookie for authentication
		if r.Header.Get("Cookie") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		manager.HandleConnection(w, r)
	})
	
	// HTTP API endpoints
	mux.HandleFunc("/api/messages/private", chatHandler.GetPrivateMessages)
	mux.HandleFunc("/api/messages/group", chatHandler.GetGroupMessages)
	mux.HandleFunc("/api/groups/invite", chatHandler.SendGroupInvite)
	mux.HandleFunc("/api/notifications", chatHandler.GetNotifications)

	server := httptest.NewServer(mux)
	return server, db, manager
}

func TestWebSocketEndpoint(t *testing.T) {
	server, db, _ := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Test 1: Unauthorized access
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	
	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	_, _, err := dialer.Dial(wsURL, nil)
	if err == nil {
		t.Error("Expected unauthorized access to fail")
	}

	// Test 2: Authorized access
	headers := http.Header{}
	headers.Set("Cookie", "session_id=test-session-1")
	
	conn, _, err := dialer.Dial(wsURL, headers)
	if err != nil {
		t.Fatalf("Failed to connect with valid session: %v", err)
	}
	defer conn.Close()

	// Test 3: Send a private message
	privateMsg := Message{
		Type:    "private",
		To:      2,
		Content: "Hello from user 1!",
	}
	
	if err := conn.WriteJSON(privateMsg); err != nil {
		t.Fatalf("Failed to send private message: %v", err)
	}

	// Test 4: Send a group message
	groupMsg := Message{
		Type:    "group",
		GroupID: "1",
		Content: "Hello group!",
	}
	
	if err := conn.WriteJSON(groupMsg); err != nil {
		t.Fatalf("Failed to send group message: %v", err)
	}

	// Test 5: Send a broadcast message
	broadcastMsg := Message{
		Type:    "broadcast",
		Content: "Hello everyone!",
	}
	
	if err := conn.WriteJSON(broadcastMsg); err != nil {
		t.Fatalf("Failed to send broadcast message: %v", err)
	}

	// Give time for message processing
	time.Sleep(100 * time.Millisecond)
	
	t.Log("WebSocket endpoint tests passed!")
}

func TestHTTPAPIEndpoints(t *testing.T) {
	server, db, _ := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Insert test messages
	_, err := db.Exec(`
		INSERT INTO Messages (sender_id, receiver_id, content, created_at) VALUES 
		(1, 2, 'Private message 1', datetime('now', '-1 hour')),
		(2, 1, 'Private message 2', datetime('now', '-30 minutes'));
		
		INSERT INTO Messages (sender_id, group_id, content, created_at) VALUES 
		(1, 1, 'Group message 1', datetime('now', '-45 minutes')),
		(2, 1, 'Group message 2', datetime('now', '-15 minutes'));
	`)
	if err != nil {
		t.Fatalf("Failed to insert test messages: %v", err)
	}

	// Test 1: Get private messages
	req, _ := http.NewRequest("GET", server.URL+"/api/messages/private?user=2&limit=10", nil)
	req.Header.Set("Cookie", "session_id=test-session-1")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get private messages: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test 2: Get group messages
	req, _ = http.NewRequest("GET", server.URL+"/api/messages/group?group=1&limit=10", nil)
	req.Header.Set("Cookie", "session_id=test-session-1")
	
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get group messages: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test 3: Send group invite
	inviteData := map[string]interface{}{
		"group_id":   1,
		"user_id":    3,
		"group_name": "Test Group",
	}
	jsonData, _ := json.Marshal(inviteData)
	
	req, _ = http.NewRequest("POST", server.URL+"/api/groups/invite", bytes.NewBuffer(jsonData))
	req.Header.Set("Cookie", "session_id=test-session-1")
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send group invite: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// Test 4: Get notifications
	req, _ = http.NewRequest("GET", server.URL+"/api/notifications?limit=10", nil)
	req.Header.Set("Cookie", "session_id=test-session-1")
	
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get notifications: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Log("HTTP API endpoint tests passed!")
}

func TestMessagePersistence(t *testing.T) {
	server, db, manager := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Connect two clients
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	
	headers1 := http.Header{}
	headers1.Set("Cookie", "session_id=test-session-1")
	conn1, _, err := dialer.Dial(wsURL, headers1)
	if err != nil {
		t.Fatalf("Failed to connect client 1: %v", err)
	}
	defer conn1.Close()

	headers2 := http.Header{}
	headers2.Set("Cookie", "session_id=test-session-2")
	conn2, _, err := dialer.Dial(wsURL, headers2)
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer conn2.Close()

	// Wait for connections to be established
	time.Sleep(100 * time.Millisecond)

	// Verify both users are online
	if !manager.IsOnline(1) || !manager.IsOnline(2) {
		t.Error("Both users should be online")
	}

	// Send a private message
	privateMsg := Message{
		Type:    "private",
		To:      2,
		Content: "Test persistence message",
	}
	
	if err := conn1.WriteJSON(privateMsg); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Wait for message processing
	time.Sleep(200 * time.Millisecond)

	// Verify message was persisted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Messages WHERE sender_id = 1 AND receiver_id = 2 AND content = ?", 
		"Test persistence message").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query messages: %v", err)
	}
	
	if count != 1 {
		t.Errorf("Expected 1 persisted message, got %d", count)
	}

	t.Log("Message persistence test passed!")
}
