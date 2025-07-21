# WebSocket Extension Guide

## Extending the WebSocket System

This guide shows how to extend the existing WebSocket implementation to support new features and message types while maintaining backward compatibility.

## Adding New Message Types

### 1. Define New Message Structure

Create a new message type for your feature:

```go
// In your feature package (e.g., pkg/posts/)
type PostMessage struct {
    Type      string `json:"type"`
    PostID    int64  `json:"post_id"`
    Action    string `json:"action"` // "like", "comment", "share"
    Content   string `json:"content,omitempty"`
    UserID    int64  `json:"user_id"`
    Timestamp int64  `json:"timestamp"`
}
```

### 2. Extend Message Handling

Add your message type to the WebSocket message router:

```go
// In pkg/ws/ws.go, extend the readPump function
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
        case "post_action": // New message type
            m.handlePostAction(c.ID, msg, encoded)
        }
    }
}
```

### 3. Implement Custom Handler

Add a custom handler method to the Manager:

```go
// Add to pkg/ws/ws.go
func (m *Manager) handlePostAction(senderID int64, msg *Message, encoded []byte) {
    // Parse post-specific data
    var postMsg PostMessage
    if err := json.Unmarshal(encoded, &postMsg); err != nil {
        return
    }
    
    // Get post owner and followers
    postOwnerID, followers := m.getPostRecipientsFromDB(postMsg.PostID)
    
    // Send to post owner if it's not the sender
    if postOwnerID != senderID {
        m.SendToUser(postOwnerID, encoded)
    }
    
    // Send to followers who are online
    for _, followerID := range followers {
        if followerID != senderID && m.IsOnline(followerID) {
            m.SendToUser(followerID, encoded)
        }
    }
}
```

## Adding New Interfaces

### 1. Define Feature-Specific Interface

Create interfaces for your feature's data access:

```go
// In your feature package
type PostDataFetcher interface {
    GetPostOwner(postID int64) (int64, error)
    GetPostFollowers(postID int64) ([]int64, error)
    GetPostDetails(postID int64) (*Post, error)
}
```

### 2. Extend Manager Constructor

Add your interface to the Manager:

```go
// Extend the Manager struct
type Manager struct {
    clients      map[int64]*Client
    mu           sync.RWMutex
    resolver     SessionResolver
    groupQuery   GroupMemberFetcher
    persister    MessagePersister
    postFetcher  PostDataFetcher  // New interface
}

// Update constructor
func NewManager(resolver SessionResolver, groupFetcher GroupMemberFetcher, 
               persister MessagePersister, postFetcher PostDataFetcher) *Manager {
    return &Manager{
        clients:     make(map[int64]*Client),
        resolver:    resolver,
        groupQuery:  groupFetcher,
        persister:   persister,
        postFetcher: postFetcher,
    }
}
```

### 3. Implement the Interface

Create the implementation:

```go
// In your feature package
type DBPostDataFetcher struct {
    DB *sql.DB
}

func NewDBPostDataFetcher(db *sql.DB) *DBPostDataFetcher {
    return &DBPostDataFetcher{DB: db}
}

func (p *DBPostDataFetcher) GetPostOwner(postID int64) (int64, error) {
    var ownerID int64
    err := p.DB.QueryRow("SELECT user_id FROM Posts WHERE id = ?", postID).Scan(&ownerID)
    return ownerID, err
}

func (p *DBPostDataFetcher) GetPostFollowers(postID int64) ([]int64, error) {
    rows, err := p.DB.Query(`
        SELECT f.follower_id 
        FROM Posts p 
        JOIN Followers f ON p.user_id = f.followee_id 
        WHERE p.id = ? AND f.is_accepted = 1
    `, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var followers []int64
    for rows.Next() {
        var followerID int64
        if err := rows.Scan(&followerID); err == nil {
            followers = append(followers, followerID)
        }
    }
    return followers, nil
}
```

## Custom Notification Types

### 1. Define Notification Structure

```go
type CustomNotification struct {
    Type      string                 `json:"type"`
    Subtype   string                 `json:"subtype"`
    Message   string                 `json:"message"`
    Data      map[string]interface{} `json:"data"`
    Timestamp int64                  `json:"timestamp"`
}
```

### 2. Create Notification Builder

```go
type NotificationBuilder struct {
    wsManager *Manager
    db        *sql.DB
}

func NewNotificationBuilder(wsManager *Manager, db *sql.DB) *NotificationBuilder {
    return &NotificationBuilder{
        wsManager: wsManager,
        db:        db,
    }
}

func (nb *NotificationBuilder) SendPostLikeNotification(postID, likerID, postOwnerID int64) {
    // Get liker's name
    likerName := nb.getUserName(likerID)
    
    notification := map[string]interface{}{
        "type":      "notification",
        "subtype":   "post_like",
        "message":   fmt.Sprintf("%s liked your post", likerName),
        "post_id":   postID,
        "liker_id":  likerID,
        "timestamp": time.Now().Unix(),
    }
    
    // Send real-time if online
    if nb.wsManager.IsOnline(postOwnerID) {
        notifier := NewDBNotificationSender(nb.wsManager)
        notifier.SendNotification(postOwnerID, notification)
    }
    
    // Store in database
    nb.storeNotification(postOwnerID, notification)
}
```

## Adding Real-time Features to Existing Endpoints

### 1. Extend Existing Handlers

Add WebSocket integration to your existing HTTP handlers:

