# Real-Time Updates Strategy

## Overview

This document outlines our strategic approach to implementing real-time features in the social network application. We employ a hybrid model that combines WebSocket-based real-time updates for critical user interactions with traditional polling/refresh mechanisms for less time-sensitive content.

## Real-Time vs. Polling Decision Matrix

| Feature                 | Real-Time? | Implementation | Reason                                                               |
| ----------------------- | ---------- | -------------- | -------------------------------------------------------------------- |
| **Private Chat**        | Yes      | WebSocket      | Immediate message delivery and typing indicators                     |
| **Group Chat**          | Yes      | WebSocket      | Group-wide real-time updates                                         |
| **Online Users List**   | Yes      | WebSocket      | So users can see who's currently online                              |
| **Notifications**       | Yes      | WebSocket      | Follow requests, group invites, event alerts should pop up instantly |
| **Typing Indicators**   | Yes      | WebSocket      | Nice UX for chat                                                     |
| **New Posts/Comments**  | Notify   | WebSocket      | Server notifies of new content, client decides when to fetch         |
| **Group Event Updates** | Notify   | WebSocket      | Server notifies of updates, client refreshes when appropriate        |

## Strategic Rationale

### Why Real-Time for Chat & Social Interactions

#### **1. User Expectations**
- Modern users expect instant messaging to be truly instant
- Social interactions (notifications) should feel immediate
- Online presence creates engagement and social connection

#### **2. Competitive Advantage**
- Real-time chat is a standard feature expectation
- Instant notifications improve user engagement
- Online status creates FOMO (Fear of Missing Out) effect

#### **3. Technical Feasibility**
- Chat messages are small, frequent, and targeted
- Notifications are event-driven and relatively infrequent
- Online status changes are simple state updates

### Why Notification-Driven Content Updates

#### **1. Smart Resource Usage**
- Server only notifies when new content actually exists
- Client fetches content when user is ready to consume it
- Eliminates unnecessary polling when no updates available

#### **2. Better User Experience**
- Users get notified of new content availability
- Client can choose optimal timing for refresh (user scrolled to top, tab became active)
- Respects user's current context and activity

#### **3. Optimal Scalability**
- Lightweight notifications instead of heavy content broadcasting
- Content fetching happens on-demand via efficient API calls
- Server resources focused on notification delivery, not content broadcasting

## Implementation Architecture

### WebSocket Event Types

```javascript
// Real-time WebSocket events
const REALTIME_EVENTS = {
  // Chat System
  PRIVATE_MESSAGE: 'private_message',
  GROUP_MESSAGE: 'group_message',
  TYPING_START: 'typing_start',
  TYPING_STOP: 'typing_stop',

  // User Presence
  USER_CONNECTED: 'user_connected',
  USER_DISCONNECTED: 'user_disconnected',

  // Notifications
  NOTIFICATION: 'notification',
  FOLLOW_REQUEST: 'follow_request',
  GROUP_INVITE: 'group_invite',
  EVENT_CREATED: 'event_created',

  // Content Update Notifications
  NEW_POSTS_AVAILABLE: 'new_posts_available',
  NEW_COMMENTS_AVAILABLE: 'new_comments_available',
  FEED_UPDATE_AVAILABLE: 'feed_update_available',
  GROUP_CONTENT_UPDATE: 'group_content_update'
};
```

### Connection Management

#### **Ping/Pong Strategy**
- **Interval**: 30 seconds
- **Purpose**: Detect stale connections and maintain connection health
- **Implementation**: Server sends ping, client responds with pong (more reliable)

```go
// Server-side ping implementation
func (m *Manager) startPingPong() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            m.mu.RLock()
            for _, client := range m.clients {
                if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                    // Connection is stale, remove client
                    m.Unregister(client.ID)
                }
            }
            m.mu.RUnlock()
        }
    }()
}
```

```javascript
// Client-side pong response (automatic in most WebSocket implementations)
wsConnection.addEventListener('ping', () => {
  // Browser automatically responds with pong
  // Manual implementation if needed:
  // wsConnection.pong();
});
```

