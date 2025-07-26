# WebSocket Integration Guide

## For Developers Implementing Other Features

This guide shows how to leverage the existing WebSocket implementation to add real-time functionality to new features like posts, events, friend requests, and more.

## Quick Start Integration

### 1. Basic Setup

The WebSocket system is already integrated into the main router. To use it in your feature:

```go
// In your handler file
import "github.com/tajjjjr/social-network/backend/pkg/ws"

// Get the WebSocket manager from your router setup
func NewYourHandler(db *sql.DB, wsManager *ws.Manager) *YourHandler {
    return &YourHandler{
        DB:        db,
        WSManager: wsManager,
    }
}
```

### 2. Sending Real-time Notifications

```go
// Create notification sender
notifier := ws.NewDBNotificationSender(wsManager)

// Send notification to a specific user
notifier.SendNotification(userID, map[string]interface{}{
    "type":      "notification",
    "subtype":   "friend_request",
    "message":   "You have a new friend request from John Doe",
    "user_id":   senderID,
    "timestamp": time.Now().Unix(),
})
```

### 3. Check User Online Status

```go
// Check if user is online before sending real-time notification
if wsManager.IsOnline(userID) {
    // Send real-time notification
    notifier.SendNotification(userID, notificationData)
} else {
    // Store notification in database for later retrieval
    storeNotificationInDB(userID, notificationData)
}
```

## Feature-Specific Integration Examples

### Friend System Integration

When implementing friend requests, follow this pattern:

```go
func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
    // ... authentication and validation ...
    
    // 1. Store friend request in database
    _, err := h.DB.Exec(`
        INSERT INTO Friend_Requests (requester_id, recipient_id, status) 
        VALUES (?, ?, 'pending')
    `, requesterID, recipientID)
    
    if err != nil {
        http.Error(w, "Failed to send friend request", http.StatusInternalServerError)
        return
    }
    
    // 2. Send real-time notification if user is online
    if h.WSManager.IsOnline(recipientID) {
        h.Notifier.SendNotification(recipientID, map[string]interface{}{
            "type":         "notification",
            "subtype":      "friend_request",
            "message":      fmt.Sprintf("%s sent you a friend request", requesterName),
            "requester_id": requesterID,
            "timestamp":    time.Now().Unix(),
        })
    }
    
    // 3. Store notification in database for offline users
    _, err = h.DB.Exec(`
        INSERT INTO Notifications (user_id, type, message) 
        VALUES (?, 'friend_request', ?)
    `, recipientID, fmt.Sprintf("%s sent you a friend request", requesterName))
    
    w.WriteHeader(http.StatusCreated)
}
```

### Posts System Integration

For real-time post updates:

```go
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    // ... create post logic ...
    
    // Get user's followers for real-time updates
    followers, err := h.getFollowers(userID)
    if err != nil {
        log.Printf("Failed to get followers: %v", err)
        return
    }
    
    // Send real-time notifications to online followers
    for _, followerID := range followers {
        if h.WSManager.IsOnline(followerID) {
            h.Notifier.SendNotification(followerID, map[string]interface{}{
                "type":      "notification",
                "subtype":   "new_post",
                "message":   fmt.Sprintf("%s shared a new post", userName),
                "post_id":   postID,
                "user_id":   userID,
                "timestamp": time.Now().Unix(),
            })
        }
    }
}
```

### Group Events Integration

For event notifications:

```go
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    // ... create event logic ...
    
    // Get group members
    memberIDs, err := h.GroupFetcher.GetGroupMemberIDs(groupID)
    if err != nil {
        log.Printf("Failed to get group members: %v", err)
        return
    }
    
    // Notify all group members about new event
    for _, memberID := range memberIDs {
        if memberID == creatorID {
            continue // Don't notify the creator
        }
        
        if h.WSManager.IsOnline(memberID) {
            h.Notifier.SendNotification(memberID, map[string]interface{}{
                "type":      "notification",
                "subtype":   "new_event",
                "message":   fmt.Sprintf("New event '%s' in %s", eventTitle, groupName),
                "event_id":  eventID,
                "group_id":  groupID,
                "timestamp": time.Now().Unix(),
            })
        }
    }
}
```

