-- migrations/000002_add_category_to_events.up.sql

ALTER TABLE events
ADD COLUMN category TEXT NOT NULL DEFAULT 'General';