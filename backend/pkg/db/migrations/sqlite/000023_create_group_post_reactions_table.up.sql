-- Create Group_Post_Reactions table
CREATE TABLE IF NOT EXISTS Group_Post_Reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    reaction_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_post_id) REFERENCES Group_Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    UNIQUE(group_post_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_group_post_reactions_group_post_id ON Group_Post_Reactions(group_post_id);
CREATE INDEX IF NOT EXISTS idx_group_post_reactions_user_id ON Group_Post_Reactions(user_id);
