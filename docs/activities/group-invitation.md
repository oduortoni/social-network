# Group Invitation Notification

## Overview

Group invitation notifications are sent when a user is invited to join a group. The recipient can accept or decline the invitation through the Activity Sidebar. Upon acceptance, the user becomes a member of the group.

## Notification Data Structure

### WebSocket Notification Format
```json
{
  "type": "notification",
  "subtype": "group_invitation",
  "user_id": 123,
  "nickname": "jane_smith",
  "avatar": "/uploads/avatars/123.jpg",
  "message": "Jane Smith invited you to join Photography Club",
  "timestamp": 1640995200,
  "requires_action": true,
  "actions": ["accept", "decline", "view_group"],
  "additional_data": {
    "notification_id": "uuid-def456",
    "inviter_id": 123,
    "group_id": "group-uuid-789",
    "group_name": "Photography Club",
    "group_avatar": "/uploads/groups/789.jpg",
    "inviter_role": "admin"
  }
}
```

### Data Fields Explanation
- **user_id**: ID of the user who sent the invitation (inviter)
- **nickname**: Display name of the inviter
- **avatar**: Profile picture path of the inviter
- **message**: Human-readable notification text
- **timestamp**: Unix timestamp when invitation was sent
- **requires_action**: Always `true` for group invitations
- **actions**: Available actions for the recipient
- **notification_id**: Unique identifier from temporary_notifications table
- **inviter_id**: Same as user_id (for clarity)
- **group_id**: ID of the group being invited to
- **group_name**: Name of the group
- **group_avatar**: Group's profile picture
- **inviter_role**: Role of the inviter in the group (admin, moderator, member)

## Required API Endpoints

### 1. Send Group Invitation
```http
POST /groups/{groupId}/invites
Content-Type: application/json
Authorization: Session cookie required

{
  "user_id": 456
}
```

**Response:**
```json
{
  "success": true,
  "message": "Group invitation sent",
  "notification_id": "uuid-def456"
}
```

### 2. Accept Group Invitation
```http
POST /group-invitations/{notificationId}/accept
Content-Type: application/json
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Group invitation accepted",
  "group": {
    "group_id": "group-uuid-789",
    "group_name": "Photography Club",
    "group_avatar": "/uploads/groups/789.jpg",
    "member_count": 25
  }
}
```

### 3. Decline Group Invitation
```http
POST /group-invitations/{notificationId}/decline
Content-Type: application/json
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "message": "Group invitation declined"
}
```

### 4. Get Pending Group Invitations
```http
GET /group-invitations
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "invitations": [
    {
      "notification_id": "uuid-def456",
      "inviter": {
        "user_id": 123,
        "nickname": "jane_smith",
        "avatar": "/uploads/avatars/123.jpg",
        "role": "admin"
      },
      "group": {
        "group_id": "group-uuid-789",
        "group_name": "Photography Club",
        "group_avatar": "/uploads/groups/789.jpg",
        "description": "A community for photography enthusiasts",
        "member_count": 24
      },
      "created_at": 1640995200
    }
  ]
}
```

## Call-to-Action Implementation

### Activity Sidebar Actions
The group invitation notification displays three action buttons:

1. Accept - Accepts the group invitation and joins the group
2. Decline - Declines the group invitation
3. View Group - Opens the group details in ContentRenderer

### Frontend Action Handlers
```javascript
// In GroupInvitationActivity component
const handleGroupInviteAction = async (action, notificationId, groupId, userId) => {
  switch (action) {
    case 'accept':
      try {
        const response = await fetch(`/group-invitations/${notificationId}/accept`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        
        const result = await response.json();
        
        // Remove notification from activity feed
        removeNotification(notificationId);
        
        // Show success message with group info
        showToast(`Welcome to ${result.group.group_name}!`, 'success');
        
        // Optionally navigate to group
        onContentChange('group_detail', { groupId });
        
      } catch (error) {
        showToast('Failed to accept invitation', 'error');
      }
      break;
      
    case 'decline':
      try {
        await fetch(`/group-invitations/${notificationId}/decline`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        
        // Remove notification from activity feed
        removeNotification(notificationId);
        showToast('Group invitation declined');
        
      } catch (error) {
        showToast('Failed to decline invitation', 'error');
      }
      break;
      
    case 'view_group':
      // Navigate to group details using ContentRenderer
      onContentChange('group_detail', { groupId });
      break;
  }
};
```

## Testing Approach

### Backend Testing

#### 1. Invitation Creation Test
```bash
# Test group invitation sending
curl -X POST http://localhost:8080/groups/group-uuid-789/invites \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=inviter_session" \
  -d '{"user_id": 456}'

# Expected: Invitation created in temporary_notifications table
# Expected: WebSocket notification sent to invited user
```

