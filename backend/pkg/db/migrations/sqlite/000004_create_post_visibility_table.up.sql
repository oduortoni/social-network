CREATE TABLE IF NOT EXISTS post_visibility (
    post_id INTEGER NOT NULL,
    viewer_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, viewer_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (viewer_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_post_visibility_post_id ON post_visibility(post_id);
CREATE INDEX IF NOT EXISTS idx_post_visibility_viewer_id ON post_visibility(viewer_id);