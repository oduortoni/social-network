-- Create Post_Visibility table
CREATE TABLE IF NOT EXISTS Post_Visibility (
    post_id INTEGER NOT NULL,
    viewer_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, viewer_id),
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (viewer_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_post_visibility_post_id ON Post_Visibility(post_id);
CREATE INDEX IF NOT EXISTS idx_post_visibility_viewer_id ON Post_Visibility(viewer_id);
