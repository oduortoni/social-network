# Group Event Created Notification

## Overview

Group event created notifications are sent to all group members when a new event is created in a group they belong to. Members can respond with "Going" or "Not Going" through the Activity Sidebar, and can view event details.

## Notification Data Structure

### WebSocket Notification Format
```json
{
  "type": "notification",
  "subtype": "group_event_created",
  "user_id": 123,
  "nickname": "alice_photo",
  "avatar": "/uploads/avatars/123.jpg",
  "message": "Alice Photo created a new event in Photography Club",
  "timestamp": 1640995200,
  "requires_action": false,
  "actions": ["view_event", "going", "not_going"],
  "additional_data": {
    "event_id": 456,
    "creator_id": 123,
    "group_id": 789,
    "group_name": "Photography Club",
    "event_title": "Photography Workshop",
    "event_description": "Learn advanced photography techniques",
    "event_time": "2024-02-15T14:00:00Z",
    "event_location": "Community Center"
  }
}
```

### Data Fields Explanation
- **user_id**: ID of the user who created the event (creator)
- **nickname**: Display name of the event creator
- **avatar**: Profile picture path of the creator
- **message**: Human-readable notification text
- **timestamp**: Unix timestamp when event was created
- **requires_action**: `false` - informational with optional response
- **actions**: Available actions for group members
- **event_id**: ID of the created event
- **creator_id**: Same as user_id (for clarity)
- **group_id**: ID of the group where event was created
- **group_name**: Name of the group
- **event_title**: Title of the event
- **event_description**: Event description
- **event_time**: Scheduled date/time of the event
- **event_location**: Event location (if specified)

## Required API Endpoints

### 1. Create Group Event
```http
POST /groups/{groupId}/events
Content-Type: application/json
Authorization: Session cookie required

{
  "title": "Photography Workshop",
  "description": "Learn advanced photography techniques",
  "event_time": "2024-02-15T14:00:00Z",
  "location": "Community Center"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Event created successfully",
  "event": {
    "event_id": 456,
    "title": "Photography Workshop",
    "description": "Learn advanced photography techniques",
    "event_time": "2024-02-15T14:00:00Z",
    "location": "Community Center",
    "created_by": 123,
    "group_id": 789
  }
}
```

### 2. Respond to Event
```http
POST /events/{eventId}/respond
Content-Type: application/json
Authorization: Session cookie required

{
  "response": "going"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Event response recorded",
  "response": "going",
  "event": {
    "event_id": 456,
    "title": "Photography Workshop",
    "going_count": 12,
    "not_going_count": 3
  }
}
```

### 3. Get Event Details
```http
GET /events/{eventId}
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "event": {
    "event_id": 456,
    "title": "Photography Workshop",
    "description": "Learn advanced photography techniques",
    "event_time": "2024-02-15T14:00:00Z",
    "location": "Community Center",
    "created_by": 123,
    "creator": {
      "nickname": "alice_photo",
      "avatar": "/uploads/avatars/123.jpg"
    },
    "group": {
      "group_id": 789,
      "group_name": "Photography Club"
    },
    "responses": {
      "going_count": 12,
      "not_going_count": 3,
      "user_response": "going"
    },
    "attendees": [
      {
        "user_id": 456,
        "nickname": "bob_wilson",
        "avatar": "/uploads/avatars/456.jpg",
        "response": "going"
      }
    ]
  }
}
```

### 4. Get Group Events
```http
GET /groups/{groupId}/events?upcoming=true
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "events": [
    {
      "event_id": 456,
      "title": "Photography Workshop",
      "event_time": "2024-02-15T14:00:00Z",
      "location": "Community Center",
      "going_count": 12,
      "not_going_count": 3,
      "user_response": "going"
    }
  ]
}
```

## Call-to-Action Implementation

### Activity Sidebar Actions
The group event notification displays three action buttons:

1. View Event - Opens the event details in ContentRenderer
2. Going - Marks user as attending the event
3. Not Going - Marks user as not attending the event

### Frontend Action Handlers
```javascript
// In GroupEventActivity component
const handleEventAction = async (action, eventId, groupId) => {
  switch (action) {
    case 'view_event':
      // Navigate to event details using ContentRenderer
      onContentChange('event_detail', { eventId });
      break;
      
    case 'going':
    case 'not_going':
      try {
        const response = await fetch(`/events/${eventId}/respond`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ response: action === 'going' ? 'going' : 'not_going' })
        });
        
        const result = await response.json();
        
        // Update notification to show user's response
        updateNotificationResponse(eventId, action);
        
        // Show success message
        showToast(`Marked as ${action === 'going' ? 'Going' : 'Not Going'}`, 'success');
        
        // Update event counts if viewing event details
        if (currentContent.type === 'event_detail' && currentContent.data.eventId === eventId) {
          updateEventCounts(result.event.going_count, result.event.not_going_count);
        }
        
      } catch (error) {
        showToast('Failed to record event response', 'error');
      }
      break;
  }
};
```

## Testing Approach

### Backend Testing

#### 1. Event Creation Test
```bash
# Test group event creation
curl -X POST http://localhost:8080/groups/789/events \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=creator_session" \
  -d '{
    "title": "Photography Workshop",
    "description": "Learn advanced techniques",
    "event_time": "2024-02-15T14:00:00Z",
    "location": "Community Center"
  }'

# Expected: Event created in Group_Events table
# Expected: WebSocket notifications sent to all group members
```

#### 2. Event Response Test
```bash
# Test event response
curl -X POST http://localhost:8080/events/456/respond \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=member_session" \
  -d '{"response": "going"}'

# Expected: Response recorded in Event_Responses table
# Expected: Response counts updated
```

