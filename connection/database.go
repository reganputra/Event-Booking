package connection

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func DbConnect() {
	var err error

	DB, err = sql.Open("sqlite3", "event_booking.db")
	if err != nil {
		panic(err)
	}

	DB.SetMaxOpenConns(5)
	err = createTable()
	if err != nil {
		panic(err)
	}

}

func createTable() error {
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
	_, err := DB.Exec(createEventsTable)
	return err
}
