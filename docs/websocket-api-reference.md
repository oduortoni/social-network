# WebSocket API Reference

## WebSocket Connection

### Endpoint
```
ws://localhost:9000/ws
```

### Authentication
All WebSocket connections require session-based authentication via cookies.

```javascript
// Frontend example
const ws = new WebSocket('ws://localhost:9000/ws');
// Session cookie is automatically included
```

### Connection Lifecycle

1. **Connection Establishment**
   - Client sends WebSocket upgrade request
   - Server validates session cookie
   - If valid, connection is established and client is registered
   - If invalid, connection is immediately closed

2. **Message Exchange**
   - Client can send messages in JSON format
   - Server routes messages based on type
   - Server can send notifications to client

3. **Connection Termination**
   - Client disconnects or connection drops
   - Server automatically unregisters client
   - Resources are cleaned up

## Message Format

All WebSocket messages use JSON format:

```json
{
  "type": "message_type",
  "content": "message content",
  "timestamp": 1640995200,
  // Additional fields based on message type
}
```

## Client-to-Server Messages

### Private Message
Send a direct message to another user.

```json
{
  "type": "private",
  "to": 123,
  "content": "Hello there!",
  "timestamp": 1640995200
}
```

**Fields:**
- `type` (string, required): Must be "private"
- `to` (integer, required): Recipient user ID
- `content` (string, required): Message content
- `timestamp` (integer, optional): Will be set by server if not provided

**Behavior:**
- Message is sent only to the specified recipient if they are online
- Message is persisted to database
- Sender receives no confirmation

### Group Message
Send a message to all members of a group.

```json
{
  "type": "group",
  "group_id": "456",
  "content": "Hello everyone!",
  "timestamp": 1640995200
}
```

**Fields:**
- `type` (string, required): Must be "group"
- `group_id` (string, required): Target group ID
- `content` (string, required): Message content
- `timestamp` (integer, optional): Will be set by server if not provided

**Behavior:**
- Message is broadcast to all accepted group members except sender
- Message is persisted to database
- Only group members receive the message

### Broadcast Message
Send a message to all connected users (admin/system use).

```json
{
  "type": "broadcast",
  "content": "System maintenance in 10 minutes",
  "timestamp": 1640995200
}
```

**Fields:**
- `type` (string, required): Must be "broadcast"
- `content` (string, required): Message content
- `timestamp` (integer, optional): Will be set by server if not provided

**Behavior:**
- Message is sent to all currently connected users
- Message is NOT persisted to database
- Typically used for system announcements

## Server-to-Client Messages

### Private Message
Received when another user sends you a private message.

```json
{
  "type": "private",
  "to": 123,
  "content": "Hello there!",
  "timestamp": 1640995200
}
```

### Group Message
Received when someone posts in a group you're a member of.

```json
{
  "type": "group",
  "group_id": "456",
  "content": "Hello everyone!",
  "timestamp": 1640995200
}
```

### Broadcast Message
Received when system sends announcement to all users.

```json
{
  "type": "broadcast",
  "content": "System maintenance in 10 minutes",
  "timestamp": 1640995200
}
```

### Notification
Received for various system events and user actions.

```json
{
  "type": "notification",
  "subtype": "group_invite",
  "message": "You were invited to join group 'Developers'",
  "group_id": 789,
  "timestamp": 1640995200
}
```

**Common Notification Subtypes:**
- `group_invite`: Group invitation received
- `friend_request`: Friend request received
- `friend_accept`: Friend request accepted
- `new_post`: New post from followed user
- `post_like`: Someone liked your post
- `post_comment`: Someone commented on your post
- `new_event`: New event in your group
- `event_rsvp`: Someone RSVP'd to your event

## HTTP API Endpoints

### Get Private Messages
Retrieve message history with a specific user.

```http
GET /api/messages/private?user={userId}&limit={limit}&offset={offset}
```

**Parameters:**
- `user` (integer, required): Other user's ID
- `limit` (integer, optional): Number of messages to retrieve (default: 50)
- `offset` (integer, optional): Number of messages to skip (default: 0)

**Response:**
```json
[
  {
    "type": "private",
    "to": 123,
    "content": "Hello there!",
    "timestamp": 1640995200
  }
]
```

### Get Group Messages
Retrieve message history for a group.

```http
GET /api/messages/group?group={groupId}&limit={limit}&offset={offset}
```

**Parameters:**
- `group` (integer, required): Group ID
- `limit` (integer, optional): Number of messages to retrieve (default: 50)
- `offset` (integer, optional): Number of messages to skip (default: 0)

**Response:**
```json
[
  {
    "type": "group",
    "group_id": "456",
    "content": "Hello everyone!",
    "timestamp": 1640995200
  }
]
```

### Send Group Invitation
Invite a user to join a group.

```http
POST /api/groups/invite
Content-Type: application/json

{
  "group_id": 1,
  "user_id": 2,
  "group_name": "Developers"
}
```

**Response:**
- `201 Created`: Invitation sent successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error

**Behavior:**
- Creates group membership record with `is_accepted: false`
- Stores notification in database
- Sends real-time notification if user is online

### Get Notifications
Retrieve user's notifications.

```http
GET /api/notifications?limit={limit}&offset={offset}
```

**Parameters:**
- `limit` (integer, optional): Number of notifications to retrieve (default: 20)
- `offset` (integer, optional): Number of notifications to skip (default: 0)

**Response:**
```json
[
  {
    "id": 1,
    "type": "group_invite",
    "message": "You were invited to join group 'Developers'",
    "is_read": false,
    "created_at": "2023-01-01T12:00:00Z"
  }
]
```

### Mark Notifications as Read
Mark all notifications as read for the authenticated user.

```http
POST /api/notifications/read
```

**Response:**
- `204 No Content`: Notifications marked as read
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error

## Error Handling

### WebSocket Errors
- **Connection Refused**: Invalid or missing session cookie
- **Message Ignored**: Invalid JSON format or missing required fields
- **Silent Failures**: Messages to offline users are silently dropped

### HTTP API Errors
All HTTP endpoints return standard HTTP status codes:

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Request successful, no content to return
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `500 Internal Server Error`: Server error

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting for production use:

- WebSocket messages: Limit messages per user per minute
- HTTP API calls: Standard rate limiting per endpoint
- Notification sending: Prevent spam notifications

## Security Considerations

1. **Authentication**: All connections require valid session cookies
2. **Authorization**: Users can only send messages to groups they're members of
3. **Input Validation**: All message content should be validated and sanitized
4. **CORS**: WebSocket upgrader currently allows all origins (development setting)

## Monitoring and Debugging
- Monitor database for message persistence
