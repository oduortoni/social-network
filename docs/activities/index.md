# Activity Notifications Documentation

## Overview

This documentation covers the required notification types for the social network, their data structures, API endpoints, and testing approaches. Each notification type includes specific call-to-action requirements for the Activity Sidebar integration.

## Required Notification Types

### Core Social Notifications
- [Follow Request](./follow-request.md) - Private profile follow requests with Accept/Decline actions
- [Follow Response](./follow-response.md) - Follow request accepted/declined notifications

### Group Management Notifications
- [Group Invitation](./group-invitation.md) - Group invitations with Accept/Decline actions
- [Group Join Request](./group-join-request.md) - Join requests for group creators to Approve/Decline
- [Group Event Created](./group-event.md) - New group events with Going/Not Going actions

### Communication Notifications
- [New Message](./new-message.md) - Private and group messages with View/Reply actions

## Implementation Priority

1. Follow Request - Foundation for private profile system
2. New Message - Critical for real-time communication
3. Group Invitation - Essential for group functionality
4. Group Join Request - Required for group management
5. Group Event Created - Needed for event features
6. Follow Response - User experience completion

## Database Schema Analysis

After examining the existing migrations in `backend/pkg/db/migrations/sqlite/`, the database already has tables that handle most notification scenarios. Here's how notifications map to existing tables:

### Existing Tables for Notifications

**Followers Table (000002):**
```sql
CREATE TABLE Followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followee_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME
);
```

**Group_Members Table (000007):**
```sql
CREATE TABLE Group_Members (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT,
    is_accepted BOOLEAN DEFAULT 0,
    invited_by INTEGER,
    requested BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Messages Table (000011):**
```sql
CREATE TABLE Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    content TEXT,
    is_emoji BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Notifications Table (000012):**
```sql
CREATE TABLE Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT,
    message TEXT,
    is_read BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Notification Type Mapping to Existing Tables

**Follow Request Notifications:**
- **Source**: `Followers` table with `status = 'pending'`
- **Query**: `SELECT * FROM Followers WHERE followee_id = ? AND status = 'pending'`
- **Actions**: Update status to 'accepted' or 'rejected'

**Group Invitation Notifications:**
- **Source**: `Group_Members` table with `is_accepted = 0` AND `invited_by IS NOT NULL`
- **Query**: `SELECT * FROM Group_Members WHERE user_id = ? AND is_accepted = 0 AND invited_by IS NOT NULL`
- **Actions**: Update `is_accepted = 1` or delete row

**Group Join Request Notifications:**
- **Source**: `Group_Members` table with `requested = 1` AND `is_accepted = 0`
- **Query**: `SELECT * FROM Group_Members WHERE group_id = ? AND requested = 1 AND is_accepted = 0`
- **Actions**: Update `is_accepted = 1` or delete row

**New Message Notifications:**
- **Source**: `Messages` table (real-time WebSocket notifications)
- **Query**: `SELECT * FROM Messages WHERE receiver_id = ? ORDER BY created_at DESC`
- **Actions**: Mark as read (frontend state, not database)

**Event Response Notifications:**
- **Source**: `Event_Responses` table or Group_Events creation
- **Query**: Event-specific queries for group members
- **Actions**: Insert into Event_Responses table

### Enhanced Notifications Table Needed

The existing `Notifications` table is too basic for rich notifications. We need to enhance it or create a new table for activity sidebar notifications:

```sql
-- Enhanced notifications table for Activity Sidebar
CREATE TABLE activity_notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL, -- 'follow_request', 'group_invitation', etc.
    from_user_id INTEGER NOT NULL,
    reference_id INTEGER, -- follower_id, group_id, message_id, etc.
    data TEXT, -- JSON string with notification-specific data
    is_read BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (from_user_id) REFERENCES Users(id) ON DELETE CASCADE
);
```

### Notification Lifecycle with Existing Tables

1. **Follow Request**: Insert into `Followers` with status 'pending' + create `activity_notifications` entry
2. **Group Invitation**: Insert into `Group_Members` with `is_accepted=0, invited_by=X` + create `activity_notifications` entry
3. **User Action**: Update source table (Followers, Group_Members) + mark `activity_notifications` as read
4. **Cleanup**: Remove read notifications older than X days from `activity_notifications`

### Benefits of This Approach

- **Leverages existing tables** for data integrity and relationships
- **Activity notifications table** provides rich WebSocket notification data
- **No data duplication** - source tables remain authoritative
- **Clean separation** between business logic (Followers, Group_Members) and notification display

## Notification Data Structure

### Base Notification Format
All notifications follow this standard structure:

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
    "notification_id": "uuid-123", -- ID from temporary_notifications table
    "target_user_id": 456
  }
}
```

