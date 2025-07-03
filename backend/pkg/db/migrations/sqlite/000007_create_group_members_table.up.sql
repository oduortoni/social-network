-- Create Group_Members table
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

CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON Group_Members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON Group_Members(user_id);
