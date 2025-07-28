# Follow Request Notification

## Overview

Follow request notifications are sent when a user wants to follow another user with a private profile. The recipient can accept or decline the request through the Activity Sidebar.

## Notification Data Structure

### WebSocket Notification Format
```json
{
  "type": "notification",
  "subtype": "follow_request",
  "user_id": 123,
  "nickname": "john_doe",
  "avatar": "/uploads/avatars/123.jpg",
  "message": "John Doe sent you a follow request",
  "timestamp": 1640995200,
  "requires_action": true,
  "actions": ["accept", "decline", "view_profile"],
  "additional_data": {
    "notification_id": "uuid-abc123", -- ID from temporary_notifications table
    "requester_id": 123,
    "target_user_id": 456
  }
}
```

### Data Fields Explanation
- **user_id**: ID of the user who sent the follow request
- **nickname**: Display name of the requester
- **avatar**: Profile picture path of the requester
- **message**: Human-readable notification text
- **timestamp**: Unix timestamp when request was sent
- **requires_action**: Always `true` for follow requests
- **actions**: Available actions for the recipient
- **notification_id**: Unique identifier from temporary_notifications table
- **requester_id**: Same as user_id (for clarity)
- **target_user_id**: ID of the user receiving the request

## Required API Endpoints

### 1. Send Follow Request
```http
POST /users/{userId}/follow
Content-Type: application/json
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Follow request sent",
  "notification_id": "uuid-abc123"
}
```

### 2. Accept Follow Request
```http
POST /follow-requests/{notificationId}/accept
Content-Type: application/json
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

### 3. Decline Follow Request
```http
POST /follow-requests/{notificationId}/decline
Content-Type: application/json
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Follow request declined"
}
```

### 4. Get Pending Follow Requests
```http
GET /follow-requests
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "requests": [
    {
      "request_id": "uuid-abc123",
      "requester": {
        "user_id": 123,
        "nickname": "john_doe",
        "avatar": "/uploads/avatars/123.jpg"
      },
      "created_at": 1640995200
    }
  ]
}
```

## Call-to-Action Implementation

### Activity Sidebar Actions
The follow request notification displays three action buttons:

1. Accept - Accepts the follow request
2. Decline - Declines the follow request
3. View Profile - Opens the requester's profile in ContentRenderer

### Frontend Action Handlers
```javascript
// In FollowRequestActivity component
const handleFollowAction = async (action, requestId, userId) => {
  switch (action) {
    case 'accept':
      try {
        await fetch(`/follow-requests/${requestId}/accept`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        // Remove notification from activity feed
        removeNotification(requestId);
        // Show success message
        showToast('Follow request accepted');
      } catch (error) {
        showToast('Failed to accept request', 'error');
      }
      break;

    case 'decline':
      try {
        await fetch(`/follow-requests/${requestId}/decline`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        // Remove notification from activity feed
        removeNotification(requestId);
        // Show success message
        showToast('Follow request declined');
      } catch (error) {
        showToast('Failed to decline request', 'error');
      }
      break;

    case 'view_profile':
      // Navigate to user profile using ContentRenderer
      onContentChange('user_profile', { userId });
      break;
  }
};
```

## Testing Approach

### Backend Testing

#### 1. Trigger Condition Test
```bash
# Test follow request creation
curl -X POST http://localhost:8080/users/456/follow \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=test_session"

# Expected: Follow request created in database
# Expected: WebSocket notification sent to target user
```

#### 2. Action Endpoint Tests
```bash
# Test accept action
curl -X POST http://localhost:8080/follow-requests/uuid-abc123/accept \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=target_user_session"

# Test decline action
curl -X POST http://localhost:8080/follow-requests/uuid-abc123/decline \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=target_user_session"

# Test get pending requests
curl -X GET http://localhost:8080/follow-requests \
  -H "Cookie: session_id=target_user_session"
```

### Frontend Testing

#### 1. Notification Display Test
- Send follow request from User A to User B
- Verify notification appears in User B's Activity Sidebar
- Check notification contains correct user data and actions

#### 2. Action Button Test
- Click "Accept" button
- Verify API call is made with correct request_id
- Check notification is removed from activity feed
- Confirm success message is displayed

#### 3. ContentRenderer Integration Test
- Click "View Profile" action
- Verify ContentRenderer loads UserProfile component
- Check correct user data is passed to component

### Integration Testing

#### 1. End-to-End Flow
1. User A sends follow request to User B (private profile)
2. User B receives real-time notification
3. User B clicks "Accept" in Activity Sidebar
4. User A receives follow acceptance notification
5. User A can now see User B's posts

#### 2. Error Scenarios
- Test with invalid request_id
- Test with expired session
- Test with already processed request
- Test WebSocket disconnection during action

## Implementation Files

### Backend Files to Create/Modify
```
backend/internal/
├── handlers/
│   └── follow_handler.go          # Follow request API endpoints
├── models/
│   └── follow_request.go          # Follow request data model
├── notifications/
│   └── follow_notifications.go    # Follow notification logic
└── websocket/
    └── ws.go                      # Add follow notification broadcasting
```

### Frontend Files to Create/Modify
```
frontend/components/
├── activity/
│   └── FollowRequestActivity.jsx  # Follow request notification component
├── layout/
│   └── ActivitySidebar.jsx        # Add follow request handling
└── profile/
    └── UserProfile.jsx            # Add follow request button
```

## Database Schema Requirements

Follow requests use the existing `Followers` table from migration 000002. This table already handles the complete follow request lifecycle with status tracking.

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

### Follow Request Data Flow
```sql
-- 1. User A (123) wants to follow User B (456) - private profile
INSERT INTO Followers (follower_id, followee_id, status)
VALUES (123, 456, 'pending');

-- 2. Create activity notification for User B
INSERT INTO activity_notifications (user_id, type, from_user_id, reference_id, data)
VALUES (456, 'follow_request', 123, LAST_INSERT_ROWID(),
        '{"requester_nickname": "john_doe", "requester_avatar": "/uploads/avatars/123.jpg"}');

-- 3. User B accepts the request
UPDATE Followers
SET status = 'accepted', accepted_at = CURRENT_TIMESTAMP
WHERE id = ? AND followee_id = 456;

-- 4. Mark notification as read
UPDATE activity_notifications
SET is_read = 1
WHERE user_id = 456 AND reference_id = ?;
```

### Follow Request Lifecycle
1. User A attempts to follow User B (private profile)
2. Row created in `Followers` table with status 'pending'
3. Activity notification created for rich WebSocket data
4. WebSocket notification sent to User B
5. User B accepts/declines through Activity Sidebar
6. `Followers` table status updated to 'accepted'/'rejected'
7. Activity notification marked as read

### Query Patterns for Follow Requests
```sql
-- Get pending follow requests for a user
SELECT f.*, u.nickname, u.avatar
FROM Followers f
JOIN Users u ON f.follower_id = u.id
WHERE f.followee_id = ? AND f.status = 'pending'
ORDER BY f.requested_at DESC;

-- Check if follow relationship exists
SELECT status FROM Followers
WHERE follower_id = ? AND followee_id = ?;

-- Get follower count
SELECT COUNT(*) FROM Followers
WHERE followee_id = ? AND status = 'accepted';
```

This follow request notification system provides the foundation for private profile functionality and demonstrates the pattern for all other notification types with clear data structures, API endpoints, and testing approaches.
