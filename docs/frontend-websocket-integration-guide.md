# Frontend WebSocket Integration Guide

## For Frontend Developers Building Real-time Features

This guide shows how to leverage the existing frontend WebSocket implementation to add real-time functionality to new React components and features.

## Quick Start Integration

### 1. Basic Component Setup

```javascript
import { wsService } from '../../lib/websocket';
import { chatAPI } from '../../lib/api';

const YourComponent = ({ user }) => {
  const [realTimeData, setRealTimeData] = useState([]);
  const [connectionStatus, setConnectionStatus] = useState('disconnected');

  useEffect(() => {
    let mounted = true;
    
    // Set up real-time message handlers
    wsService.onMessage('your_message_type', handleYourMessage);
    wsService.onMessage('connection_status', (event) => {
      if (mounted) setConnectionStatus(event.status);
    });
    
    // Connect with delay for server readiness
    const connectTimer = setTimeout(() => {
      if (mounted) wsService.connect();
    }, 500);
    
    return () => {
      mounted = false;
      clearTimeout(connectTimer);
    };
  }, []);

  const handleYourMessage = (message) => {
    setRealTimeData(prev => [...prev, message]);
  };
};
```

### 2. Sending Real-time Updates

```javascript
const sendUpdate = (data) => {
  if (!wsService.isConnected()) {
    alert('Not connected to server');
    return;
  }
  
  // Send via WebSocket for real-time delivery
  wsService.sendMessage('your_message_type', JSON.stringify(data));
};
```

## Feature-Specific Integration Examples

### Friend System Integration

```javascript
const FriendRequestComponent = ({ user }) => {
  const [friendRequests, setFriendRequests] = useState([]);
  const [friends, setFriends] = useState([]);

  useEffect(() => {
    // Load existing data
    loadFriendRequests();
    loadFriends();
    
    // Set up real-time handlers
    wsService.onMessage('notification', handleFriendNotification);
    wsService.onMessage('friend_update', handleFriendUpdate);
  }, []);

  const handleFriendNotification = (notification) => {
    if (notification.subtype === 'friend_request') {
      // Add new friend request to state
      setFriendRequests(prev => [...prev, {
        id: notification.requester_id,
        name: notification.requester_name,
        timestamp: notification.timestamp
      }]);
      
      // Show toast notification
      showToast(`${notification.requester_name} sent you a friend request`);
    }
  };

  const sendFriendRequest = async (targetUserId) => {
    try {
      // Send HTTP request
      await fetch('/api/friends/request', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ user_id: targetUserId })
      });
      
      // Real-time notification is handled by backend
      showToast('Friend request sent!');
    } catch (error) {
      console.error('Failed to send friend request:', error);
    }
  };

  const acceptFriendRequest = async (requestId) => {
    try {
      await fetch(`/api/friends/accept/${requestId}`, {
        method: 'POST',
        credentials: 'include'
      });
      
      // Remove from pending requests
      setFriendRequests(prev => prev.filter(req => req.id !== requestId));
      
      // Real-time update will be sent to both users by backend
    } catch (error) {
      console.error('Failed to accept friend request:', error);
    }
  };
};
```

### Posts/Feed System Integration

```javascript
const PostFeedComponent = ({ user }) => {
  const [posts, setPosts] = useState([]);
  const [newPostContent, setNewPostContent] = useState('');

  useEffect(() => {
    loadFeed();
    
    // Real-time post updates
    wsService.onMessage('notification', handlePostNotification);
    wsService.onMessage('post_update', handlePostUpdate);
  }, []);

  const handlePostNotification = (notification) => {
    switch (notification.subtype) {
      case 'new_post':
        // Friend posted something new
        showToast(`${notification.user_name} shared a new post`);
        // Optionally refresh feed or add post to top
        break;
        
      case 'post_like':
        // Someone liked your post
        showToast(`${notification.liker_name} liked your post`);
        updatePostLikes(notification.post_id, notification.like_count);
        break;
        
      case 'post_comment':
        // Someone commented on your post
        showToast(`${notification.commenter_name} commented on your post`);
        updatePostComments(notification.post_id);
        break;
    }
  };

  const createPost = async () => {
    if (!newPostContent.trim()) return;
    
    try {
      const response = await fetch('/api/posts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ content: newPostContent })
      });
      
      const newPost = await response.json();
      
      // Add to local state immediately
      setPosts(prev => [newPost, ...prev]);
      setNewPostContent('');
      
      // Backend will send real-time notifications to followers
    } catch (error) {
      console.error('Failed to create post:', error);
    }
  };

  const likePost = async (postId) => {
    try {
      await fetch(`/api/posts/${postId}/like`, {
        method: 'POST',
        credentials: 'include'
      });
      
      // Update local state optimistically
      setPosts(prev => prev.map(post => 
        post.id === postId 
          ? { ...post, likes: post.likes + 1, liked_by_user: true }
          : post
      ));
      
      // Backend sends real-time notification to post owner
    } catch (error) {
      console.error('Failed to like post:', error);
    }
  };
};
```

### Group Events Integration

