# Group Join Request Notification

## Overview

Group join request notifications are sent to group admins/moderators when a user requests to join their group. The admin can approve or decline the request through the Activity Sidebar. Upon approval, the user becomes a member of the group.

## Notification Data Structure

### WebSocket Notification Format
```json
{
  "type": "notification",
  "subtype": "group_join_request",
  "user_id": 456,
  "nickname": "bob_wilson",
  "avatar": "/uploads/avatars/456.jpg",
  "message": "Bob Wilson wants to join Photography Club",
  "timestamp": 1640995200,
  "requires_action": true,
  "actions": ["approve", "decline", "view_profile"],
  "additional_data": {
    "requester_id": 456,
    "group_id": 789,
    "group_name": "Photography Club",
    "group_avatar": "/uploads/groups/789.jpg",
    "request_message": "I love photography and would like to join your community"
  }
}
```

### Data Fields Explanation
- **user_id**: ID of the user requesting to join (requester)
- **nickname**: Display name of the requester
- **avatar**: Profile picture path of the requester
- **message**: Human-readable notification text
- **timestamp**: Unix timestamp when request was sent
- **requires_action**: Always `true` for join requests
- **actions**: Available actions for the group admin
- **requester_id**: Same as user_id (for clarity)
- **group_id**: ID of the group being requested to join
- **group_name**: Name of the group
- **group_avatar**: Group's profile picture
- **request_message**: Optional message from the requester

## Required API Endpoints

### 1. Send Group Join Request
```http
POST /groups/{groupId}/join
Content-Type: application/json
Authorization: Session cookie required

{
  "message": "I love photography and would like to join your community"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Join request sent to group admins",
  "group": {
    "group_id": 789,
    "group_name": "Photography Club"
  }
}
```

### 2. Approve Group Join Request
```http
POST /join-requests/{requestId}/accept
Content-Type: application/json
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Join request approved",
  "new_member": {
    "user_id": 456,
    "nickname": "bob_wilson",
    "avatar": "/uploads/avatars/456.jpg"
  }
}
```

### 3. Decline Group Join Request
```http
POST /join-requests/{requestId}/decline
Content-Type: application/json
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Join request declined"
}
```

### 4. Get Pending Join Requests (For Group Admins)
```http
GET /groups/{groupId}/join-requests
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "requests": [
    {
      "requester": {
        "user_id": 456,
        "nickname": "bob_wilson",
        "avatar": "/uploads/avatars/456.jpg"
      },
      "group": {
        "group_id": 789,
        "group_name": "Photography Club"
      },
      "request_message": "I love photography and would like to join",
      "created_at": 1640995200
    }
  ]
}
```

## Call-to-Action Implementation

### Activity Sidebar Actions
The group join request notification displays three action buttons:

1. Approve - Approves the join request and adds user to group
2. Decline - Declines the join request
3. View Profile - Opens the requester's profile in ContentRenderer

### Frontend Action Handlers
```javascript
// In GroupJoinRequestActivity component
const handleJoinRequestAction = async (action, requestId, userId, groupId) => {
  switch (action) {
    case 'approve':
      try {
        const response = await fetch(`/join-requests/${requestId}/accept`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        
        const result = await response.json();
        
        // Remove notification from activity feed
        removeNotification(requestId);
        
        // Show success message with new member info
        showToast(`${result.new_member.nickname} joined the group!`, 'success');
        
        // Optionally refresh group member list if viewing group
        if (currentContent.type === 'group_detail' && currentContent.data.groupId === groupId) {
          refreshGroupMembers();
        }
        
      } catch (error) {
        showToast('Failed to approve join request', 'error');
      }
      break;
      
    case 'decline':
      try {
        await fetch(`/join-requests/${requestId}/decline`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        
        // Remove notification from activity feed
        removeNotification(requestId);
        showToast('Join request declined');
        
      } catch (error) {
        showToast('Failed to decline join request', 'error');
      }
      break;
      
    case 'view_profile':
      // Navigate to requester's profile using ContentRenderer
      onContentChange('user_profile', { userId });
      break;
  }
};
```

## Testing Approach

### Backend Testing

#### 1. Join Request Creation Test
```bash
# Test group join request sending
curl -X POST http://localhost:8080/groups/789/join \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=requester_session" \
  -d '{"message": "I would like to join this group"}'

# Expected: Request created in Group_Members table with requested=1
# Expected: WebSocket notification sent to group admins
```

#### 2. Action Endpoint Tests
```bash
# Test approve join request
curl -X POST http://localhost:8080/join-requests/request-id-123/accept \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=admin_session"

# Test decline join request
curl -X POST http://localhost:8080/join-requests/request-id-123/decline \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=admin_session"

# Test get pending requests for group
curl -X GET http://localhost:8080/groups/789/join-requests \
  -H "Cookie: session_id=admin_session"
```

