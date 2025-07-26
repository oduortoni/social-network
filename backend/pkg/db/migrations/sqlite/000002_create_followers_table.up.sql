-- Create Followers table
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

-- Create INDEX on Followers table for faster lookups
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON Followers(follower_id);
CREATE INDEX IF NOT EXISTS idx_followers_followee_id ON Followers(followee_id);