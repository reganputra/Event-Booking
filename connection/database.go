package connection

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func DbConnect() {
	var err error

	db, err = sql.Open("sqlite3", "event_booking.db")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(5)
	createTable()

}

func createTable() {
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER
	);
	`
	_, err := db.Exec(createEventsTable)
	if err != nil {
		panic(err)
	}
}
