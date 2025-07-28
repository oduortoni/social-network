# New Message Notification

## Overview

New message notifications are sent when a user receives a private message or group chat message. These notifications provide real-time awareness of incoming messages and allow users to quickly navigate to the conversation.

## Notification Data Structure

### WebSocket Notification Format
```json
{
  "type": "notification",
  "subtype": "new_message",
  "user_id": 123,
  "nickname": "john_doe",
  "avatar": "/uploads/avatars/123.jpg",
  "message": "John Doe sent you a message",
  "timestamp": 1640995200,
  "requires_action": false,
  "actions": ["view_message", "reply"],
  "additional_data": {
    "message_id": "msg-uuid-456",
    "conversation_type": "private", // or "group"
    "conversation_id": "conv-uuid-789",
    "message_preview": "Hey, how are you doing?",
    "group_name": null // Only for group messages
  }
}
```

### Data Fields Explanation
- **user_id**: ID of the user who sent the message
- **nickname**: Display name of the message sender
- **avatar**: Profile picture path of the sender
- **message**: Human-readable notification text
- **timestamp**: Unix timestamp when message was sent
- **requires_action**: Always `false` for message notifications (informational)
- **actions**: Available actions for the recipient
- **message_id**: Unique identifier of the actual message
- **conversation_type**: "private" for direct messages, "group" for group chat
- **conversation_id**: ID of the conversation (user_id for private, group_id for group)
- **message_preview**: First 50 characters of the message content
- **group_name**: Name of the group (only for group messages)

## Required API Endpoints

### 1. Send Private Message
```http
POST /chats/{userId}/messages
Content-Type: application/json
Authorization: Session cookie required

{
  "content": "Hey, how are you doing?",
  "message_type": "text"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message sent",
  "message_id": "msg-uuid-456",
  "timestamp": 1640995200
}
```

### 2. Get Message History
```http
GET /chats/{userId}/messages?limit=50&offset=0
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "messages": [
    {
      "message_id": "msg-uuid-456",
      "sender_id": 123,
      "content": "Hey, how are you doing?",
      "message_type": "text",
      "timestamp": 1640995200,
      "is_read": false
    }
  ],
  "has_more": false
}
```

### 3. Mark Messages as Read
```http
POST /chats/{userId}/messages/read
Content-Type: application/json
Authorization: Session cookie required

{
  "message_ids": ["msg-uuid-456", "msg-uuid-457"]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Messages marked as read"
}
```

### 4. Get Active Conversations
```http
GET /chats
Authorization: Session cookie required
```

**Response:**
```json
{
  "success": true,
  "conversations": [
    {
      "conversation_id": "conv-uuid-789",
      "type": "private",
      "participant": {
        "user_id": 123,
        "nickname": "john_doe",
        "avatar": "/uploads/avatars/123.jpg"
      },
      "last_message": {
        "content": "Hey, how are you doing?",
        "timestamp": 1640995200,
        "sender_id": 123
      },
      "unread_count": 2
    }
  ]
}
```

## Call-to-Action Implementation

### Activity Sidebar Actions
The new message notification displays two action buttons:

1. View Message - Opens the chat conversation in ContentRenderer
2. Reply - Opens the chat conversation with input focused

### Frontend Action Handlers
```javascript
// In NewMessageActivity component
const handleMessageAction = async (action, conversationId, userId, conversationType) => {
  switch (action) {
    case 'view_message':
      if (conversationType === 'private') {
        // Navigate to private chat
        onContentChange('chat', { userId });
      } else {
        // Navigate to group chat
        onContentChange('group_chat', { groupId: conversationId });
      }
      
      // Mark messages as read when viewing
      try {
        await fetch(`/chats/${userId}/messages/read`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ message_ids: [messageId] })
        });
        // Remove notification from activity feed
        removeNotification(notificationId);
      } catch (error) {
        console.error('Failed to mark message as read:', error);
      }
      break;
      
    case 'reply':
      // Same as view_message but with input focused
      if (conversationType === 'private') {
        onContentChange('chat', { userId, focusInput: true });
      } else {
        onContentChange('group_chat', { groupId: conversationId, focusInput: true });
      }
      break;
  }
};
```

