DROP TABLE IF EXISTS waitlist_entries;

ALTER TABLE events
DROP COLUMN IF EXISTS capacity;
