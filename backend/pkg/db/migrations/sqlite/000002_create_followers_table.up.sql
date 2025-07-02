-- Create FOLLOWERS table
CREATE TABLE IF NOT EXISTS FOLLOWERS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followee_id INTEGER NOT NULL,
    is_accepted BOOLEAN DEFAULT 0,
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME,
    FOREIGN KEY (follower_id) REFERENCES USERS(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES USERS(id) ON DELETE CASCADE,
    UNIQUE (follower_id, followee_id)
);

-- Create INDEX on FOLLOWERS table for faster lookups
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON FOLLOWERS(follower_id);
CREATE INDEX IF NOT EXISTS idx_followers_followee_id ON FOLLOWERS(followee_id);