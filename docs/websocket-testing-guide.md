# WebSocket Testing Guide

## Overview

This guide covers testing strategies for the WebSocket implementation, including unit tests, integration tests, and manual testing procedures.

## Test Structure

The WebSocket package includes comprehensive tests:

```
backend/internal/websocket/
├── ws_test.go                    # Unit tests for core functionality
├── integration_test.go           # Integration tests with database
└── websocket_routes_test.go      # Route and endpoint testing
```

## Running Tests

### Run All WebSocket Tests
```bash
cd backend
go test ./pkg/ws/... -v
```

### Run Specific Test Files
```bash
# Unit tests only
go test ./pkg/ws -run TestManager -v

# Integration tests only  
go test ./pkg/ws -run TestWebSocketConnection -v

# Route tests only
go test ./pkg/ws -run TestWebSocketEndpoint -v
```

### Run with Coverage
```bash
go test ./pkg/ws/... -cover -v
```

## Unit Tests

### Manager Registration Tests
Tests basic client registration and unregistration:

```go
func TestManagerRegisterUnregister(t *testing.T) {
    manager := NewManager(nil, nil, nil)
    client := NewClient(123, nil)
    
    // Test registration
    manager.Register(client)
    if !manager.IsOnline(123) {
        t.Error("Client should be online after registration")
    }
    
    // Test unregistration
    manager.Unregister(123)
    if manager.IsOnline(123) {
        t.Error("Client should be offline after unregistration")
    }
}
```

### Broadcast Functionality Tests
Tests message broadcasting to multiple clients:

```go
func TestManagerBroadcast(t *testing.T) {
    manager := NewManager(nil, nil, nil)
    
    // Create mock clients
    client1 := NewClient(1, nil)
    client2 := NewClient(2, nil)
    
    manager.Register(client1)
    manager.Register(client2)
    
    // Test online status
    if !manager.IsOnline(1) || !manager.IsOnline(2) {
        t.Error("Clients should be online after registration")
    }
    
    // Test getting online user IDs
    ids := manager.OnlineUserIDs()
    if len(ids) != 2 {
        t.Errorf("Expected 2 online users, got %d", len(ids))
    }
}
```

## Integration Tests

### Full WebSocket Connection Test
Tests complete WebSocket lifecycle with database integration:

```go
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
        r.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
        manager.HandleConnection(w, r)
    }))
    defer server.Close()
    
    // Test WebSocket connection
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
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
}
```

### Message Persistence Test
Tests that messages are properly stored in the database:

```go
func TestMessagePersistence(t *testing.T) {
    server, db, manager := setupTestServer(t)
    defer server.Close()
    defer db.Close()
    
    // Connect clients and send messages
    // ... connection setup ...
    
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
}
```

## Route Testing

### HTTP API Endpoint Tests
Tests all HTTP endpoints with proper authentication:

```go
func TestHTTPAPIEndpoints(t *testing.T) {
    server, db, _ := setupTestServer(t)
    defer server.Close()
    defer db.Close()
    
    // Test private messages endpoint
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
    
    // Test other endpoints...
}
```

### Authentication Tests
Tests that endpoints properly handle authentication:

```go
func TestWebSocketAuthentication(t *testing.T) {
    // Test unauthorized access
    wsURL := "ws://localhost:9000/ws"
    _, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err == nil {
        t.Error("Expected unauthorized access to fail")
    }
    
    // Test authorized access
    headers := http.Header{}
    headers.Set("Cookie", "session_id=valid-session")
    
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
    if err != nil {
        t.Fatalf("Failed to connect with valid session: %v", err)
    }
    defer conn.Close()
}
```

## Manual Testing

### WebSocket Connection Testing

1. **Start the server:**
```bash
cd backend
go run server.go
```

2. **Test with browser console:**
```javascript
// Open browser console and connect
const ws = new WebSocket('ws://localhost:9000/ws');

ws.onopen = () => console.log('Connected');
ws.onmessage = (event) => console.log('Received:', JSON.parse(event.data));
ws.onerror = (error) => console.log('Error:', error);

// Send test messages
ws.send(JSON.stringify({
    type: "broadcast",
    content: "Hello from browser!"
}));
```

3. **Test with curl for HTTP endpoints:**
```bash
# Test notifications endpoint
curl -X GET "http://localhost:9000/api/notifications" \
  -H "Cookie: session_id=your-session-id" \
  -v

# Test group invite
curl -X POST "http://localhost:9000/api/groups/invite" \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=your-session-id" \
  -d '{"group_id": 1, "user_id": 2, "group_name": "Test Group"}'
```

### Load Testing

For performance testing, you can use tools like:

1. **WebSocket Load Testing with Artillery:**
```yaml
# artillery-websocket.yml
config:
  target: 'ws://localhost:9000'
  phases:
    - duration: 60
      arrivalRate: 10

scenarios:
  - name: "WebSocket messaging"
    engine: ws
    flow:
      - connect:
          url: "/ws"
      - send:
          payload: '{"type":"broadcast","content":"Load test message"}'
      - think: 1
```

2. **HTTP Load Testing:**
```bash
# Using Apache Bench
ab -n 1000 -c 10 -H "Cookie: session_id=test" \
  http://localhost:9000/api/notifications

# Using curl in a loop
for i in {1..100}; do
  curl -s "http://localhost:9000/api/notifications" \
    -H "Cookie: session_id=test" > /dev/null
done
```

## Test Data Setup

### Database Test Setup
```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    
    // Create test tables
    _, err = db.Exec(`
        CREATE TABLE sessions (
            id TEXT PRIMARY KEY,
            user_id INTEGER NOT NULL,
            expires_at DATETIME
        );
        CREATE TABLE Messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sender_id INTEGER,
            receiver_id INTEGER,
            group_id INTEGER,
            content TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        -- Add other required tables...
    `)
    
    if err != nil {
        t.Fatalf("Failed to create test tables: %v", err)
    }
    
    // Insert test data
    _, err = db.Exec(`
        INSERT INTO sessions (id, user_id, expires_at) VALUES 
        ('test-session-1', 1, datetime('now', '+1 day')),
        ('test-session-2', 2, datetime('now', '+1 day'));
    `)
    
    return db
}
```

## Mocking for Unit Tests

### Mock WebSocket Manager
```go
type MockWSManager struct {
    OnlineUsers map[int64]bool
    SentMessages []MockMessage
}

func (m *MockWSManager) IsOnline(userID int64) bool {
    return m.OnlineUsers[userID]
}

func (m *MockWSManager) SendToUser(userID int64, msg []byte) {
    m.SentMessages = append(m.SentMessages, MockMessage{
        UserID: userID,
        Data:   msg,
    })
}
```

### Mock Notification Sender
```go
type MockNotificationSender struct {
    NotificationsSent []MockNotification
}

func (m *MockNotificationSender) SendNotification(userID int64, data map[string]interface{}) {
    m.NotificationsSent = append(m.NotificationsSent, MockNotification{
        UserID: userID,
        Data:   data,
    })
}
```

## Continuous Integration

### GitHub Actions Example
```yaml
name: WebSocket Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.19
    
    - name: Run WebSocket Tests
      run: |
        cd backend
        go test ./pkg/ws/... -v -cover
    
    - name: Upload Coverage
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
```

This testing guide ensures comprehensive coverage of the WebSocket implementation and provides patterns for testing integrations with other features.