#### **Connection Recovery**
- Automatic reconnection on connection loss
- Exponential backoff for failed reconnections
- State synchronization on reconnection

### Message Broadcasting Patterns

#### **1. Direct User Messaging**
```go
// Send to specific user
wsManager.SendToUser(userID, message)
```

#### **2. Group Broadcasting**
```go
// Send to all group members except sender
wsManager.SendToGroup(groupID, message, excludeUserID)
```

#### **3. Notification Broadcasting**
```go
// Send notification to specific user
wsManager.SendNotificationToUser(userID, notification)
```

## Notification-Driven Content Updates

### Strategy Overview
Instead of traditional polling, the server sends lightweight notifications when new content is available. The client then decides when and how to fetch the actual content based on user context.

### Content Update Notifications

#### **New Posts Available**
```json
{
  "type": "new_posts_available",
  "count": 3,
  "latest_post_id": "uuid-123",
  "timestamp": 1640995200,
  "data": {
    "feed_type": "following", // or "public", "group"
    "group_id": null
  }
}
```

#### **New Comments Available**
```json
{
  "type": "new_comments_available",
  "post_id": "uuid-456",
  "count": 2,
  "timestamp": 1640995200
}
```

### Client-Side Content Refresh Logic

```javascript
// Smart content refresh based on user context
class ContentRefreshManager {
  constructor() {
    this.pendingUpdates = new Map();
    this.isUserActive = true;
    this.isAtTopOfFeed = false;
  }

  handleContentNotification(notification) {
    // Store pending update
    this.pendingUpdates.set(notification.type, notification);

    // Decide when to refresh based on user context
    if (this.shouldRefreshImmediately(notification)) {
      this.refreshContent(notification.type);
    } else {
      this.showUpdateIndicator(notification);
    }
  }

  shouldRefreshImmediately(notification) {
    return (
      this.isUserActive &&
      this.isAtTopOfFeed &&
      notification.count <= 3 // Don't overwhelm with too many updates
    );
  }

  showUpdateIndicator(notification) {
    // Show "New posts available" banner
    // User can click to refresh manually
  }
}
```

## Traditional Polling Fallback Strategy

### Content Updates (Posts/Comments)

#### **Trigger-Based Refresh**
- User action triggers (pull-to-refresh, scroll to top)
- Page focus/visibility change
- Manual refresh button

#### **Smart Caching**
```javascript
// Cache strategy for posts
const CACHE_DURATION = {
  POSTS: 5 * 60 * 1000,      // 5 minutes
  COMMENTS: 2 * 60 * 1000,   // 2 minutes
  PROFILES: 10 * 60 * 1000   // 10 minutes
};
```

#### **Background Sync**
- Periodic background updates for active tabs
- Sync on network reconnection
- Optimistic updates for user actions

### API Polling Implementation

```javascript
// Intelligent polling with backoff
class ContentPoller {
  constructor(endpoint, interval = 30000) {
    this.endpoint = endpoint;
    this.baseInterval = interval;
    this.currentInterval = interval;
    this.isActive = true;
  }
  
  async poll() {
    if (!this.isActive) return;
    
    try {
      const data = await fetch(this.endpoint);
      this.currentInterval = this.baseInterval; // Reset on success
      return data;
    } catch (error) {
      this.currentInterval = Math.min(this.currentInterval * 2, 300000); // Max 5 min
      throw error;
    }
  }
}
```

## Performance Optimizations

### WebSocket Optimizations

#### **Message Batching**
- Batch multiple notifications into single message
- Debounce rapid typing indicators
- Compress message payloads for large groups

#### **Connection Pooling**
- Limit concurrent WebSocket connections per user
- Graceful degradation for connection limits
- Priority queuing for critical messages

### Polling Optimizations

#### **Conditional Requests**
- Use ETags for cache validation
- If-Modified-Since headers
- 304 Not Modified responses

#### **Data Pagination**
- Cursor-based pagination for infinite scroll
- Limit payload sizes
- Progressive loading strategies

## Error Handling & Resilience

### WebSocket Error Scenarios

