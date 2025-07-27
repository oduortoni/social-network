# Follow Response Notification

## Overview

Follow response notifications are sent to users when their follow requests are accepted or declined by other users. These are informational notifications that provide feedback on the outcome of follow requests and allow users to engage with new followers.

## Notification Data Structure

### WebSocket Notification Format (Accepted)
```json
{
  "type": "notification",
  "subtype": "follow_accepted",
  "user_id": 456,
  "nickname": "jane_doe",
  "avatar": "/uploads/avatars/456.jpg",
  "message": "Jane Doe accepted your follow request",
  "timestamp": 1640995200,
  "requires_action": false,
  "actions": ["view_profile", "send_message"],
  "additional_data": {
    "responder_id": 456,
    "follower_relationship_id": 123,
    "response_type": "accepted",
    "can_message": true,
    "mutual_followers": 5
  }
}
```

### WebSocket Notification Format (Declined)
```json
{
  "type": "notification",
  "subtype": "follow_declined",
  "user_id": 456,
  "nickname": "jane_doe",
  "avatar": "/uploads/avatars/456.jpg",
  "message": "Jane Doe declined your follow request",
  "timestamp": 1640995200,
  "requires_action": false,
  "actions": ["view_profile"],
  "additional_data": {
    "responder_id": 456,
    "response_type": "declined",
    "can_message": false
  }
}
```

### Data Fields Explanation
- **user_id**: ID of the user who responded to the follow request
- **nickname**: Display name of the responder
- **avatar**: Profile picture path of the responder
- **message**: Human-readable notification text
- **timestamp**: Unix timestamp when response was given
- **requires_action**: Always `false` - informational notifications
- **actions**: Available actions for the original requester
- **responder_id**: Same as user_id (for clarity)
- **follower_relationship_id**: ID from Followers table (for accepted requests)
- **response_type**: "accepted" or "declined"
- **can_message**: Whether the user can send messages to the responder
- **mutual_followers**: Number of mutual followers (for accepted requests)

## Required API Endpoints

Follow response notifications are triggered by existing follow request endpoints and don't require new API endpoints. They are generated when the follow request endpoints are used:

### Existing Endpoints That Trigger Follow Responses

#### 1. Accept Follow Request (Triggers follow_accepted notification)
```http
POST /follow-requests/{notificationId}/accept
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Follow request accepted",
  "new_follower": {
    "user_id": 123,
    "nickname": "john_doe",
    "avatar": "/uploads/avatars/123.jpg"
  }
}
```

**Side Effect:** Sends `follow_accepted` notification to the original requester (user_id: 123)

#### 2. Decline Follow Request (Triggers follow_declined notification)
```http
POST /follow-requests/{notificationId}/decline
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Follow request declined"
}
```

**Side Effect:** Sends `follow_declined` notification to the original requester

### Supporting Endpoints for Actions

#### 3. Get User Profile (For "View Profile" action)
```http
GET /users/{userId}
Authorization: Session cookie required
```

#### 4. Send Message (For "Send Message" action - accepted requests only)
```http
POST /chats/{userId}/messages
Content-Type: application/json
Authorization: Session cookie required

{
  "content": "Thanks for accepting my follow request!",
  "message_type": "text"
}
```

## Call-to-Action Implementation

### Activity Sidebar Actions

#### For Accepted Follow Requests:
1. **View Profile** - Opens the responder's profile in ContentRenderer
2. **Send Message** - Opens chat with the new follower

#### For Declined Follow Requests:
1. **View Profile** - Opens the responder's profile in ContentRenderer

### Frontend Action Handlers
```javascript
// In FollowResponseActivity component
const handleFollowResponseAction = async (action, userId, responseType) => {
  switch (action) {
    case 'view_profile':
      // Navigate to responder's profile using ContentRenderer
      onContentChange('user_profile', { userId });
      
      // Mark notification as read when viewing profile
      markNotificationAsRead();
      break;
      
    case 'send_message':
      if (responseType === 'accepted') {
        // Navigate to chat with the new follower
        onContentChange('chat', { userId, focusInput: true });
        
        // Mark notification as read when opening chat
        markNotificationAsRead();
      }
      break;
  }
};

// Helper function to mark notification as read
const markNotificationAsRead = async () => {
  try {
    // Remove notification from activity feed (frontend state)
    removeNotification(notificationId);
    
    // Optionally call API to mark as read if using persistent notifications
    await fetch(`/notifications/${notificationId}/read`, {
      method: 'POST',
      credentials: 'include'
    });
  } catch (error) {
    console.error('Failed to mark notification as read:', error);
  }
};
```

## Testing Approach

### Backend Testing

#### 1. Follow Accept Response Test
```bash
# 1. User A sends follow request to User B
curl -X POST http://localhost:8080/users/456/follow \
  -H "Cookie: session_id=user_a_session"

# 2. User B accepts the request (should trigger follow_accepted notification to User A)
curl -X POST http://localhost:8080/follow-requests/{requestId}/accept \
  -H "Cookie: session_id=user_b_session"

# Expected: Followers table updated with status 'accepted'
# Expected: WebSocket notification sent to User A with follow_accepted
```

#### 2. Follow Decline Response Test
```bash
# 1. User A sends follow request to User B
curl -X POST http://localhost:8080/users/456/follow \
  -H "Cookie: session_id=user_a_session"

# 2. User B declines the request (should trigger follow_declined notification to User A)
curl -X POST http://localhost:8080/follow-requests/{requestId}/decline \
  -H "Cookie: session_id=user_b_session"

# Expected: Followers table updated with status 'rejected'
# Expected: WebSocket notification sent to User A with follow_declined
```

