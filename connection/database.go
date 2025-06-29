package connection

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//var DB *sql.DB

func DbConnect() *sql.DB {
	var err error
	db, err := sql.Open("sqlite3", "event_booking.db")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(5)
	err = createTable(db)
	if err != nil {
		panic(err)
	}
	return db

}

func createTable(db *sql.DB) error {
	createUsersTable := `
CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	email TEXT NOT NULL UNIQUE,
    	password TEXT NOT NULL,
    	 role TEXT NOT NULL DEFAULT 'user'
);`
	_, err := db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(createEventsTable)
	if err != nil {
		return err
	}

	createRegistrationsTable := `
    CREATE TABLE IF NOT EXISTS registrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        event_id INTEGER,
        user_id INTEGER,
        FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
	_, err = db.Exec(createRegistrationsTable)
	return err
}
