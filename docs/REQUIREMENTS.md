# Social Network Project Requirements

## Project Overview
Create a Facebook-like social network with comprehensive social features including user profiles, posts, groups, real-time chat, and notifications.

## Technology Stack

### Frontend
- **Framework**: JavaScript framework (Next.js, Vue.js, Svelte, or Mithril)
- **Languages**: HTML, CSS, JavaScript
- **Focus**: Responsiveness and performance
- **Architecture**: Single Page Application (SPA) with dynamic content loading

### Backend
- **Language**: Go
- **Database**: SQLite
- **Real-time**: WebSocket (Gorilla WebSocket)
- **Authentication**: Sessions and cookies
- **File Handling**: JPEG, PNG, GIF support

### Allowed Packages
- Standard Go packages
- Gorilla WebSocket
- golang-migrate / sql-migration / migration
- sqlite3
- bcrypt
- gofrs/uuid or google/uuid

## Core Features

### 1. Authentication System
**Status**: ✅ Implemented

#### Registration Requirements
- [x] Email (required)
- [x] Password (required)
- [x] First Name (required)
- [x] Last Name (required)
- [x] Date of Birth (required)
- [x] Avatar/Image (optional)
- [x] Nickname (optional)
- [x] About Me (optional)

#### Login System
- [x] Session-based authentication
- [x] Cookie management
- [x] Persistent login until logout
- [x] Logout functionality available at all times

### 2. User Profiles
**Status**: 🔄 In Progress

#### Profile Information
- [ ] User information display (excluding password)
- [ ] User activity (posts made by user)
- [ ] Followers and following lists
- [ ] Profile privacy toggle (public/private)

#### Profile Types
- [ ] **Public Profile**: Visible to all users
- [ ] **Private Profile**: Visible to followers only
- [ ] Profile privacy settings in user's own profile

### 3. Followers System
**Status**: 🔄 In Progress

#### Follow Functionality
- [x] Follow/unfollow other users
- [ ] Follow request system for private profiles
- [ ] Accept/decline follow requests
- [ ] Automatic following for public profiles (bypass request)

### 4. Posts System
**Status**: 📋 Planned

#### Post Creation
- [ ] Text posts
- [ ] Image/GIF attachments (JPEG, PNG, GIF)
- [ ] Comment system on posts

#### Post Privacy Levels
- [ ] **Public**: Visible to all users
- [ ] **Almost Private**: Visible to followers only
- [ ] **Private**: Visible to selected followers only

### 5. Groups System
**Status**: 📋 Planned

#### Group Management
- [ ] Create groups with title and description
- [ ] Invite users to groups
- [ ] Accept/decline group invitations
- [ ] Request to join groups
- [ ] Creator approval for join requests
- [ ] Browse all groups section

#### Group Features
- [ ] Group-specific posts and comments
- [ ] Group member-only content visibility
- [ ] Group chat room

#### Events System
- [ ] Create events within groups
- [ ] Event details: title, description, date/time
- [ ] RSVP options: Going/Not Going
- [ ] Event notifications to group members

### 6. Chat System
**Status**: 🔄 In Progress

#### Private Messaging
- [x] Real-time WebSocket communication
- [x] Message persistence
- [ ] Private messages between followers
- [ ] Emoji support
- [ ] Message history

#### Group Chat
- [ ] Group chat rooms
- [ ] Real-time group messaging
- [ ] Group member participation

### 7. Notification System
**Status**: ✅ Implemented (Basic)

#### Current Features
- [x] Real-time WebSocket notifications
- [x] User connection/disconnection alerts
- [x] Bell icon with unread count badge
- [x] Notification panel with message history

#### Required Notifications
- [ ] Follow request notifications (private profiles)
- [ ] Group invitation notifications
- [ ] Group join request notifications (for group creators)
- [ ] Event creation notifications (for group members)
- [ ] Custom notifications as needed

#### Notification Display
- [x] Visible on every page
- [x] Different display from private messages
- [ ] Accept/decline actions for requests

## Technical Requirements

### Database Design
- [x] SQLite database implementation
- [x] Entity relationship diagram planning
- [x] Database migrations
- [x] User sessions and authentication tables
- [x] WebSocket connection tracking

### Image Handling
- [ ] Support for JPEG, PNG, GIF formats
- [ ] File storage system
- [ ] Database path storage
- [ ] Image optimization and serving

### WebSocket Implementation
- [x] Real-time client connections
- [x] Private chat functionality
- [x] Connection management
- [x] Message broadcasting
- [ ] Group chat implementation
- [ ] Notification delivery

### Security
- [x] Password encryption (bcrypt)
- [x] Session management
- [x] Cookie security
- [ ] Input validation and sanitization
- [ ] File upload security

## Development Progress

### Completed ✅
- User authentication system
- Session and cookie management
- Basic WebSocket infrastructure
- Real-time notification system
- User connection tracking
- Database schema and migrations
- Frontend-backend integration

### In Progress 🔄
- User profile system
- Follow/unfollow functionality
- Private messaging system
- Connected users display

### Planned 📋
- Posts and comments system
- Groups and events
- Comprehensive notification types
- Image upload and handling
- Group chat functionality

## Architecture Decisions

### Frontend Architecture
- Single Page Application (SPA)
- Dynamic content loading without URL changes
- Real-time updates via WebSocket
- Centralized state management
- Component-based UI structure

### Backend Architecture
- RESTful API design
- WebSocket for real-time features
- Session-based authentication
- Modular service architecture
- Database abstraction layer

### Data Flow
- WebSocket for real-time updates
- REST API for CRUD operations
- Event-driven notification system
- Optimistic UI updates
- Background data synchronization

## Testing Strategy
- [x] Unit tests for WebSocket functionality
- [x] Integration tests for API endpoints
- [x] Database migration testing
- [ ] End-to-end testing for user flows
- [ ] Performance testing for real-time features

## Deployment Considerations
- [ ] Docker containerization
- [ ] Environment configuration
- [ ] Database migration scripts
- [ ] Static file serving
- [ ] WebSocket scaling considerations

---

## Next Steps
1. Complete user profile system
2. Implement posts and comments
3. Add comprehensive notification types
4. Develop groups and events functionality
5. Enhance chat system with emoji support
6. Add image upload capabilities
7. Implement comprehensive testing
8. Prepare for deployment
