-- migrations/000005_add_name_and_rating_to_events.up.sql

ALTER TABLE events ADD COLUMN IF NOT EXISTS average_rating REAL DEFAULT 0;