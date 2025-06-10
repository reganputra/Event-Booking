package model

import (
	"go-rest-api/connection"
	"go-rest-api/helper"
	"time"
)

type Event struct {
	Id          int64
	Name        string    `bind:"required"`
	Description string    `bind:"required"`
	Location    string    `bind:"required"`
	Date        time.Time `bind:"required"`
	UserIds     int
}

var events []Event

func (e *Event) Save() error {

	insert := "INSERT INTO events (name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)"
	stmt, err := connection.DB.Prepare(insert)
	helper.PanicIfError(err)
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.Date, e.UserIds)
	helper.PanicIfError(err)

	lastInsertId, err := result.LastInsertId()
	helper.PanicIfError(err)

	e.Id = lastInsertId

	return nil
}

func GetAllEvents() []Event {
	query := "SELECT * FROM events"
	stmt, err := connection.DB.Query(query)
	helper.PanicIfError(err)
	defer stmt.Close()

	return events
}