```go
type PostHandler struct {
    DB           *sql.DB
    WSManager    *ws.Manager
    Notifier     *ws.NotificationSender
    PostFetcher  PostDataFetcher
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
    // ... existing like logic ...
    
    // Add real-time notification
    if h.WSManager.IsOnline(postOwnerID) {
        h.Notifier.SendNotification(postOwnerID, map[string]interface{}{
            "type":      "notification",
            "subtype":   "post_like",
            "message":   fmt.Sprintf("%s liked your post", likerName),
            "post_id":   postID,
            "liker_id":  likerID,
            "timestamp": time.Now().Unix(),
        })
    }
    
    // Broadcast to followers
    followers, _ := h.PostFetcher.GetPostFollowers(postID)
    for _, followerID := range followers {
        if h.WSManager.IsOnline(followerID) {
            h.Notifier.SendNotification(followerID, map[string]interface{}{
                "type":      "post_update",
                "subtype":   "post_like",
                "post_id":   postID,
                "liker_id":  likerID,
                "timestamp": time.Now().Unix(),
            })
        }
    }
}
```

## Advanced Message Routing

### 1. Conditional Message Routing

Route messages based on user preferences or relationships:

```go
func (m *Manager) SendConditionalNotification(userID int64, notification map[string]interface{}) {
    // Check user notification preferences
    if !m.userWantsNotification(userID, notification["subtype"].(string)) {
        return
    }
    
    // Check if users are connected (friends/followers)
    senderID := notification["sender_id"].(int64)
    if !m.usersAreConnected(userID, senderID) {
        return
    }
    
    // Send notification
    if m.IsOnline(userID) {
        notifier := NewDBNotificationSender(m)
        notifier.SendNotification(userID, notification)
    }
}
```

### 2. Batch Message Processing

Handle multiple recipients efficiently:

```go
func (m *Manager) BroadcastToFollowers(userID int64, message map[string]interface{}) {
    followers := m.getFollowers(userID)
    
    // Batch process online users
    onlineFollowers := make([]int64, 0)
    for _, followerID := range followers {
        if m.IsOnline(followerID) {
            onlineFollowers = append(onlineFollowers, followerID)
        }
    }
    
    // Send to all online followers
    notifier := NewDBNotificationSender(m)
    for _, followerID := range onlineFollowers {
        notifier.SendNotification(followerID, message)
    }
}
```

## Performance Optimizations

### 1. Message Queuing

Implement message queuing for high-volume scenarios:

```go
type MessageQueue struct {
    queue   chan QueuedMessage
    workers int
    manager *Manager
}

type QueuedMessage struct {
    UserID  int64
    Message map[string]interface{}
}

func NewMessageQueue(manager *Manager, workers int) *MessageQueue {
    mq := &MessageQueue{
        queue:   make(chan QueuedMessage, 1000),
        workers: workers,
        manager: manager,
    }
    
    // Start worker goroutines
    for i := 0; i < workers; i++ {
        go mq.worker()
    }
    
    return mq
}

func (mq *MessageQueue) worker() {
    notifier := NewDBNotificationSender(mq.manager)
    for msg := range mq.queue {
        if mq.manager.IsOnline(msg.UserID) {
            notifier.SendNotification(msg.UserID, msg.Message)
        }
    }
}

func (mq *MessageQueue) QueueNotification(userID int64, message map[string]interface{}) {
    select {
    case mq.queue <- QueuedMessage{UserID: userID, Message: message}:
        // Queued successfully
    default:
        // Queue is full, handle overflow
        log.Printf("Message queue full, dropping notification for user %d", userID)
    }
}
```

### 2. Connection Pooling

Optimize database connections for high-frequency operations:

```go
type OptimizedManager struct {
    *Manager
    dbPool *sql.DB
    cache  *sync.Map // For caching frequently accessed data
}

func (om *OptimizedManager) GetCachedUserData(userID int64) (*UserData, error) {
    // Check cache first
    if cached, ok := om.cache.Load(userID); ok {
        return cached.(*UserData), nil
    }
    
    // Fetch from database
    userData, err := om.fetchUserDataFromDB(userID)
    if err != nil {
        return nil, err
    }
    
    // Cache for future use
    om.cache.Store(userID, userData)
    return userData, nil
}
```

## Testing Extensions

### 1. Test New Message Types

```go
func TestCustomMessageType(t *testing.T) {
    manager := setupTestManager(t)
    
    // Test custom message handling
    customMsg := map[string]interface{}{
        "type":    "post_action",
        "post_id": 123,
        "action":  "like",
        "user_id": 456,
    }
    
    // Send message and verify routing
    // ... test implementation ...
}
```

### 2. Test Performance

```go
func BenchmarkMessageBroadcast(b *testing.B) {
    manager := setupTestManager(b)
    
    // Create many mock clients
    for i := 0; i < 1000; i++ {
        client := NewClient(int64(i), nil)
        manager.Register(client)
    }
    
    message := []byte(`{"type":"broadcast","content":"test"}`)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        manager.BroadcastToAll(message)
    }
}
```

## Migration Guide

When extending the WebSocket system:

1. **Backward Compatibility**: Ensure existing message types continue to work
2. **Database Migrations**: Add new tables/columns as needed
3. **Gradual Rollout**: Deploy extensions incrementally
4. **Monitoring**: Add metrics for new message types
5. **Documentation**: Update API documentation

This extension guide provides the foundation for adding any new real-time features to the social network while maintaining the existing WebSocket infrastructure.