#### 2. Action Endpoint Tests
```bash
# Test accept invitation
curl -X POST http://localhost:8080/group-invitations/uuid-def456/accept \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=invitee_session"

# Test decline invitation
curl -X POST http://localhost:8080/group-invitations/uuid-def456/decline \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=invitee_session"

# Test get pending invitations
curl -X GET http://localhost:8080/group-invitations \
  -H "Cookie: session_id=invitee_session"
```

#### 3. Permission Tests
```bash
# Test invitation by non-member (should fail)
curl -X POST http://localhost:8080/groups/group-uuid-789/invites \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=non_member_session" \
  -d '{"user_id": 456}'

# Expected: 403 Forbidden error
```

### Frontend Testing

#### 1. Notification Display Test
- Admin sends group invitation to User B
- Verify notification appears in User B's Activity Sidebar
- Check notification contains correct group data and inviter info

#### 2. Action Button Test
- Click "Accept" button
- Verify API call is made with correct notification_id
- Check user is added to group members
- Confirm notification is removed from activity feed
- Verify success message displays group name

#### 3. ContentRenderer Integration Test
- Click "View Group" action
- Verify ContentRenderer loads GroupDetail component
- Check correct group data is passed to component

### Integration Testing

#### 1. End-to-End Flow
1. Group admin invites User B to Photography Club
2. User B receives real-time notification
3. User B clicks "Accept" in Activity Sidebar
4. User B becomes group member
5. User B can now see group posts and events

#### 2. Permission Validation
1. Non-member tries to invite someone to group
2. System rejects invitation with proper error
3. Only admins/moderators can send invitations

#### 3. Error Scenarios
- Test with invalid notification_id
- Test with deleted group
- Test with already processed invitation
- Test invitation to user already in group

## Implementation Files

### Backend Files to Create/Modify
```
backend/internal/
├── handlers/
│   └── group_handler.go           # Group invitation API endpoints
├── models/
│   ├── group.go                   # Group data model
│   └── group_member.go            # Group membership model
├── notifications/
│   └── group_notifications.go     # Group notification logic
└── websocket/
    └── ws.go                      # Add group notification broadcasting
```

### Frontend Files to Create/Modify
```
frontend/components/
├── activity/
│   └── GroupInvitationActivity.jsx # Group invitation notification component
├── groups/
│   ├── GroupDetail.jsx             # Update to show invitation button
│   └── GroupMembersList.jsx        # Group members management
└── layout/
    └── ActivitySidebar.jsx         # Add group invitation handling
```

## Database Schema Requirements

Group invitations use the existing `Group_Members` table from migration 000007. This table already handles both group invitations and join requests with the `invited_by` and `requested` fields.

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
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (invited_by) REFERENCES Users(id) ON DELETE SET NULL
);
```

### Group Invitation Data Flow
```sql
-- 1. Admin (123) invites User B (456) to Photography Club (789)
INSERT INTO Group_Members (group_id, user_id, role, is_accepted, invited_by, requested)
VALUES (789, 456, 'member', 0, 123, 0);

-- 2. Create activity notification for User B
INSERT INTO activity_notifications (user_id, type, from_user_id, reference_id, data)
VALUES (456, 'group_invitation', 123, 789,
        '{"group_name": "Photography Club", "inviter_role": "admin"}');

-- 3. User B accepts the invitation
UPDATE Group_Members
SET is_accepted = 1
WHERE group_id = 789 AND user_id = 456;

-- 4. Mark notification as read
UPDATE activity_notifications
SET is_read = 1
WHERE user_id = 456 AND reference_id = 789 AND type = 'group_invitation';
```

### Group Invitation Lifecycle
1. Group admin/moderator invites User B to group
2. Row created in `Group_Members` with `is_accepted=0, invited_by=admin_id`
3. Activity notification created for rich WebSocket data
4. WebSocket notification sent to User B
5. User B accepts/declines through Activity Sidebar
6. If accepted: Update `is_accepted=1` in Group_Members
7. If declined: Delete row from Group_Members
8. Activity notification marked as read

### Query Patterns for Group Invitations
```sql
-- Get pending group invitations for a user
SELECT gm.*, g.name as group_name, g.avatar as group_avatar,
       u.nickname as inviter_nickname, u.avatar as inviter_avatar
FROM Group_Members gm
JOIN Groups g ON gm.group_id = g.id
JOIN Users u ON gm.invited_by = u.id
WHERE gm.user_id = ? AND gm.is_accepted = 0 AND gm.invited_by IS NOT NULL
ORDER BY gm.created_at DESC;

-- Check if user is already in group
SELECT is_accepted FROM Group_Members
WHERE group_id = ? AND user_id = ?;

-- Get group member count
SELECT COUNT(*) FROM Group_Members
WHERE group_id = ? AND is_accepted = 1;
```

### Permission Requirements
- Only group admins and moderators can send invitations
- Users cannot be invited to groups they're already members of
- Invitations expire after a configurable time period (e.g., 30 days)

This group invitation system provides the foundation for group membership management and demonstrates how temporary notifications create permanent relationships when accepted.
