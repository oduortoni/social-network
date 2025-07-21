# Frontend WebSocket Implementation Documentation

## Overview

This document provides comprehensive documentation for the frontend WebSocket implementation in the social network application. The frontend WebSocket system enables real-time messaging, notifications, and live updates in the React/Next.js application.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Components](#core-components)
3. [WebSocket Service](#websocket-service)
4. [React Components](#react-components)
5. [API Integration](#api-integration)
6. [State Management](#state-management)
7. [Usage Examples](#usage-examples)
8. [Best Practices](#best-practices)

## Architecture Overview

The frontend WebSocket implementation follows a service-oriented architecture with React hooks for state management and real-time updates.

```
React Components
    ↓ useEffect hooks
WebSocket Service (wsService)
    ↓ WebSocket connection
Backend WebSocket Server
    ↓ Real-time updates
React State Updates
    ↓ UI re-rendering
User Interface
```

### Key Design Principles

- **Service Layer Pattern**: Centralized WebSocket management through `wsService`
- **React Hooks Integration**: Seamless integration with React component lifecycle
- **Connection State Management**: Robust connection handling with visual feedback
- **Message Type Routing**: Organized message handling by type
- **Error Recovery**: Automatic reconnection with exponential backoff

## Core Components

### 1. WebSocket Service (`lib/websocket.js`)

The central service that manages all WebSocket communication:

```javascript
class WebSocketService {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.messageHandlers = new Map();
    this.isConnecting = false;
    this.shouldReconnect = true;
  }
}
```

**Key Features:**
- Connection lifecycle management
- Automatic reconnection with exponential backoff
- Message type routing and handler registration
- Connection state tracking and events
- Error handling and recovery

### 2. API Service (`lib/api.js`)

HTTP API client for chat history and persistent data:

```javascript
export const chatAPI = {
  getPrivateMessages: (userId, limit, offset) => { /* ... */ },
  getGroupMessages: (groupId, limit, offset) => { /* ... */ },
  sendGroupInvite: (groupId, userId, groupName) => { /* ... */ },
  getNotifications: (limit, offset) => { /* ... */ },
  markNotificationsRead: () => { /* ... */ }
};
```

### 3. React Components

#### ChatInterface (`components/chat/ChatInterface.jsx`)
Real-time messaging interface with WebSocket integration.

#### NotificationCenter (`components/notifications/NotificationCenter.jsx`)
Real-time notification system with dropdown UI.

## WebSocket Service

### Connection Management

```javascript
// Connect to WebSocket
wsService.connect();

// Check connection status
if (wsService.isConnected()) {
  // Connection is active
}

// Disconnect
wsService.disconnect();
```

### Message Handling

```javascript
// Register message handlers
wsService.onMessage('private', (message) => {
  console.log('Private message received:', message);
});

wsService.onMessage('notification', (notification) => {
  console.log('Notification received:', notification);
});

// Send messages
wsService.sendMessage('private', 'Hello!', recipientId);
wsService.sendMessage('group', 'Hello group!', null, groupId);
wsService.sendMessage('broadcast', 'System announcement');
```

### Connection Status Events

```javascript
wsService.onMessage('connection_status', (statusEvent) => {
  const status = statusEvent.status; // 'connecting', 'connected', 'disconnected'
  updateUIConnectionStatus(status);
});
```

## React Components

### ChatInterface Integration

```javascript
import { wsService } from '../../lib/websocket';
import { chatAPI } from '../../lib/api';

const ChatInterface = ({ user }) => {
  const [messages, setMessages] = useState([]);
  const [connectionStatus, setConnectionStatus] = useState('disconnected');

  useEffect(() => {
    let mounted = true;
    
    // Set up message handlers
    wsService.onMessage('private', handlePrivateMessage);
    wsService.onMessage('group', handleGroupMessage);
    wsService.onMessage('broadcast', handleBroadcastMessage);
    wsService.onMessage('notification', handleNotification);
    
    // Connection status tracking
    wsService.onMessage('connection_status', (message) => {
      if (mounted) {
        setConnectionStatus(message.status);
      }
    });
    
    // Delayed connection for server readiness
    const connectTimer = setTimeout(() => {
      if (mounted) {
        setConnectionStatus('connecting');
        wsService.connect();
      }
    }, 500);
    
    return () => {
      mounted = false;
      clearTimeout(connectTimer);
      wsService.disconnect();
    };
  }, []);

  const sendMessage = () => {
    if (!wsService.isConnected()) {
      alert('Not connected to chat server');
      return;
    }
    
    wsService.sendMessage('private', messageContent, recipientId);
  };
};
```

### NotificationCenter Integration

```javascript
const NotificationCenter = () => {
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    // Load existing notifications
    loadNotifications();
    
    // Listen for real-time notifications
    wsService.onMessage('notification', (notification) => {
      setNotifications(prev => [notification, ...prev]);
      setUnreadCount(prev => prev + 1);
    });
  }, []);

  const loadNotifications = async () => {
    try {
      const data = await chatAPI.getNotifications();
      setNotifications(data || []);
      setUnreadCount(data?.filter(n => !n.is_read).length || 0);
    } catch (error) {
      console.error('Failed to load notifications:', error);
    }
  };
};
```

## API Integration

### HTTP + WebSocket Pattern

The frontend uses a hybrid approach combining HTTP APIs for persistent data and WebSocket for real-time updates:

```javascript
// Load chat history via HTTP
const loadChatHistory = async (chatType, chatId) => {
  try {
    let history;
    if (chatType === 'private') {
      history = await chatAPI.getPrivateMessages(chatId);
    } else if (chatType === 'group') {
      history = await chatAPI.getGroupMessages(chatId);
    }
    setMessages(history || []);
  } catch (error) {
    console.error('Failed to load chat history:', error);
  }
};

// Real-time message updates via WebSocket
wsService.onMessage('private', (message) => {
  setMessages(prev => [...prev, message]);
});
```

### Environment Configuration

```javascript
// WebSocket URL derivation
const wsUrl = process.env.NEXT_PUBLIC_API_URL?.replace('http', 'ws') + '/ws';

// API base URL
const API_BASE = process.env.NEXT_PUBLIC_API_URL;
```

## State Management

### Connection State

```javascript
const [connectionStatus, setConnectionStatus] = useState('disconnected');

// Status values: 'connecting', 'connected', 'disconnected'
```

### Message State

```javascript
const [messages, setMessages] = useState([]);
const [activeChat, setActiveChat] = useState(null);

// activeChat format: { type: 'private', id: userId } or { type: 'group', id: groupId }
```

### Notification State

```javascript
const [notifications, setNotifications] = useState([]);
const [unreadCount, setUnreadCount] = useState(0);
```

## Usage Examples

### Basic WebSocket Integration

```javascript
import { wsService } from '../lib/websocket';

const MyComponent = () => {
  useEffect(() => {
    // Connect and set up handlers
    wsService.connect();
    wsService.onMessage('custom_event', handleCustomEvent);
    
    return () => wsService.disconnect();
  }, []);

  const handleCustomEvent = (message) => {
    // Handle custom message type
  };
};
```

### Sending Different Message Types

```javascript
// Private message
wsService.sendMessage('private', 'Hello!', 123);

// Group message
wsService.sendMessage('group', 'Hello group!', null, '456');

// Broadcast message (admin only)
wsService.sendMessage('broadcast', 'System announcement');
```

### Connection Status UI

```javascript
const ConnectionIndicator = () => {
  const [status, setStatus] = useState('disconnected');

  useEffect(() => {
    wsService.onMessage('connection_status', (event) => {
      setStatus(event.status);
    });
  }, []);

  return (
    <div className={`w-3 h-3 rounded-full ${
      status === 'connected' ? 'bg-green-500' : 
      status === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
    }`} />
  );
};
```

## Best Practices

### 1. Component Lifecycle Management

```javascript
useEffect(() => {
  let mounted = true;
  
  // Set up handlers
  wsService.onMessage('event', (data) => {
    if (mounted) {
      setState(data);
    }
  });
  
  return () => {
    mounted = false;
    // Cleanup if needed
  };
}, []);
```

### 2. Error Handling

```javascript
const sendMessageSafely = (type, content, target) => {
  if (!wsService.isConnected()) {
    showErrorMessage('Not connected to server');
    return;
  }
  
  try {
    wsService.sendMessage(type, content, target);
  } catch (error) {
    console.error('Failed to send message:', error);
    showErrorMessage('Failed to send message');
  }
};
```

### 3. Connection State Feedback

```javascript
// Always provide visual feedback for connection state
const MessageInput = () => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  
  return (
    <input
      disabled={connectionStatus !== 'connected'}
      placeholder={
        connectionStatus === 'connected' 
          ? "Type a message..." 
          : "Waiting for connection..."
      }
      className={connectionStatus !== 'connected' ? 'bg-gray-100' : ''}
    />
  );
};
```

### 4. Memory Management

```javascript
// Clear timeouts and intervals
useEffect(() => {
  const timer = setTimeout(() => wsService.connect(), 500);
  
  return () => {
    clearTimeout(timer);
    wsService.disconnect();
  };
}, []);
```

This frontend WebSocket implementation provides a robust foundation for real-time features while maintaining clean separation of concerns and excellent user experience.
