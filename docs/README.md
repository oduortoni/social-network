# Social Network Documentation

## Overview

This documentation covers the complete social network backend implementation, with a focus on the real-time messaging and WebSocket system.

## WebSocket Implementation

The WebSocket system provides real-time messaging, notifications, and live updates throughout the application.

### 📚 Documentation Index

| Document | Description | Audience |
|----------|-------------|----------|
| **[WebSocket Implementation](websocket-implementation.md)** | Complete overview of the WebSocket system architecture | All developers |
| **[Integration Guide](websocket-integration-guide.md)** | How to add real-time features to your components | Feature developers |
| **[API Reference](websocket-api-reference.md)** | Complete API documentation and message formats | Frontend & API developers |
| **[Testing Guide](websocket-testing-guide.md)** | Testing strategies and examples | QA & developers |
| **[Extension Guide](websocket-extension-guide.md)** | How to extend the WebSocket system | Senior developers |

### 🚀 Quick Start

1. **For Frontend Developers**: Start with [API Reference](websocket-api-reference.md) to understand message formats
2. **For Backend Developers**: Read [Integration Guide](websocket-integration-guide.md) to add real-time features
3. **For Testing**: Follow [Testing Guide](websocket-testing-guide.md) for comprehensive testing strategies
4. **For Architecture**: Review [WebSocket Implementation](websocket-implementation.md) for system overview

### 🔧 Common Use Cases

#### Adding Real-time Notifications
```go
// Quick example - see Integration Guide for details
notifier := ws.NewDBNotificationSender(wsManager)
notifier.SendNotification(userID, map[string]interface{}{
    "type":      "notification",
    "subtype":   "friend_request",
    "message":   "You have a new friend request",
    "timestamp": time.Now().Unix(),
})
```

#### Checking User Online Status
```go
if wsManager.IsOnline(userID) {
    // Send real-time notification
} else {
    // Store for later retrieval
}
```

#### Broadcasting to Groups
```go
wsManager.BroadcastToGroup(senderID, groupID, messageData)
```

### 📋 Implementation Status

| Feature | Status | Documentation |
|---------|--------|---------------|
| **Real-time Messaging** | ✅ Complete | [API Reference](websocket-api-reference.md) |
| **Notifications** | ✅ Complete | [Integration Guide](websocket-integration-guide.md) |
| **Group Chat** | ✅ Complete | [WebSocket Implementation](websocket-implementation.md) |
| **HTTP API** | ✅ Complete | [API Reference](websocket-api-reference.md) |
| **Testing Suite** | ✅ Complete | [Testing Guide](websocket-testing-guide.md) |
| **Friend System** | 🔄 In Progress | [Integration Guide](websocket-integration-guide.md#friend-system-integration) |
| **Posts/Feed** | ⏳ Planned | [Integration Guide](websocket-integration-guide.md#posts-system-integration) |
| **Events** | ⏳ Planned | [Integration Guide](websocket-integration-guide.md#group-events-integration) |

### 🏗️ Architecture Overview

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

### 🔗 Related Systems

- **Authentication**: Session-based authentication via cookies
- **Database**: SQLite with comprehensive schema for social features
- **HTTP API**: RESTful endpoints for data retrieval and management
- **Frontend**: Next.js application with real-time WebSocket integration

### 📞 Support

For questions about the WebSocket implementation:

1. **Check the documentation** - Most questions are answered in the guides above
2. **Review the tests** - See [Testing Guide](websocket-testing-guide.md) for examples
3. **Look at integration examples** - [Integration Guide](websocket-integration-guide.md) has real-world patterns

### 🔄 Development Workflow

1. **Planning**: Review [WebSocket Implementation](websocket-implementation.md) for architecture
2. **Development**: Follow [Integration Guide](websocket-integration-guide.md) for implementation
3. **Testing**: Use [Testing Guide](websocket-testing-guide.md) for validation
4. **Extension**: Reference [Extension Guide](websocket-extension-guide.md) for advanced features

---

## Other Documentation

- **Database Schema**: See `backend-file-structure.md` for complete database design
- **API Endpoints**: Full REST API documentation (coming soon)
- **Frontend Components**: React component documentation (coming soon)
- **Deployment Guide**: Production deployment instructions (coming soon)

---

*This documentation was created for the social network WebSocket implementation. Last updated: 2024*