### Database to WebSocket Mapping
```sql
-- Database row in temporary_notifications
{
  "id": "uuid-123",
  "type": "follow_request",
  "from_user_id": 123,
  "to_user_id": 456,
  "status": "pending",
  "data": {"requester_nickname": "john_doe", "requester_avatar": "/uploads/avatars/123.jpg"}
}

-- Becomes WebSocket notification
{
  "type": "notification",
  "subtype": "follow_request",
  "user_id": 123,
  "nickname": "john_doe", -- from data.requester_nickname
  "avatar": "/uploads/avatars/123.jpg", -- from data.requester_avatar
  "additional_data": {"notification_id": "uuid-123"}
}
```

### Action Categories

#### High Priority (Requires User Decision)
- Follow Request → Accept, Decline, View Profile
- Group Invitation → Accept, Decline, View Group
- Group Join Request → Approve, Decline, View Profile

#### Medium Priority (Engagement Actions)
- Group Event → View Event, Going, Not Going
- New Message → View Message, Reply

#### Low Priority (Informational)
- Follow Response → View Profile, Send Message

## Activity Sidebar Integration

### Call-to-Action Mapping
Each notification type maps to specific ContentRenderer navigation:

```javascript
// Activity click handlers by notification type
const NOTIFICATION_ACTIONS = {
  'follow_request': {
    'accept': () => handleFollowAction('accept'),
    'decline': () => handleFollowAction('decline'),
    'view_profile': (userId) => onContentChange('user_profile', { userId })
  },
  'group_invitation': {
    'accept': () => handleGroupInvite('accept'),
    'decline': () => handleGroupInvite('decline'),
    'view_group': (groupId) => onContentChange('group_detail', { groupId })
  },
  'new_message': {
    'view_message': (userId) => onContentChange('chat', { userId }),
    'reply': (userId) => onContentChange('chat', { userId })
  }
};
```

## Testing Strategy

### Backend Testing
Each notification type requires testing:
1. Trigger Conditions - When notifications are created
2. Data Validation - Correct data structure and content
3. WebSocket Delivery - Real-time notification sending
4. Action Endpoints - API endpoints for user actions

### Frontend Testing
1. Notification Display - Correct rendering in Activity Sidebar
2. Action Handling - Call-to-action button functionality
3. Content Navigation - ContentRenderer integration
4. Real-time Updates - WebSocket message handling

### Integration Testing
1. End-to-End Flows - Complete notification lifecycle
2. Cross-Component - Activity Sidebar → ContentRenderer → API
3. Real-time Sync - Multiple users receiving notifications
4. Error Scenarios - Network failures, invalid actions

## Database Design Strategy

### Temporary Notifications Table
All actionable notifications use a single `temporary_notifications` table (detailed schema above).

### Benefits of This Approach
- Single source of truth for all temporary notifications
- Automatic cleanup of processed/expired notifications
- Consistent data structure across notification types
- Efficient querying with type and status indexes
- Flexible data storage with JSON column for type-specific fields

### Permanent vs Temporary Data
- Temporary: Follow requests, group invitations, event responses
- Permanent: Followers relationships, chat messages, user profiles

## File Structure for Implementation

### Backend Files
```
backend/internal/
├── notifications/
│   ├── temporary_notification_service.go  # Handles temp notifications
│   ├── follow_handler.go
│   ├── group_handler.go
│   ├── message_handler.go
│   └── cleanup_service.go                 # Cleanup expired notifications
├── websocket/
│   └── notification_broadcaster.go
└── api/
    └── notification_routes.go
```

### Frontend Files
```
frontend/components/
├── activity/
│   ├── FollowRequestActivity.jsx
│   ├── GroupInviteActivity.jsx
│   ├── MessageActivity.jsx
│   └── ActivityActionHandler.js
└── layout/
    └── ActivitySidebar.jsx
```

## WebSocket Integration

### Real-time Delivery
All notifications are delivered via WebSocket with this pattern:

```json
{
  "type": "activity",
  "notification": {
    // Full notification object
  }
}
```

### Client Handling
```javascript
wsService.onMessage('activity', (data) => {
  const notification = data.notification;
  addToActivityFeed(notification);
  updateNotificationBadge();

  if (notification.requires_action) {
    showHighPriorityIndicator();
  }
});
```

## Additional API Endpoints Needed

Based on the existing endpoint.md structure, these additional endpoints are needed for notifications:

### Notification Management
- `GET /notifications` - Already exists
- `POST /notifications/{notificationId}/read` - Already exists

### Missing Endpoints for Activity System
- `POST /events/{eventId}/respond` - Already exists (for group events)
- `GET /group-invitations` - Already exists
- `POST /group-invitations/{invitationId}/accept` - Already exists
- `POST /group-invitations/{invitationId}/decline` - Already exists
- `GET /groups/{groupId}/join-requests` - Already exists
- `POST /join-requests/{requestId}/accept` - Already exists
- `POST /join-requests/{requestId}/decline` - Already exists

### WebSocket Endpoint
- `ws://your-api-domain.com/chat` - Already exists

All required endpoints for our notification system already exist in endpoint.md! Our notification documentation now aligns perfectly with the established API structure.

---

Navigate to specific notification documentation for detailed implementation guidance. Each section provides complete data structures, API specifications, and testing approaches.