# Social Network Documentation

## Overview

This documentation covers the complete social network backend implementation, with a focus on the real-time messaging and WebSocket system.

## WebSocket Implementation

The WebSocket system provides real-time messaging, notifications, and live updates throughout the application.

### Documentation Index

#### Backend WebSocket Documentation
| Document | Description | Audience |
|----------|-------------|----------|
| **[WebSocket Implementation](websocket-implementation.md)** | Complete overview of the backend WebSocket system architecture | Backend developers |
| **[Backend Integration Guide](websocket-integration-guide.md)** | How to add real-time features to backend components | Backend developers |
| **[API Reference](websocket-api-reference.md)** | Complete API documentation and message formats | All developers |
| **[Backend Testing Guide](websocket-testing-guide.md)** | Backend testing strategies and examples | Backend QA & developers |
| **[Extension Guide](websocket-extension-guide.md)** | How to extend the backend WebSocket system | Senior backend developers |

#### Frontend WebSocket Documentation
| Document | Description | Audience |
|----------|-------------|----------|
| **[Frontend Implementation](frontend-websocket-implementation.md)** | Complete overview of the frontend WebSocket integration | Frontend developers |
| **[Frontend Integration Guide](frontend-websocket-integration-guide.md)** | How to add real-time features to React components | Frontend developers |
| **[Frontend Testing Guide](frontend-websocket-testing-guide.md)** | Frontend testing strategies and React component testing | Frontend QA & developers |

### Quick Start

#### For Frontend Developers
1. **Start with**: [Frontend Implementation](frontend-websocket-implementation.md) to understand React integration
2. **Integration**: [Frontend Integration Guide](frontend-websocket-integration-guide.md) to add real-time features
3. **Testing**: [Frontend Testing Guide](frontend-websocket-testing-guide.md) for React component testing
4. **API Reference**: [API Reference](websocket-api-reference.md) for message formats

#### For Backend Developers
1. **Start with**: [WebSocket Implementation](websocket-implementation.md) for system architecture
2. **Integration**: [Backend Integration Guide](websocket-integration-guide.md) to add real-time features
3. **Testing**: [Backend Testing Guide](websocket-testing-guide.md) for Go testing strategies
4. **Extension**: [Extension Guide](websocket-extension-guide.md) for advanced features

### Common Use Cases

#### Backend: Adding Real-time Notifications
```go
// Quick example - see Backend Integration Guide for details
notifier := ws.NewDBNotificationSender(wsManager)
notifier.SendNotification(userID, map[string]interface{}{
    "type":      "notification",
    "subtype":   "friend_request",
    "message":   "You have a new friend request",
    "timestamp": time.Now().Unix(),
})
```

#### Backend: Checking User Online Status
```go
if wsManager.IsOnline(userID) {
    // Send real-time notification
} else {
    // Store for later retrieval
}
```

#### Frontend: React Component Integration
```javascript
// Quick example - see Frontend Integration Guide for details
import { wsService } from '../lib/websocket';

const MyComponent = () => {
  useEffect(() => {
    wsService.onMessage('notification', handleNotification);
    wsService.connect();
    return () => wsService.disconnect();
  }, []);

  const sendMessage = () => {
    wsService.sendMessage('private', 'Hello!', recipientId);
  };
};
```

### Implementation Status

| Feature | Backend Status | Frontend Status | Documentation |
|---------|----------------|-----------------|---------------|
| **Real-time Messaging** | ✅ Complete | ✅ Complete | [API Reference](websocket-api-reference.md) |
| **Notifications** | ✅ Complete | ✅ Complete | [Backend](websocket-integration-guide.md) \| [Frontend](frontend-websocket-integration-guide.md) |
| **Group Chat** | ✅ Complete | ✅ Complete | [Backend](websocket-implementation.md) \| [Frontend](frontend-websocket-implementation.md) |
| **HTTP API** | ✅ Complete | ✅ Complete | [API Reference](websocket-api-reference.md) |
| **Testing Suite** | ✅ Complete | ✅ Complete | [Backend](websocket-testing-guide.md) \| [Frontend](frontend-websocket-testing-guide.md) |
| **Friend System** | 🔄 In Progress | 🔄 In Progress | [Backend](websocket-integration-guide.md#friend-system-integration) \| [Frontend](frontend-websocket-integration-guide.md#friend-system-integration) |
| **Posts/Feed** | ⏳ Planned | ⏳ Planned | [Backend](websocket-integration-guide.md#posts-system-integration) \| [Frontend](frontend-websocket-integration-guide.md#posts-feed-system-integration) |
| **Events** | ⏳ Planned | ⏳ Planned | [Backend](websocket-integration-guide.md#group-events-integration) \| [Frontend](frontend-websocket-integration-guide.md#group-events-integration) |

### Architecture Overview

```
Frontend (React/Next.js)
    ↓ WebSocket Connection
Backend WebSocket Manager
    ↓ Message Routing
┌─────────────────┬─────────────────┬─────────────────┐
│   Private Chat  │   Group Chat    │   Notifications │
└─────────────────┴─────────────────┴─────────────────┘
    ↓ Persistence
Database (SQLite)
```

### Related Systems

- **Authentication**: Session-based authentication via cookies
- **Database**: SQLite with comprehensive schema for social features
- **HTTP API**: RESTful endpoints for data retrieval and management
- **Frontend**: Next.js application with real-time WebSocket integration

### Support

For questions about the WebSocket implementation:

1. **Check the documentation** - Most questions are answered in the guides above
2. **Review the tests** - See [Testing Guide](websocket-testing-guide.md) for examples
3. **Look at integration examples** - [Integration Guide](websocket-integration-guide.md) has real-world patterns

### Development Workflow

#### Backend Development
1. **Planning**: Review [WebSocket Implementation](websocket-implementation.md) for architecture
2. **Development**: Follow [Backend Integration Guide](websocket-integration-guide.md) for implementation
3. **Testing**: Use [Backend Testing Guide](websocket-testing-guide.md) for validation
4. **Extension**: Reference [Extension Guide](websocket-extension-guide.md) for advanced features

#### Frontend Development
1. **Planning**: Review [Frontend Implementation](frontend-websocket-implementation.md) for React architecture
2. **Development**: Follow [Frontend Integration Guide](frontend-websocket-integration-guide.md) for React components
3. **Testing**: Use [Frontend Testing Guide](frontend-websocket-testing-guide.md) for component testing
4. **API Integration**: Reference [API Reference](websocket-api-reference.md) for message formats
