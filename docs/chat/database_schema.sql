-- =====================================================
-- SOCIAL NETWORK DATABASE SCHEMA
-- Complete database schema with all tables and indexes
-- Generated from migration files
-- =====================================================

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- =====================================================
-- CORE USER MANAGEMENT TABLES
-- =====================================================

-- Users table - Core user information
CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    date_of_birth DATE,
    avatar TEXT,
    nickname TEXT,
    about_me TEXT,
    is_profile_public INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table - User authentication sessions
CREATE TABLE IF NOT EXISTS Sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- =====================================================
-- SOCIAL CONNECTIONS
-- =====================================================

-- Followers table - User follow relationships
CREATE TABLE IF NOT EXISTS Followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followee_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME,
    FOREIGN KEY (follower_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES Users(id) ON DELETE CASCADE,
    UNIQUE (follower_id, followee_id)
);

-- =====================================================
-- POSTS AND CONTENT
-- =====================================================

-- Posts table - User posts/content
CREATE TABLE IF NOT EXISTS Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    privacy TEXT NOT NULL DEFAULT 'public',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Post_Visibility table - Controls who can see specific posts
CREATE TABLE IF NOT EXISTS Post_Visibility (
    post_id INTEGER NOT NULL,
    viewer_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, viewer_id),
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (viewer_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Comments table - Comments on posts
CREATE TABLE IF NOT EXISTS Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- =====================================================
-- GROUPS SYSTEM
-- =====================================================

-- Groups table - User groups/communities
CREATE TABLE IF NOT EXISTS Groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    creator_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Group_Members table - Group membership and roles
CREATE TABLE IF NOT EXISTS Group_Members (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT,
    is_accepted BOOLEAN DEFAULT 0,
    invited_by INTEGER,
    requested BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (invited_by) REFERENCES Users(id) ON DELETE SET NULL
);

-- Group_Posts table - Posts within groups
CREATE TABLE IF NOT EXISTS Group_Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT,
    image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- =====================================================
-- EVENTS SYSTEM
-- =====================================================

-- Group_Events table - Events within groups
CREATE TABLE IF NOT EXISTS Group_Events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    event_time DATETIME,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES Users(id) ON DELETE CASCADE
);

-- Event_Responses table - User responses to events (going/not going)
CREATE TABLE IF NOT EXISTS Event_Responses (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT,
    responded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES Group_Events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- =====================================================
-- MESSAGING SYSTEM
-- =====================================================

-- Messages table - Private and group messages
CREATE TABLE IF NOT EXISTS Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    content TEXT,
    is_emoji BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE
);

-- =====================================================
-- NOTIFICATIONS SYSTEM
-- =====================================================

-- Notifications table - Basic notification system (legacy)
CREATE TABLE IF NOT EXISTS Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT,
    message TEXT,
    is_read BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Activity_Notifications table - Rich notification system (new)