## Custom Message Types

You can extend the WebSocket system with custom message types:

### 1. Define New Message Type

```go
// In your feature package
type CustomMessage struct {
    Type      string                 `json:"type"`
    Subtype   string                 `json:"subtype"`
    Data      map[string]interface{} `json:"data"`
    Timestamp int64                  `json:"timestamp"`
}
```

### 2. Send Custom Messages

```go
func (h *YourHandler) sendCustomNotification(userID int64, data map[string]interface{}) {
    if h.WSManager.IsOnline(userID) {
        customMsg := map[string]interface{}{
            "type":      "custom_feature",
            "subtype":   "specific_action",
            "data":      data,
            "timestamp": time.Now().Unix(),
        }
        
        h.Notifier.SendNotification(userID, customMsg)
    }
}
```

## Best Practices

### 1. Always Check Online Status

```go
// Good: Check before sending real-time notifications
if wsManager.IsOnline(userID) {
    sendRealTimeNotification(userID, data)
} else {
    storeForLaterRetrieval(userID, data)
}

// Bad: Sending without checking
sendRealTimeNotification(userID, data) // May fail silently
```

### 2. Graceful Error Handling

```go
func (h *Handler) sendNotificationSafely(userID int64, data map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Notification sending failed: %v", r)
        }
    }()
    
    if h.WSManager.IsOnline(userID) {
        h.Notifier.SendNotification(userID, data)
    }
}
```

### 3. Batch Notifications

```go
// For multiple users, batch the operations
func (h *Handler) notifyMultipleUsers(userIDs []int64, data map[string]interface{}) {
    onlineUsers := make([]int64, 0)
    
    // First, check who's online
    for _, userID := range userIDs {
        if h.WSManager.IsOnline(userID) {
            onlineUsers = append(onlineUsers, userID)
        }
    }
    
    // Then send notifications
    for _, userID := range onlineUsers {
        h.Notifier.SendNotification(userID, data)
    }
}
```

### 4. Notification Types Convention

Use consistent notification types across features:

```go
const (
    NotificationFriendRequest = "friend_request"
    NotificationFriendAccept  = "friend_accept"
    NotificationNewPost       = "new_post"
    NotificationPostLike      = "post_like"
    NotificationPostComment   = "post_comment"
    NotificationGroupInvite   = "group_invite"
    NotificationNewEvent      = "new_event"
    NotificationEventRSVP     = "event_rsvp"
)
```

## Testing Your Integration

### 1. Unit Testing

```go
func TestYourFeatureNotification(t *testing.T) {
    // Create mock WebSocket manager
    mockManager := &MockWSManager{}
    mockNotifier := &MockNotifier{}
    
    handler := NewYourHandler(db, mockManager, mockNotifier)
    
    // Test your feature
    handler.YourMethod(w, r)
    
    // Verify notification was sent
    assert.True(t, mockNotifier.NotificationSent)
}
```

### 2. Integration Testing

```go
func TestRealTimeNotification(t *testing.T) {
    // Use the existing WebSocket test setup
    manager := ws.NewManager(
        ws.NewDBSessionResolver(db),
        ws.NewDBGroupMemberFetcher(db),
        ws.NewDBMessagePersister(db),
    )
    
    // Connect test client
    client := connectTestClient(t, manager)
    
    // Trigger your feature
    triggerYourFeature()
    
    // Verify real-time notification received
    notification := readNotificationFromClient(t, client)
    assert.Equal(t, "your_notification_type", notification.Type)
}
```

## Common Patterns

### 1. Database + Real-time Pattern

```go
// 1. Update database
err := updateDatabase(data)
if err != nil {
    return err
}

// 2. Send real-time notification
if wsManager.IsOnline(userID) {
    sendNotification(userID, data)
}

// 3. Store notification for offline users
storeNotification(userID, data)
```

### 2. Conditional Notification Pattern

```go
// Only notify if user has enabled this notification type
if userHasNotificationEnabled(userID, notificationType) {
    if wsManager.IsOnline(userID) {
        sendNotification(userID, data)
    }
}
```

This integration guide provides everything you need to add real-time functionality to your features using the existing WebSocket implementation.
