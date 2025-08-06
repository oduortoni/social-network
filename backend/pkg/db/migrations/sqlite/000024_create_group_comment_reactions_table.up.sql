-- Create Group_Comment_Reactions table
CREATE TABLE IF NOT EXISTS Group_Comment_Reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    reaction_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES Group_Post_Comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    UNIQUE(comment_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_group_comment_reactions_comment_id ON Group_Comment_Reactions(comment_id);
CREATE INDEX IF NOT EXISTS idx_group_comment_reactions_user_id ON Group_Comment_Reactions(user_id);
