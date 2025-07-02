-- Create USERS table
CREATE TABLE IF NOT EXISTS USERS (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    date_of_birth DATE,
    avatar TEXT,
    nickname TEXT,
    about_me TEXT,
    is_profile_public BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create SESSIONS table
CREATE TABLE IF NOT EXISTS SESSIONS (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE
);