#### 3. Permission Tests
```bash
# Test join request by existing member (should fail)
curl -X POST http://localhost:8080/groups/789/join \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=existing_member_session"

# Expected: 400 Bad Request - Already a member

# Test approve by non-admin (should fail)
curl -X POST http://localhost:8080/join-requests/request-id-123/accept \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=regular_member_session"

# Expected: 403 Forbidden - Insufficient permissions
```

### Frontend Testing

#### 1. Notification Display Test
- User B sends join request to Photography Club
- Verify notification appears in admin's Activity Sidebar
- Check notification contains correct requester data and group info

#### 2. Action Button Test
- Click "Approve" button
- Verify API call is made with correct request_id
- Check user is added to group members
- Confirm notification is removed from activity feed
- Verify success message displays new member name

#### 3. ContentRenderer Integration Test
- Click "View Profile" action
- Verify ContentRenderer loads UserProfile component
- Check correct requester data is passed to component

### Integration Testing

#### 1. End-to-End Flow
1. User B requests to join Photography Club
2. Group admin receives real-time notification
3. Admin clicks "Approve" in Activity Sidebar
4. User B becomes group member
5. User B can now see group posts and events

#### 2. Multiple Admin Notification
1. User B requests to join group with multiple admins
2. All admins receive notification
3. First admin approves request
4. Other admins' notifications are automatically removed

#### 3. Error Scenarios
- Test with invalid request_id
- Test with deleted group
- Test with already processed request
- Test approval by user who lost admin privileges

## Implementation Files

### Backend Files to Create/Modify
```
backend/internal/
├── handlers/
│   └── group_handler.go           # Group join request API endpoints
├── models/
│   └── group_member.go            # Update for join request handling
├── notifications/
│   └── group_notifications.go     # Group join request notification logic
└── websocket/
    └── ws.go                      # Add join request broadcasting
```

### Frontend Files to Create/Modify
```
frontend/components/
├── activity/
│   └── GroupJoinRequestActivity.jsx # Join request notification component
├── groups/
│   ├── GroupDetail.jsx              # Update to show join request button
│   └── GroupJoinRequestsList.jsx    # Admin view of pending requests
└── layout/
    └── ActivitySidebar.jsx          # Add join request handling
```

## Database Schema Requirements

Group join requests use the existing `Group_Members` table from migration 000007, using the `requested` field to distinguish from invitations.

### Existing Group_Members Table Structure
```sql
-- From 000007_create_group_members_table.up.sql
CREATE TABLE Group_Members (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT,
    is_accepted BOOLEAN DEFAULT 0,
    invited_by INTEGER,
    requested BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);
```

### Group Join Request Data Flow
```sql
-- 1. User B (456) requests to join Photography Club (789)
INSERT INTO Group_Members (group_id, user_id, role, is_accepted, invited_by, requested)
VALUES (789, 456, 'member', 0, NULL, 1);

-- 2. Create activity notifications for all group admins
INSERT INTO activity_notifications (user_id, type, from_user_id, reference_id, data)
SELECT u.id, 'group_join_request', 456, 789,
       '{"group_name": "Photography Club", "request_message": "I love photography"}'
FROM Group_Members gm
JOIN Users u ON gm.user_id = u.id
WHERE gm.group_id = 789 AND gm.is_accepted = 1 AND gm.role IN ('admin', 'moderator');

-- 3. Admin approves the request
UPDATE Group_Members 
SET is_accepted = 1, requested = 0
WHERE group_id = 789 AND user_id = 456;

-- 4. Mark all related notifications as read
UPDATE activity_notifications 
SET is_read = 1 
WHERE reference_id = 789 AND type = 'group_join_request' AND from_user_id = 456;
```

### Group Join Request Lifecycle
1. User B requests to join group
2. Row created in `Group_Members` with `is_accepted=0, requested=1`
3. Activity notifications created for all group admins/moderators
4. WebSocket notifications sent to admins
5. Admin approves/declines through Activity Sidebar
6. If approved: Update `is_accepted=1, requested=0` in Group_Members
7. If declined: Delete row from Group_Members
8. All related activity notifications marked as read

### Query Patterns for Group Join Requests
```sql
-- Get pending join requests for a group (admin view)
SELECT gm.*, u.nickname, u.avatar
FROM Group_Members gm
JOIN Users u ON gm.user_id = u.id
WHERE gm.group_id = ? AND gm.requested = 1 AND gm.is_accepted = 0
ORDER BY gm.created_at DESC;

-- Check if user has pending request for group
SELECT * FROM Group_Members 
WHERE group_id = ? AND user_id = ? AND requested = 1 AND is_accepted = 0;

-- Get all admins/moderators for notification
SELECT user_id FROM Group_Members 
WHERE group_id = ? AND is_accepted = 1 AND role IN ('admin', 'moderator');
```

### Permission Requirements
- Only non-members can send join requests
- Only group admins and moderators can approve/decline requests
- Users cannot have multiple pending requests for the same group
- Join requests expire after a configurable time period (e.g., 30 days)

This group join request system complements the group invitation system and provides complete group membership management functionality.
