# Activity Sidebar Implementation Guide

## Overview

The Activity Sidebar is a notification-aware component system that displays actionable user activities with intelligent routing and endpoint integration. Each activity type has its own specialized component that knows how to handle its specific data and actions.

## Architecture

### Core Principles
- **Notification-Aware Components**: Each activity component subscribes to relevant WebSocket notifications
- **Smart Routing**: Components know which routes to navigate to for different actions
- **Endpoint Integration**: Components have built-in knowledge of API endpoints they need to call
- **Self-Contained Actions**: Each component handles its own action logic and state management

### Component Hierarchy
```
ActivitySidebar
├── ActivityFeed
│   ├── FollowRequestActivity
│   ├── NewFollowerActivity
│   ├── PostEngagementActivity
│   ├── GroupInviteActivity
│   ├── UserConnectionActivity
│   └── GenericActivity
└── ActivityFilters
```

## Activity Component Specifications

### Base Activity Component

All activity components extend from a base `ActivityItem` component:

```javascript
// Base props structure
interface ActivityItemProps {
  activity: {
    id: string;
    type: string;
    subtype: string;
    actor: UserData;
    target?: TargetData;
    message: string;
    timestamp: number;
    is_read: boolean;
    requires_action: boolean;
    actions: string[];
    additional_data: object;
  };
  onAction: (activityId: string, action: string, data?: object) => void;
  onMarkRead: (activityId: string) => void;
}
```

### 1. FollowRequestActivity Component

**Purpose**: Handle incoming and outgoing follow requests

**WebSocket Subscriptions**:
- `notification` type with `subtype: "follow_request"`
- `notification` type with `subtype: "follow_request_accepted"`
- `notification` type with `subtype: "follow_request_declined"`

**API Endpoints**:
- `POST /api/follow/accept` - Accept follow request
- `POST /api/follow/decline` - Decline follow request
- `GET /api/users/{id}/profile` - View user profile

**Routes**:
- `/profile/{user_id}` - View requester's profile
- `/followers` - View followers list

**Actions**:
```javascript
const actions = {
  accept: async (requestId) => {
    await fetch('/api/follow/accept', {
      method: 'POST',
      body: JSON.stringify({ request_id: requestId })
    });
  },
  decline: async (requestId) => {
    await fetch('/api/follow/decline', {
      method: 'POST', 
      body: JSON.stringify({ request_id: requestId })
    });
  },
  viewProfile: (userId) => {
    router.push(`/profile/${userId}`);
  }
};
```

**Data Requirements**:
```javascript
{
  type: "follow_request",
  additional_data: {
    request_id: "uuid",
    requester_id: 123
  }
}
```

### 2. NewFollowerActivity Component

**Purpose**: Display new followers and provide follow-back options

**WebSocket Subscriptions**:
- `notification` type with `subtype: "new_follower"`

**API Endpoints**:
- `POST /api/follow/request` - Send follow request back
- `GET /api/users/{id}/profile` - View follower profile

**Routes**:
- `/profile/{user_id}` - View follower's profile
- `/followers` - View all followers

**Actions**:
```javascript
const actions = {
  followBack: async (userId) => {
    await fetch('/api/follow/request', {
      method: 'POST',
      body: JSON.stringify({ target_user_id: userId })
    });
  },
  viewProfile: (userId) => {
    router.push(`/profile/${userId}`);
  },
  viewFollowers: () => {
    router.push('/followers');
  }
};
```

### 3. PostEngagementActivity Component

**Purpose**: Handle likes, comments, and mentions on posts

**WebSocket Subscriptions**:
- `notification` type with `subtype: "post_like"`
- `notification` type with `subtype: "post_comment"`
- `notification` type with `subtype: "post_mention"`

**API Endpoints**:
- `GET /api/posts/{id}` - View specific post
- `POST /api/posts/{id}/like` - Like a post
- `POST /api/posts/{id}/comments` - Add comment

**Routes**:
- `/posts/{post_id}` - View full post
- `/profile/{user_id}` - View actor's profile

