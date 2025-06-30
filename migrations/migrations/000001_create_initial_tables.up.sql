-- migrations/000001_create_initial_tables.up.sql

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY, -- Menggunakan SERIAL untuk auto-increment di PostgreSQL
                                     email TEXT NOT NULL UNIQUE,
                                     password TEXT NOT NULL,
                                     role TEXT NOT NULL DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS events (
                                      id SERIAL PRIMARY KEY,
                                      name TEXT NOT NULL,
                                      description TEXT NOT NULL,
                                      location TEXT NOT NULL,
                                      dateTime TIMESTAMP NOT NULL, -- Menggunakan TIMESTAMP untuk tanggal/waktu di PostgreSQL
                                      user_id INTEGER,
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS registrations (
                                             id SERIAL PRIMARY KEY,
                                             event_id INTEGER NOT NULL,
                                             user_id INTEGER NOT NULL,
                                             FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (event_id, user_id) -- Tambahkan constraint UNIQUE untuk mencegah duplikasi pendaftaran
    );