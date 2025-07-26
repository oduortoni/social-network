# Notification System Documentation

## Overview

The notification system provides real-time notifications to users through WebSocket connections. It features a clean bell icon interface with red badge indicators and supports extensible notification types for various social network events.

## Architecture

### Backend Components
- **WebSocket Manager**: Handles user connections and broadcasting
- **Notification Broadcasting**: Sends notifications to individual users, groups, or all users
- **Session Management**: Authenticates users and manages connections

### Frontend Components
- **useSimpleNotifications Hook**: Manages notification state and WebSocket listening
- **Header Bell Icon**: Displays notifications with red badge and dropdown panel
- **Real-time Updates**: Instant notification delivery and UI updates

## Notification Format

All notifications follow a standardized JSON format:

```json
{
  "type": "notification",
  "subtype": "notification_category",
  "user_id": 123,
  "user_name": "John Doe",
  "message": "Human-readable notification message",
  "timestamp": 1640995200,
  "additional_data": {
    "post_id": 456,
    "group_id": 789
  }
}
```

### Required Fields
- **type**: Always "notification"
- **subtype**: Notification category (see types below)
- **message**: User-friendly message text
- **timestamp**: Unix timestamp

### Optional Fields
- **user_id**: ID of user who triggered the notification
- **user_name**: Display name of triggering user
- **additional_data**: Extra context data specific to notification type

## Notification Types

### Currently Implemented
- **user_connected**: User comes online
- **user_disconnected**: User goes offline

### You can add other types for example
- **follow_request**: New follow request received
- **follow_accepted**: Follow request accepted
- **group_invite**: Invitation to join group
- **new_message**: New private or group message
- **post_like**: Someone liked your post
- **post_comment**: Someone commented on your post
- **event_invite**: Invitation to event

## Backend Usage

### Broadcasting to All Users

```go
// Example: Post creation notification
postNotification := map[string]interface{}{
    "type":      "notification",
    "subtype":   "new_post",
    "user_id":   authorID,
    "user_name": authorName,
    "message":   authorName + " created a new post",
    "timestamp": time.Now().Unix(),
    "additional_data": map[string]interface{}{
        "post_id": postID,
        "post_title": postTitle,
    },
}

// Broadcast to all users except the author
wsManager.broadcastNotificationToAll(postNotification, authorID)
```

### Sending to Specific User

```go
// Example: Follow request notification
followRequestNotification := map[string]interface{}{
    "type":      "notification",
    "subtype":   "follow_request",
    "user_id":   senderID,
    "user_name": senderName,
    "message":   senderName + " sent you a follow request",
    "timestamp": time.Now().Unix(),
    "additional_data": map[string]interface{}{
        "request_id": requestID,
    },
}

// Send to specific user
wsManager.SendNotificationToUser(recipientID, followRequestNotification)
```

### Sending to Group Members

```go
// Example: Group message notification
groupMessageNotification := map[string]interface{}{
    "type":      "notification",
    "subtype":   "group_message",
    "user_id":   senderID,
    "user_name": senderName,
    "message":   senderName + " sent a message to " + groupName,
    "timestamp": time.Now().Unix(),
    "additional_data": map[string]interface{}{
        "group_id": groupID,
        "message_id": messageID,
    },
}

// Send to all group members except sender
wsManager.SendNotificationToGroup(groupID, groupMessageNotification, senderID)
```

## Frontend Usage

### Using the Notification Hook

```javascript
import { useSimpleNotifications } from '../hooks/useNotifications';

const MyComponent = () => {
  const { unreadCount, notifications, markAllAsRead, clearNotifications } = useSimpleNotifications();

  return (
    <div>
      <span>Unread: {unreadCount}</span>
      <button onClick={markAllAsRead}>Mark All Read</button>
      <button onClick={clearNotifications}>Clear All</button>
      
      {notifications.map((notification, index) => (
        <div key={index}>
          <p>{notification.message}</p>
          <small>{new Date(notification.timestamp * 1000).toLocaleString()}</small>
        </div>
      ))}
    </div>
  );
};
```

### Bell Icon Integration

The notification system is automatically integrated into the Header component. The bell icon:
- Shows red badge when unread notifications exist
- Displays count (1, 2, 3... or 9+ for 10+)
- Opens dropdown panel on click
- Auto-marks notifications as read when opened

## Extending the System

### Adding New Notification Types

1. **Define the notification type** in your backend handler:

```go
// Example: Post like notification
likeNotification := map[string]interface{}{
    "type":      "notification",
    "subtype":   "post_like",
    "user_id":   likerID,
    "user_name": likerName,
    "message":   likerName + " liked your post",
    "timestamp": time.Now().Unix(),
    "additional_data": map[string]interface{}{
        "post_id": postID,
        "like_id": likeID,
    },
}

// Send to post author
wsManager.SendNotificationToUser(postAuthorID, likeNotification)
```

2. **Handle in frontend** (optional custom handling):

```javascript
// The system automatically handles all notification types
// Custom handling only needed for special UI behavior

useEffect(() => {
  const handleSpecialNotification = (notification) => {
    if (notification.subtype === 'post_like') {
      // Custom handling for post likes
      console.log('Someone liked your post!');
    }
  };

  wsService.onMessage('notification', handleSpecialNotification);
}, []);
```

### Custom Notification Actions

For notifications requiring user action (accept/decline):

```go
// Backend: Action-required notification
actionNotification := map[string]interface{}{
    "type":      "notification",
    "subtype":   "follow_request",
    "user_id":   senderID,
    "user_name": senderName,
    "message":   senderName + " sent you a follow request",
    "timestamp": time.Now().Unix(),
    "requires_action": true,
    "additional_data": map[string]interface{}{
        "request_id": requestID,
        "actions": []string{"accept", "decline"},
    },
}
```

## Best Practices

### Backend
- Always exclude the triggering user from broadcasts
- Include relevant IDs in additional_data for frontend routing
- Use descriptive, user-friendly messages
- Set appropriate timestamps
- Handle WebSocket connection failures gracefully

### Frontend
- Keep notification messages concise and clear
- Provide visual feedback for user actions
- Implement proper error handling for failed notifications
- Consider notification persistence for offline users
- Respect user notification preferences

### Performance
- Limit notification history (current: 20 notifications)
- Use efficient broadcasting to avoid overwhelming users
- Implement notification batching for high-frequency events
- Consider notification priorities for important events

