-- Enable the pgcrypto extension to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Drop existing foreign key constraints
ALTER TABLE events DROP CONSTRAINT IF EXISTS events_user_id_fkey;
ALTER TABLE registrations DROP CONSTRAINT IF EXISTS registrations_event_id_fkey;
ALTER TABLE registrations DROP CONSTRAINT IF EXISTS registrations_user_id_fkey;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_event_id_fkey;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_user_id_fkey;
ALTER TABLE waitlist_entries DROP CONSTRAINT IF EXISTS waitlist_entries_event_id_fkey;
ALTER TABLE waitlist_entries DROP CONSTRAINT IF EXISTS waitlist_entries_user_id_fkey;

-- Alter tables to use UUIDs
-- Note: This migration is destructive and will result in data loss.

-- Alter 'users' table
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();

-- Alter 'events' table
ALTER TABLE events DROP CONSTRAINT IF EXISTS events_pkey CASCADE;
ALTER TABLE events DROP COLUMN IF EXISTS id;
ALTER TABLE events ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE events DROP COLUMN IF EXISTS user_id;
ALTER TABLE events ADD COLUMN user_id UUID;

-- Alter 'registrations' table
ALTER TABLE registrations DROP CONSTRAINT IF EXISTS registrations_pkey CASCADE;
ALTER TABLE registrations DROP COLUMN IF EXISTS id;
ALTER TABLE registrations ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE registrations DROP COLUMN IF EXISTS event_id;
ALTER TABLE registrations ADD COLUMN event_id UUID;
ALTER TABLE registrations DROP COLUMN IF EXISTS user_id;
ALTER TABLE registrations ADD COLUMN user_id UUID;

-- Alter 'reviews' table
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_pkey CASCADE;
ALTER TABLE reviews DROP COLUMN IF EXISTS id;
ALTER TABLE reviews ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE reviews DROP COLUMN IF EXISTS event_id;
ALTER TABLE reviews ADD COLUMN event_id UUID;
ALTER TABLE reviews DROP COLUMN IF EXISTS user_id;
ALTER TABLE reviews ADD COLUMN user_id UUID;

-- Alter 'waitlist_entries' table
ALTER TABLE waitlist_entries DROP CONSTRAINT IF EXISTS waitlist_entries_pkey CASCADE;
ALTER TABLE waitlist_entries DROP COLUMN IF EXISTS id;
ALTER TABLE waitlist_entries ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE waitlist_entries DROP COLUMN IF EXISTS event_id;
ALTER TABLE waitlist_entries ADD COLUMN event_id UUID;
ALTER TABLE waitlist_entries DROP COLUMN IF EXISTS user_id;
ALTER TABLE waitlist_entries ADD COLUMN user_id UUID;

-- Re-create foreign key constraints
ALTER TABLE events ADD CONSTRAINT fk_events_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE registrations ADD CONSTRAINT fk_registrations_events FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE registrations ADD CONSTRAINT fk_registrations_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE reviews ADD CONSTRAINT fk_reviews_events FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE reviews ADD CONSTRAINT fk_reviews_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE waitlist_entries ADD CONSTRAINT fk_waitlist_entries_events FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE waitlist_entries ADD CONSTRAINT fk_waitlist_entries_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;