#### 3. Notification Data Validation
```bash
# Verify notification contains correct data structure
# Check WebSocket message includes:
# - Correct responder information
# - Appropriate actions based on response type
# - Mutual follower count (for accepted requests)
```

### Frontend Testing

#### 1. Notification Display Test
- User A sends follow request to User B
- User B accepts the request
- Verify `follow_accepted` notification appears in User A's Activity Sidebar
- Check notification contains correct responder data and available actions

#### 2. Action Button Test (Accepted)
- Click "View Profile" button
- Verify ContentRenderer loads UserProfile component for responder
- Click "Send Message" button
- Verify ContentRenderer loads ChatWindow component with responder

#### 3. Action Button Test (Declined)
- User B declines follow request
- Verify `follow_declined` notification appears in User A's Activity Sidebar
- Check only "View Profile" action is available (no "Send Message")
- Click "View Profile" and verify correct profile loads

### Integration Testing

#### 1. End-to-End Accept Flow
1. User A sends follow request to User B (private profile)
2. User B receives follow request notification
3. User B accepts request through Activity Sidebar
4. User A receives follow_accepted notification
5. User A can now see User B's posts and send messages

#### 2. End-to-End Decline Flow
1. User A sends follow request to User B (private profile)
2. User B receives follow request notification
3. User B declines request through Activity Sidebar
4. User A receives follow_declined notification
5. User A cannot see User B's posts or send messages

#### 3. Multiple Request Scenarios
- Test rapid accept/decline of multiple requests
- Verify each generates appropriate response notification
- Check notification ordering and timestamps

## Implementation Files

### Backend Files to Create/Modify
```
backend/internal/
├── handlers/
│   └── follow_handler.go          # Add response notification logic
├── notifications/
│   └── follow_notifications.go    # Add follow response notifications
└── websocket/
    └── ws.go                      # Add follow response broadcasting
```

### Frontend Files to Create/Modify
```
frontend/components/
├── activity/
│   └── FollowResponseActivity.jsx # Follow response notification component
└── layout/
    └── ActivitySidebar.jsx        # Add follow response handling
```

## Database Schema Requirements

Follow response notifications use the existing `Followers` table and are triggered by status changes in that table. No additional database tables are needed.

### Existing Followers Table Structure
```sql
-- From 000002_create_followers_table.up.sql
CREATE TABLE Followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followee_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME,
    FOREIGN KEY (follower_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES Users(id) ON DELETE CASCADE,
    UNIQUE (follower_id, followee_id)
);
```

### Follow Response Notification Trigger Logic
```sql
-- When follow request is accepted
UPDATE Followers 
SET status = 'accepted', accepted_at = CURRENT_TIMESTAMP 
WHERE id = ? AND followee_id = ?;

-- Trigger: Create activity notification for original requester
INSERT INTO activity_notifications (user_id, type, from_user_id, reference_id, data)
VALUES (
    (SELECT follower_id FROM Followers WHERE id = ?), -- Original requester
    'follow_accepted',
    ?, -- User who accepted (followee_id)
    ?, -- Followers table ID
    JSON_OBJECT(
        'responder_nickname', ?,
        'responder_avatar', ?,
        'mutual_followers', (SELECT COUNT(*) FROM mutual_followers_view WHERE user1 = ? AND user2 = ?)
    )
);

-- When follow request is declined
UPDATE Followers 
SET status = 'rejected' 
WHERE id = ? AND followee_id = ?;

-- Trigger: Create activity notification for original requester
INSERT INTO activity_notifications (user_id, type, from_user_id, reference_id, data)
VALUES (
    (SELECT follower_id FROM Followers WHERE id = ?), -- Original requester
    'follow_declined',
    ?, -- User who declined (followee_id)
    ?, -- Followers table ID
    JSON_OBJECT(
        'responder_nickname', ?,
        'responder_avatar', ?
    )
);
```

### Follow Response Lifecycle
1. User A sends follow request to User B
2. Follow request stored in `Followers` table with status 'pending'
3. User B accepts/declines through Activity Sidebar
4. `Followers` table status updated to 'accepted'/'rejected'
5. Follow response notification created for User A
6. WebSocket notification sent to User A
7. User A can view profile or send message (if accepted)

### Query Patterns for Follow Responses
```sql
-- Get mutual followers count for accepted follow notification
SELECT COUNT(*) as mutual_count
FROM Followers f1
JOIN Followers f2 ON f1.followee_id = f2.followee_id
WHERE f1.follower_id = ? AND f2.follower_id = ? 
AND f1.status = 'accepted' AND f2.status = 'accepted';

-- Check if user can message another user (following relationship exists)
SELECT COUNT(*) > 0 as can_message
FROM Followers 
WHERE follower_id = ? AND followee_id = ? AND status = 'accepted';

-- Get follow relationship details for notification
SELECT f.*, u.nickname, u.avatar
FROM Followers f
JOIN Users u ON f.followee_id = u.id
WHERE f.id = ?;
```

### Notification Behavior
- **Follow Accepted**: Enables messaging, shows mutual followers, provides engagement actions
- **Follow Declined**: Informational only, limited to profile viewing
- **Auto-cleanup**: Response notifications can be cleaned up after user interaction
- **Privacy Respect**: Declined notifications don't reveal private profile information

This follow response system completes the follow request lifecycle and provides users with clear feedback on their social interactions.
