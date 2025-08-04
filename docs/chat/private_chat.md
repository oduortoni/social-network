# Private Chat System Documentation

## Overview

The private chat system enables real-time messaging between authenticated users within the social network platform. Users can initiate conversations with mutual followers or users with public profiles, with messages delivered instantly via WebSocket connections and persisted to the database.

## Architecture

### Frontend Components
- **Header Navigation**: Clickable "Chats" icon for easy access
- **Chat List Page** (`/chats`): Displays all messageable users
- **Individual Chat Page** (`/chats/[userId]`): Focused 1-on-1 messaging interface
- **ChatInterface Component**: Reusable chat UI with WebSocket integration

### Backend Services
- **WebSocket Manager**: Handles real-time message routing
- **Permission Checker**: Validates user chat permissions
- **Message Persister**: Database operations for message storage
- **Chat API Endpoints**: REST endpoints for chat history and user lists

## User Flow

### 1. Accessing Chats
Users click the "Chats" icon in the navigation header, which navigates to `/chats` displaying all users they can message.

### 2. User Selection
The system fetches messageable users from `/api/users/messageable` and displays them in a clean list format.

**API Request:**
```http
GET /api/users/messageable
Authorization: Session Cookie
```

**API Response:**
```json
[
  {
    "id": 2,
    "nickname": "john_doe",
    "avatar": "avatar_filename.jpg"
  },
  {
    "id": 3,
    "nickname": "jane_smith", 
    "avatar": "no profile photo"
  }
]
```

### 3. Opening Individual Chat
Clicking on a user navigates to `/chats/[userId]?nickname=username`, opening a dedicated chat interface.

### 4. Message History Loading
The system loads existing conversation history from the database.

**API Request:**
```http
GET /api/messages/private/2
Authorization: Session Cookie
```

**API Response:**
```json
[
  {
    "from": 1,
    "to": 2,
    "content": "Hello there!",
    "timestamp": 1672531200
  },
  {
    "from": 2,
    "to": 1,
    "content": "Hi! How are you?",
    "timestamp": 1672531260
  }
]
```

## Real-Time Messaging

### WebSocket Connection
The frontend establishes a WebSocket connection to `/ws` upon user authentication.

**Connection Handshake:**
```javascript
// Frontend WebSocket initialization
const ws = new WebSocket('ws://localhost:9000/ws');
ws.onopen = () => console.log('Connected to chat server');
```

### Message Sending
When users type and send messages, the system implements optimistic updates for immediate UI feedback.

**WebSocket Message Format:**
```json
{
  "type": "private",
  "to": 2,
  "content": "Hello there!",
  "timestamp": 1672531200
}
```

**Backend Processing:**
1. Validate user permissions
2. Save message to database
3. Forward to both sender and recipient if online

### Message Reception
Recipients receive messages in real-time without page refreshes.

**Received Message Format:**
```json
{
  "type": "private",
  "from": 1,
  "to": 2,
  "content": "Hello there!",
  "timestamp": 1672531200
}
```

## Design Decisions

### 1. Focused Chat Interface
**Decision**: Individual chat pages show only the conversation between two users, without sidebar distractions.

**Rationale**: Provides focused, distraction-free messaging experience similar to modern chat applications.

### 2. Optimistic Updates
**Decision**: Messages appear immediately in sender's interface before server confirmation.

**Rationale**: Creates responsive user experience while maintaining data consistency through server validation.

### 3. Permission-Based Messaging
**Decision**: Users can only message mutual followers or users with public profiles.

**Rationale**: Prevents spam while allowing organic conversation growth through the social network.

### 4. WebSocket + REST Hybrid
**Decision**: Use WebSocket for real-time messaging and REST for chat history/user lists.

**Rationale**: Leverages strengths of both protocols - WebSocket for real-time updates, REST for reliable data fetching.

### 5. Message Ordering
**Decision**: Display messages in chronological order (oldest first).

**Rationale**: Natural conversation flow that users expect from messaging applications.

## Technical Implementation

### Frontend State Management
The ChatInterface component uses React hooks for state management:

```javascript
const [messages, setMessages] = useState([]);
const [activeChat, setActiveChat] = useState(null);
const [connectionStatus, setConnectionStatus] = useState('disconnected');
```

### WebSocket Handler Pattern
Stable message handlers prevent duplicate event listeners:

```javascript
const handlePrivateMessage = useCallback((message) => {
  // Message filtering and UI updates
}, []); // No dependencies to prevent handler recreation
```

### Backend Message Routing
The WebSocket manager routes messages based on user permissions:

```go
// Validate permissions
allowed, err := m.PermissionChecker.CanUsersChat(senderID, recipientID)
if allowed {
    // Save to database
    m.persister.SaveMessage(senderID, message)
    // Send to both users
    m.SendToUser(recipientID, encodedMessage)
    m.SendToUser(senderID, encodedMessage)
}
```

## Database Schema

### Messages Table
```sql
CREATE TABLE Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES Users(id),
    FOREIGN KEY (receiver_id) REFERENCES Users(id)
);
```

## Error Handling

### Permission Denied
When users attempt to message unauthorized recipients:

**Response:**
```json
{
  "type": "error",
  "content": "You are not permitted to message this user.",
  "timestamp": 1672531200
}
```

### Connection Issues
The frontend gracefully handles WebSocket disconnections with visual indicators and automatic reconnection attempts.

## Security Considerations

### Input Sanitization
All user input is sanitized on the backend to prevent XSS attacks:

```go
content = html.EscapeString(content)
```

### Permission Validation
Every message is validated against user permissions before processing or storage.

### Session Management
WebSocket connections are authenticated using session cookies, ensuring only logged-in users can access chat features.

## Conclusion

The private chat system provides a robust, real-time messaging experience that integrates seamlessly with the social network's permission model. The combination of WebSocket real-time updates, optimistic UI patterns, and focused user interface design creates an engaging communication platform for users.