#### 3. Event Details Test
```bash
# Test event details retrieval
curl -X GET http://localhost:8080/events/456 \
  -H "Cookie: session_id=member_session"

# Expected: Complete event details with response counts
```

#### 4. Permission Tests
```bash
# Test event creation by non-member (should fail)
curl -X POST http://localhost:8080/groups/789/events \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=non_member_session" \
  -d '{"title": "Unauthorized Event"}'

# Expected: 403 Forbidden error
```

### Frontend Testing

#### 1. Notification Display Test
- Admin creates event in Photography Club
- Verify notification appears in all members' Activity Sidebars
- Check notification contains correct event data and creator info

#### 2. Action Button Test
- Click "Going" button
- Verify API call is made with correct event_id
- Check user's response is recorded
- Confirm notification updates to show user's response
- Verify success message displays

#### 3. ContentRenderer Integration Test
- Click "View Event" action
- Verify ContentRenderer loads EventDetail component
- Check correct event data is passed to component

### Integration Testing

#### 1. End-to-End Flow
1. Group admin creates Photography Workshop event
2. All group members receive real-time notifications
3. Member clicks "Going" in Activity Sidebar
4. Response is recorded and counts updated
5. Member can view event details with attendee list

#### 2. Multiple Response Updates
1. Member initially responds "Going"
2. Member changes response to "Not Going"
3. Response counts update correctly
4. Event attendee list reflects changes

#### 3. Error Scenarios
- Test with invalid event_id
- Test response by non-group member
- Test event creation in non-existent group
- Test response to deleted event

## Implementation Files

### Backend Files to Create/Modify
```
backend/internal/
├── handlers/
│   └── event_handler.go           # Event API endpoints
├── models/
│   ├── group_event.go             # Event data model
│   └── event_response.go          # Event response model
├── notifications/
│   └── event_notifications.go     # Event notification logic
└── websocket/
    └── ws.go                      # Add event notification broadcasting
```

### Frontend Files to Create/Modify
```
frontend/components/
├── activity/
│   └── GroupEventActivity.jsx     # Event notification component
├── events/
│   ├── EventDetail.jsx            # Event details component
│   └── EventsList.jsx             # Group events list
└── layout/
    └── ActivitySidebar.jsx        # Add event notification handling
```

## Database Schema Requirements

Group events use the existing `Group_Events` and `Event_Responses` tables from migrations 000009 and 000010.

### Existing Group_Events Table Structure
```sql
-- From 000009_create_group_events_table.up.sql
CREATE TABLE Group_Events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    event_time DATETIME,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES Users(id) ON DELETE CASCADE
);
```

### Existing Event_Responses Table Structure
```sql
-- From 000010_create_event_responses_table.up.sql
CREATE TABLE Event_Responses (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT,
    responded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES Group_Events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);
```

### Group Event Notification Data Flow
```sql
-- 1. Admin (123) creates event in Photography Club (789)
INSERT INTO Group_Events (group_id, title, description, event_time, created_by)
VALUES (789, 'Photography Workshop', 'Learn advanced techniques', '2024-02-15 14:00:00', 123);

-- 2. Create activity notifications for all group members (except creator)
INSERT INTO activity_notifications (user_id, type, from_user_id, reference_id, data)
SELECT gm.user_id, 'group_event_created', 123, LAST_INSERT_ROWID(),
       '{"event_title": "Photography Workshop", "group_name": "Photography Club", "event_time": "2024-02-15T14:00:00Z"}'
FROM Group_Members gm
WHERE gm.group_id = 789 AND gm.is_accepted = 1 AND gm.user_id != 123;

-- 3. Member (456) responds to event
INSERT INTO Event_Responses (event_id, user_id, response)
VALUES (?, 456, 'going')
ON CONFLICT (event_id, user_id) DO UPDATE SET 
    response = excluded.response, 
    responded_at = CURRENT_TIMESTAMP;
```

### Group Event Lifecycle
1. Group member creates event in group
2. Event stored in `Group_Events` table
3. Activity notifications created for all group members (except creator)
4. WebSocket notifications sent to all members
5. Members respond through Activity Sidebar or event details
6. Responses stored in `Event_Responses` table
7. Event details show response counts and attendee lists

### Query Patterns for Group Events
```sql
-- Get upcoming events for a group
SELECT ge.*, u.nickname as creator_name, u.avatar as creator_avatar,
       COUNT(CASE WHEN er.response = 'going' THEN 1 END) as going_count,
       COUNT(CASE WHEN er.response = 'not_going' THEN 1 END) as not_going_count
FROM Group_Events ge
JOIN Users u ON ge.created_by = u.id
LEFT JOIN Event_Responses er ON ge.id = er.event_id
WHERE ge.group_id = ? AND ge.event_time > CURRENT_TIMESTAMP
GROUP BY ge.id
ORDER BY ge.event_time ASC;

-- Get user's response to specific event
SELECT response FROM Event_Responses 
WHERE event_id = ? AND user_id = ?;

-- Get all attendees for an event
SELECT u.id, u.nickname, u.avatar, er.response
FROM Event_Responses er
JOIN Users u ON er.user_id = u.id
WHERE er.event_id = ? AND er.response = 'going'
ORDER BY u.nickname;
```

### Permission Requirements
- Only group members can create events
- Only group members can respond to events
- Event creators can edit/delete their events
- Group admins can manage all group events

This group event system provides comprehensive event management with notification support and response tracking for group activities.