## Testing Approach

### Backend Testing

#### 1. Message Creation Test
```bash
# Test private message sending
curl -X POST http://localhost:8080/chats/456/messages \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=sender_session" \
  -d '{"content": "Hello there!", "message_type": "text"}'

# Expected: Message stored in database
# Expected: WebSocket notification sent to recipient
```

#### 2. Message History Test
```bash
# Test message retrieval
curl -X GET http://localhost:8080/chats/456/messages?limit=10 \
  -H "Cookie: session_id=user_session"

# Expected: Returns message history with correct format
```

#### 3. Read Status Test
```bash
# Test marking messages as read
curl -X POST http://localhost:8080/chats/456/messages/read \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=recipient_session" \
  -d '{"message_ids": ["msg-uuid-456"]}'

# Expected: Messages marked as read in database
```

### Frontend Testing

#### 1. Notification Display Test
- Send message from User A to User B
- Verify notification appears in User B's Activity Sidebar
- Check notification contains correct sender data and message preview

#### 2. Action Button Test
- Click "View Message" button
- Verify ContentRenderer loads ChatWindow component
- Check correct conversation is opened
- Confirm message is marked as read

#### 3. Real-time Updates Test
- Send multiple messages in quick succession
- Verify each message generates a notification
- Check notifications are properly ordered by timestamp

### Integration Testing

#### 1. End-to-End Flow
1. User A sends message to User B
2. User B receives real-time notification
3. User B clicks "View Message" in Activity Sidebar
4. ChatWindow opens with conversation history
5. Message is marked as read automatically

#### 2. Group Message Flow
1. User A sends message to Group X
2. All group members receive notifications
3. User B clicks "View Message"
4. Group chat opens in ContentRenderer

#### 3. Error Scenarios
- Test with invalid user_id
- Test with deleted conversation
- Test WebSocket disconnection during message sending
- Test message sending to blocked user

## Implementation Files

### Backend Files to Create/Modify
```
backend/internal/
├── handlers/
│   └── message_handler.go         # Message API endpoints
├── models/
│   ├── message.go                 # Message data model
│   └── conversation.go            # Conversation data model
├── notifications/
│   └── message_notifications.go   # Message notification logic
└── websocket/
    └── ws.go                      # Add message broadcasting
```

### Frontend Files to Create/Modify
```
frontend/components/
├── activity/
│   └── NewMessageActivity.jsx     # Message notification component
├── chat/
│   ├── ChatWindow.jsx             # Update to handle focus
│   └── GroupChatWindow.jsx        # Group chat component
└── layout/
    └── ActivitySidebar.jsx        # Add message notification handling
```

## Database Schema Requirements

Message notifications do NOT use the temporary_notifications table since messages are permanent data that should be preserved.

### Messages Table
```sql
CREATE TABLE messages (
    id VARCHAR(36) PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER, -- NULL for group messages
    group_id INTEGER, -- NULL for private messages
    content TEXT NOT NULL,
    message_type ENUM('text', 'image', 'file') DEFAULT 'text',
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    
    INDEX idx_private_conversation (sender_id, recipient_id, created_at),
    INDEX idx_group_conversation (group_id, created_at),
    INDEX idx_unread_messages (recipient_id, is_read)
);
```

### Message Read Status
```sql
-- For group messages, track read status per user
CREATE TABLE message_read_status (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    message_id VARCHAR(36) NOT NULL,
    user_id INTEGER NOT NULL,
    read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_read_status (message_id, user_id)
);
```

### Message Notification Lifecycle
1. User A sends message to User B
2. Message stored in `messages` table
3. WebSocket notification sent to User B (not stored in temporary_notifications)
4. User B views message through Activity Sidebar
5. Message marked as read in `messages` table
6. Notification removed from Activity Sidebar (frontend only)

This approach treats message notifications as real-time alerts rather than persistent actionable items, which aligns with typical messaging system behavior.