**Actions**:
```javascript
const actions = {
  viewPost: (postId) => {
    router.push(`/posts/${postId}`);
  },
  likeBack: async (postId) => {
    await fetch(`/api/posts/${postId}/like`, {
      method: 'POST'
    });
  },
  reply: (postId) => {
    router.push(`/posts/${postId}#comment`);
  }
};
```

### 4. GroupInviteActivity Component

**Purpose**: Handle group invitations and join requests

**WebSocket Subscriptions**:
- `notification` type with `subtype: "group_invite"`
- `notification` type with `subtype: "group_join_request"`

**API Endpoints**:
- `POST /api/groups/{id}/accept-invite` - Accept group invitation
- `POST /api/groups/{id}/decline-invite` - Decline group invitation
- `GET /api/groups/{id}` - View group details

**Routes**:
- `/groups/{group_id}` - View group page
- `/groups` - View all groups

**Actions**:
```javascript
const actions = {
  acceptInvite: async (groupId, inviteId) => {
    await fetch(`/api/groups/${groupId}/accept-invite`, {
      method: 'POST',
      body: JSON.stringify({ invite_id: inviteId })
    });
  },
  declineInvite: async (groupId, inviteId) => {
    await fetch(`/api/groups/${groupId}/decline-invite`, {
      method: 'POST',
      body: JSON.stringify({ invite_id: inviteId })
    });
  },
  viewGroup: (groupId) => {
    router.push(`/groups/${groupId}`);
  }
};
```

### 5. UserConnectionActivity Component

**Purpose**: Display user online/offline status and provide messaging options

**WebSocket Subscriptions**:
- `notification` type with `subtype: "user_connected"`
- `notification` type with `subtype: "user_disconnected"`

**API Endpoints**:
- `POST /api/messages/private` - Send private message
- `GET /api/users/{id}/profile` - View user profile

**Routes**:
- `/chat/{user_id}` - Open private chat
- `/profile/{user_id}` - View user profile

**Actions**:
```javascript
const actions = {
  sendMessage: (userId) => {
    router.push(`/chat/${userId}`);
  },
  viewProfile: (userId) => {
    router.push(`/profile/${userId}`);
  }
};
```

## Activity Feed Management

### ActivityFeed Component

**Responsibilities**:
- Fetch initial activity data from API
- Subscribe to real-time activity updates via WebSocket
- Manage activity state (read/unread, filtering)
- Handle batch operations (mark all as read)

**API Integration**:
```javascript
// Fetch activities
GET /api/activities?limit=20&offset=0&filter=unread

// Mark activities as read
POST /api/activities/read
Body: { activity_ids: ["uuid1", "uuid2"] }

// Get activity count
GET /api/activities/count?filter=unread
```

**WebSocket Integration**:
```javascript
// Listen for new activities
wsService.onMessage('activity', (activity) => {
  setActivities(prev => [activity, ...prev]);
  updateUnreadCount(prev => prev + 1);
});

// Listen for activity updates
wsService.onMessage('activity_update', (update) => {
  updateActivityStatus(update.activity_id, update.status);
});
```

### Activity Filtering

**Filter Types**:
- `all` - Show all activities
- `unread` - Show only unread activities
- `action_required` - Show only activities requiring user action
- `social` - Show only social interactions (follows, likes)
- `content` - Show only content-related activities
- `groups` - Show only group-related activities

**Implementation**:
```javascript
const filterActivities = (activities, filter) => {
  switch (filter) {
    case 'unread':
      return activities.filter(a => !a.is_read);
    case 'action_required':
      return activities.filter(a => a.requires_action);
    case 'social':
      return activities.filter(a => 
        ['follow_request', 'new_follower'].includes(a.type)
      );
    // ... other filters
  }
};
```

## State Management

### Activity State Structure
```javascript
const activityState = {
  activities: [],
  unreadCount: 0,
  filter: 'all',
  loading: false,
  error: null,
  hasMore: true,
  offset: 0
};
```

### Action Handlers
```javascript
const handleActivityAction = async (activityId, action, data) => {
  try {
    // Find the activity component
    const activity = activities.find(a => a.id === activityId);
    const component = getActivityComponent(activity.type);
    
    // Execute the action
    await component.actions[action](data);
    
    // Update activity state
    markActivityAsActioned(activityId, action);
    
    // Show success feedback
    showToast(`Action ${action} completed successfully`);
  } catch (error) {
    showToast(`Failed to ${action}: ${error.message}`, 'error');
  }
};
```

## Error Handling

### Network Errors
- Retry failed API calls with exponential backoff
- Show offline indicators when network is unavailable
- Queue actions for retry when connection is restored

### Action Failures
- Display specific error messages for different failure types
- Provide retry options for failed actions
- Maintain activity state consistency

### WebSocket Disconnections
- Gracefully handle WebSocket disconnections
- Fetch missed activities on reconnection
- Show connection status indicators

## Performance Considerations

### Lazy Loading
- Load activities in batches of 20
- Implement infinite scroll for older activities
- Cache activity data in memory with TTL

### Real-Time Updates
- Debounce rapid activity updates
- Batch multiple activities arriving simultaneously
- Limit maximum activities in memory (e.g., 100)

### Component Optimization
- Use React.memo for activity components
- Implement virtual scrolling for large activity lists
- Optimize re-renders with proper key props

## Testing Strategy

### Unit Tests
- Test each activity component's action handlers
- Verify API endpoint calls with correct data
- Test routing behavior for different actions

### Integration Tests
- Test WebSocket subscription and message handling
- Verify activity state management
- Test error handling scenarios

### E2E Tests
- Test complete user flows (receive notification → take action)
- Verify cross-component interactions
- Test real-time updates across multiple browser tabs

## Implementation Phases

### Phase 1: Core Infrastructure
1. Base ActivityItem component
2. ActivityFeed container
3. Basic API integration
4. WebSocket subscription system

### Phase 2: Essential Activities
1. FollowRequestActivity
2. NewFollowerActivity
3. UserConnectionActivity

### Phase 3: Content Activities
1. PostEngagementActivity
2. GroupInviteActivity

### Phase 4: Advanced Features
1. Activity filtering
2. Batch operations
3. Performance optimizations
4. Advanced error handling

This documentation provides a comprehensive guide for implementing a robust, scalable Activity Sidebar system that integrates seamlessly with the existing notification infrastructure while providing rich, actionable user experiences.
