package websocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create minimal tables for testing
	_, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE NOT NULL,
            nickname TEXT,
            first_name TEXT,
            last_name TEXT,
            avatar TEXT
        );
        CREATE TABLE sessions (
            id TEXT PRIMARY KEY,
            user_id INTEGER NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
        CREATE TABLE Messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sender_id INTEGER,
            receiver_id INTEGER,
            group_id INTEGER,
            content TEXT,
            created_at DATETIME
        );
        CREATE TABLE Group_Members (
            user_id INTEGER,
            group_id INTEGER,
            is_accepted INTEGER DEFAULT 0
        );
    `)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Insert test user and session
	_, err = db.Exec(`
		INSERT INTO users (id, email, nickname, first_name, last_name, avatar) VALUES
		(123, 'test@example.com', 'testuser', 'Test', 'User', 'test-avatar.jpg');

		INSERT INTO sessions (id, user_id) VALUES ('test-session', 123);
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	return db
}

func TestWebSocketConnection(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create manager with real dependencies
	manager := NewManager(
		NewDBSessionResolver(db),
		NewDBGroupMemberFetcher(db),
		NewDBMessagePersister(db),
	)

	// Create test server with session cookie
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add test session cookie
		r.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
		manager.HandleConnection(w, r)
	}))
	defer server.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Test connection with timeout
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Test message sending
	testMsg := Message{
		Type:    "broadcast",
		Content: "Hello, World!",
	}

	if err := conn.WriteJSON(testMsg); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Give some time for message processing
	time.Sleep(100 * time.Millisecond)
}