#### **Connection Failures**
```javascript
// Exponential backoff reconnection
const reconnect = (attempt = 1) => {
  const delay = Math.min(1000 * Math.pow(2, attempt), 30000);
  setTimeout(() => {
    if (attempt <= 5) {
      establishConnection();
    } else {
      // Fall back to polling mode
      enablePollingMode();
    }
  }, delay);
};
```

#### **Message Delivery Failures**
- Message queuing for offline scenarios
- Retry mechanisms with exponential backoff
- Fallback to HTTP API for critical messages

### Polling Error Scenarios

#### **API Failures**
- Graceful degradation to cached data
- User notification of connectivity issues
- Automatic retry with backoff

#### **Network Instability**
- Adaptive polling intervals based on network conditions
- Offline mode with local storage
- Sync queue for pending updates

## Monitoring & Analytics

### Real-Time Metrics
- WebSocket connection count and duration
- Message delivery success rates
- Average message latency
- Connection failure rates

### Polling Metrics
- API response times
- Cache hit/miss ratios
- Polling frequency optimization
- Data freshness metrics

## Security Considerations

### WebSocket Security
- Authentication token validation on connection
- Message origin verification
- Rate limiting for message sending
- Input sanitization for all messages

### Polling Security
- API rate limiting
- CSRF protection
- Data access authorization
- Secure caching strategies

## Future Enhancements

### Potential Real-Time Additions
- **Live Post Reactions**: Real-time like/reaction counts
- **Collaborative Features**: Real-time document editing
- **Live Events**: Real-time event updates and participation

### Advanced Optimizations
- **WebRTC for Direct Chat**: Peer-to-peer messaging for reduced server load
- **Server-Sent Events**: Alternative to WebSocket for one-way updates
- **GraphQL Subscriptions**: Real-time data subscriptions

## Implementation Timeline

### Phase 1: Core Real-Time Features
- [x] WebSocket infrastructure
- [x] Private messaging
- [x] Online user tracking
- [x] Basic notifications

### Phase 2: Enhanced Features
- [ ] Group chat implementation
- [ ] Typing indicators
- [ ] Advanced notification types
- [ ] Connection resilience

### Phase 3: Optimization
- [ ] Message batching
- [ ] Performance monitoring
- [ ] Advanced caching strategies
- [ ] Scalability improvements

## Technical Implementation Details

### WebSocket Message Format

```json
{
  "type": "private_message",
  "timestamp": 1640995200,
  "from_user_id": 123,
  "to_user_id": 456,
  "data": {
    "content": "Hello!",
    "message_id": "uuid-here"
  }
}
```

### Notification Message Format

```json
{
  "type": "notification",
  "subtype": "follow_request",
  "user_id": 123,
  "user_name": "John Doe",
  "message": "John Doe sent you a follow request",
  "timestamp": 1640995200,
  "additional_data": {
    "request_id": "uuid-here",
    "actions": ["accept", "decline"]
  }
}
```

### Frontend Integration

```javascript
// Real-time hook for components
const useRealTimeUpdates = (eventTypes) => {
  const [data, setData] = useState({});

  useEffect(() => {
    eventTypes.forEach(type => {
      wsService.onMessage(type, (message) => {
        setData(prev => ({
          ...prev,
          [type]: [...(prev[type] || []), message]
        }));
      });
    });
  }, [eventTypes]);

  return data;
};
```

## Conclusion

This hybrid approach balances user experience expectations with technical constraints and scalability requirements. By focusing real-time capabilities on truly interactive features while using efficient polling for content updates, we create a responsive social network that can scale effectively while providing excellent user experience.

The strategy prioritizes immediate feedback for social interactions while maintaining performance and resource efficiency for content consumption patterns.

### Key Benefits
- **Optimal Resource Usage**: WebSocket bandwidth used only where it adds real value
- **Scalable Architecture**: Can handle growth without exponential resource increases
- **Excellent UX**: Real-time where users expect it, efficient elsewhere
- **Maintainable Code**: Clear separation of concerns between real-time and polling features
- **Cost Effective**: Lower infrastructure costs compared to full real-time approach
