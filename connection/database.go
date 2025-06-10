package connection

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-rest-api/helper"
)

var DB *sql.DB

func DbConnect() {
	var err error

	DB, err = sql.Open("sqlite3", "event_booking.db")
	helper.PanicIfError(err)

	DB.SetMaxOpenConns(5)
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
	_, err := DB.Exec(createEventsTable)
	helper.PanicIfError(err)
}
