-- This migration is destructive and cannot be perfectly reversed.
-- This 'down' migration will attempt to restore the schema to its previous state,
-- but all data will be lost.

-- Drop the new foreign key constraints
ALTER TABLE events DROP CONSTRAINT IF EXISTS fk_events_users;
ALTER TABLE registrations DROP CONSTRAINT IF EXISTS fk_registrations_events;
ALTER TABLE registrations DROP CONSTRAINT IF EXISTS fk_registrations_users;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_events;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_users;
ALTER TABLE waitlist_entries DROP CONSTRAINT IF EXISTS fk_waitlist_entries_events;
ALTER TABLE waitlist_entries DROP CONSTRAINT IF EXISTS fk_waitlist_entries_users;

-- Alter tables back to using SERIAL integer primary keys
-- Alter 'users' table
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users ADD COLUMN id SERIAL PRIMARY KEY;

-- Alter 'events' table
ALTER TABLE events DROP COLUMN IF EXISTS id;
ALTER TABLE events ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE events DROP COLUMN IF EXISTS user_id;
ALTER TABLE events ADD COLUMN user_id INTEGER;

-- Alter 'registrations' table
ALTER TABLE registrations DROP COLUMN IF EXISTS id;
ALTER TABLE registrations ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE registrations DROP COLUMN IF EXISTS event_id;
ALTER TABLE registrations ADD COLUMN event_id INTEGER;
ALTER TABLE registrations DROP COLUMN IF EXISTS user_id;
ALTER TABLE registrations ADD COLUMN user_id INTEGER;

-- Alter 'reviews' table
ALTER TABLE reviews DROP COLUMN IF EXISTS id;
ALTER TABLE reviews ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE reviews DROP COLUMN IF EXISTS event_id;
ALTER TABLE reviews ADD COLUMN event_id INTEGER;
ALTER TABLE reviews DROP COLUMN IF EXISTS user_id;
ALTER TABLE reviews ADD COLUMN user_id INTEGER;

-- Alter 'waitlist_entries' table
ALTER TABLE waitlist_entries DROP COLUMN IF EXISTS id;
ALTER TABLE waitlist_entries ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE waitlist_entries DROP COLUMN IF EXISTS event_id;
ALTER TABLE waitlist_entries ADD COLUMN event_id INTEGER;
ALTER TABLE waitlist_entries DROP COLUMN IF EXISTS user_id;
ALTER TABLE waitlist_entries ADD COLUMN user_id INTEGER;

-- Re-create the original foreign key constraints
ALTER TABLE events ADD CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE registrations ADD CONSTRAINT registrations_event_id_fkey FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE registrations ADD CONSTRAINT registrations_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE reviews ADD CONSTRAINT reviews_event_id_fkey FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE reviews ADD CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE waitlist_entries ADD CONSTRAINT waitlist_entries_event_id_fkey FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE waitlist_entries ADD CONSTRAINT waitlist_entries_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;