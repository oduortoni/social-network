-- Create Event_Responses table
CREATE TABLE IF NOT EXISTS Event_Responses (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT,
    responded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES Group_Events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_event_responses_event_id ON Event_Responses(event_id);
CREATE INDEX IF NOT EXISTS idx_event_responses_user_id ON Event_Responses(user_id);
