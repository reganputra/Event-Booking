CREATE TABLE IF NOT EXISTS waitlist_entries (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (event_id, user_id) -- A user can only be on the waitlist for an event once
);

-- Add capacity to events table
ALTER TABLE events
ADD COLUMN IF NOT EXISTS capacity INTEGER DEFAULT 0; -- Default 0 means no limit unless specified otherwise.
