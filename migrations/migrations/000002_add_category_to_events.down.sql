-- migrations/000002_add_category_to_events.down.sql

ALTER TABLE events
DROP COLUMN category;