CREATE TABLE IF NOT EXISTS Activity_Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK(type IN (
        'follow_request',
        'follow_accepted', 
        'follow_declined',
        'group_invitation',
        'group_join_request',
        'group_event_created',
        'new_message'
    )),
    from_user_id INTEGER NOT NULL,
    reference_id INTEGER, -- Points to Followers.id, Group_Members composite key, Message.id, etc.
    data TEXT, -- JSON string with notification-specific data
    is_read BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (from_user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- =====================================================
-- DATABASE INDEXES FOR PERFORMANCE
-- =====================================================

-- Followers table indexes
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON Followers(follower_id);
CREATE INDEX IF NOT EXISTS idx_followers_followee_id ON Followers(followee_id);

-- Posts table indexes
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON Posts(user_id);

-- Post_Visibility table indexes
CREATE INDEX IF NOT EXISTS idx_post_visibility_post_id ON Post_Visibility(post_id);
CREATE INDEX IF NOT EXISTS idx_post_visibility_viewer_id ON Post_Visibility(viewer_id);

-- Comments table indexes
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON Comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON Comments(user_id);

-- Groups table indexes
CREATE INDEX IF NOT EXISTS idx_groups_creator_id ON Groups(creator_id);

-- Group_Members table indexes
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON Group_Members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON Group_Members(user_id);

-- Group_Posts table indexes
CREATE INDEX IF NOT EXISTS idx_group_posts_group_id ON Group_Posts(group_id);
CREATE INDEX IF NOT EXISTS idx_group_posts_user_id ON Group_Posts(user_id);

-- Group_Events table indexes
CREATE INDEX IF NOT EXISTS idx_group_events_group_id ON Group_Events(group_id);
CREATE INDEX IF NOT EXISTS idx_group_events_created_by ON Group_Events(created_by);

-- Event_Responses table indexes
CREATE INDEX IF NOT EXISTS idx_event_responses_event_id ON Event_Responses(event_id);
CREATE INDEX IF NOT EXISTS idx_event_responses_user_id ON Event_Responses(user_id);

-- Messages table indexes
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON Messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON Messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_group_id ON Messages(group_id);

-- Notifications table indexes (legacy)
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON Notifications(user_id);

-- Activity_Notifications table indexes (new system)
CREATE INDEX IF NOT EXISTS idx_activity_notifications_user_id ON Activity_Notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_notifications_type ON Activity_Notifications(type);
CREATE INDEX IF NOT EXISTS idx_activity_notifications_user_type ON Activity_Notifications(user_id, type);
CREATE INDEX IF NOT EXISTS idx_activity_notifications_created_at ON Activity_Notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_activity_notifications_is_read ON Activity_Notifications(user_id, is_read);

-- =====================================================
-- SAMPLE DATA (OPTIONAL - FOR DEVELOPMENT/TESTING)
-- =====================================================

-- Uncomment the following section if you want to insert sample data

/*
-- Sample users
INSERT OR IGNORE INTO Users (id, email, password, first_name, last_name, nickname, avatar) VALUES
    (1, 'john@example.com', 'hashed_password_1', 'John', 'Doe', 'john_doe', '/uploads/avatars/123.jpg'),
    (2, 'jane@example.com', 'hashed_password_2', 'Jane', 'Smith', 'jane_smith', '/uploads/avatars/456.jpg'),
    (3, 'bob@example.com', 'hashed_password_3', 'Bob', 'Wilson', 'bob_wilson', '/uploads/avatars/789.jpg');

-- Sample groups
INSERT OR IGNORE INTO Groups (id, title, description, creator_id) VALUES
    (1, 'Photography Club', 'A group for photography enthusiasts', 1),
    (2, 'Tech Discussions', 'Discuss the latest in technology', 2);

-- Sample group memberships
INSERT OR IGNORE INTO Group_Members (group_id, user_id, role, is_accepted) VALUES
    (1, 1, 'admin', 1),
    (1, 2, 'member', 1),
    (2, 2, 'admin', 1),
    (2, 3, 'member', 1);
*/

-- =====================================================
-- DATABASE SCHEMA COMPLETE
-- =====================================================

-- This schema includes:
-- ✅ 13 Core tables (Users, Sessions, Followers, Posts, etc.)
-- ✅ All foreign key relationships
-- ✅ Performance indexes on frequently queried columns
-- ✅ Check constraints for data validation
-- ✅ Both legacy and new notification systems
-- ✅ Complete social network functionality:
--    - User management and authentication
--    - Social connections (followers/following)
--    - Posts and comments with privacy controls
--    - Groups and group posts
--    - Events and event responses
--    - Private and group messaging
--    - Comprehensive notification system

-- To use this schema:
-- 1. Create a new SQLite database
-- 2. Run this entire file to create all tables and indexes
-- 3. Optionally uncomment and run the sample data section
-- 4. Your social network database is ready!