```javascript
const GroupEventsComponent = ({ groupId, user }) => {
  const [events, setEvents] = useState([]);
  const [rsvps, setRsvps] = useState({});

  useEffect(() => {
    loadGroupEvents();
    
    // Real-time event updates
    wsService.onMessage('notification', handleEventNotification);
    wsService.onMessage('event_update', handleEventUpdate);
  }, []);

  const handleEventNotification = (notification) => {
    switch (notification.subtype) {
      case 'new_event':
        // New event in group
        showToast(`New event: ${notification.event_title}`);
        loadGroupEvents(); // Refresh events
        break;
        
      case 'event_rsvp':
        // Someone RSVP'd to your event
        showToast(`${notification.user_name} RSVP'd to ${notification.event_title}`);
        updateEventRSVP(notification.event_id, notification.response);
        break;
    }
  };

  const createEvent = async (eventData) => {
    try {
      const response = await fetch(`/api/groups/${groupId}/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(eventData)
      });
      
      const newEvent = await response.json();
      setEvents(prev => [newEvent, ...prev]);
      
      // Backend sends notifications to all group members
    } catch (error) {
      console.error('Failed to create event:', error);
    }
  };

  const rsvpToEvent = async (eventId, response) => {
    try {
      await fetch(`/api/events/${eventId}/rsvp`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ response })
      });
      
      // Update local RSVP state
      setRsvps(prev => ({ ...prev, [eventId]: response }));
      
      // Backend sends notification to event creator
    } catch (error) {
      console.error('Failed to RSVP:', error);
    }
  };
};
```

## Custom Hook Patterns

### useWebSocketConnection Hook

```javascript
import { useState, useEffect } from 'react';
import { wsService } from '../lib/websocket';

export const useWebSocketConnection = () => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const handleConnectionStatus = (event) => {
      setConnectionStatus(event.status);
      setIsConnected(event.status === 'connected');
    };

    wsService.onMessage('connection_status', handleConnectionStatus);
    
    // Initial connection
    const timer = setTimeout(() => wsService.connect(), 500);
    
    return () => {
      clearTimeout(timer);
    };
  }, []);

  return { connectionStatus, isConnected };
};
```

### useRealTimeNotifications Hook

```javascript
export const useRealTimeNotifications = () => {
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    const handleNotification = (notification) => {
      setNotifications(prev => [notification, ...prev]);
      setUnreadCount(prev => prev + 1);
      
      // Show toast notification
      showToast(notification.message);
    };

    wsService.onMessage('notification', handleNotification);
  }, []);

  const markAsRead = async () => {
    try {
      await chatAPI.markNotificationsRead();
      setUnreadCount(0);
    } catch (error) {
      console.error('Failed to mark notifications as read:', error);
    }
  };

  return { notifications, unreadCount, markAsRead };
};
```

### useRealTimeData Hook

```javascript
export const useRealTimeData = (messageType, initialData = []) => {
  const [data, setData] = useState(initialData);

  useEffect(() => {
    const handleMessage = (message) => {
      setData(prev => [...prev, message]);
    };

    wsService.onMessage(messageType, handleMessage);
  }, [messageType]);

  const sendData = (content, target = null) => {
    if (wsService.isConnected()) {
      wsService.sendMessage(messageType, content, target);
    }
  };

  return { data, setData, sendData };
};
```

## UI Components for Real-time Features

### Connection Status Component

```javascript
const ConnectionStatus = () => {
  const { connectionStatus } = useWebSocketConnection();

  return (
    <div className="flex items-center space-x-2">
      <div className={`w-2 h-2 rounded-full ${
        connectionStatus === 'connected' ? 'bg-green-500' : 
        connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
      }`} />
      <span className="text-sm text-gray-600">
        {connectionStatus === 'connected' ? 'Online' : 
         connectionStatus === 'connecting' ? 'Connecting...' : 'Offline'}
      </span>
    </div>
  );
};
```

### Toast Notification Component

```javascript
const ToastNotification = ({ message, type = 'info', onClose }) => {
  useEffect(() => {
    const timer = setTimeout(onClose, 5000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div className={`fixed top-4 right-4 p-4 rounded-lg shadow-lg ${
      type === 'success' ? 'bg-green-500' : 
      type === 'error' ? 'bg-red-500' : 'bg-blue-500'
    } text-white`}>
      <div className="flex items-center justify-between">
        <span>{message}</span>
        <button onClick={onClose} className="ml-4 text-white">×</button>
      </div>
    </div>
  );
};
```

## Testing Real-time Features

### Mock WebSocket Service

```javascript
// For testing
export const mockWSService = {
  isConnected: () => true,
  connect: jest.fn(),
  disconnect: jest.fn(),
  sendMessage: jest.fn(),
  onMessage: jest.fn(),
  simulateMessage: (type, data) => {
    // Simulate receiving a message
    const handlers = mockWSService.messageHandlers.get(type);
    if (handlers) handlers.forEach(handler => handler(data));
  },
  messageHandlers: new Map()
};
```

### Component Testing

```javascript
import { render, screen } from '@testing-library/react';
import { mockWSService } from '../__mocks__/websocket';

jest.mock('../lib/websocket', () => ({
  wsService: mockWSService
}));

test('handles real-time notifications', () => {
  render(<YourComponent />);
  
  // Simulate receiving a notification
  mockWSService.simulateMessage('notification', {
    type: 'notification',
    subtype: 'friend_request',
    message: 'You have a new friend request'
  });
  
  expect(screen.getByText('You have a new friend request')).toBeInTheDocument();
});
```

This integration guide provides everything you need to add real-time functionality to your React components using the existing WebSocket infrastructure.
