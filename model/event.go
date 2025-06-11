package model

import (
	"context"
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

func (e *Event) Save(ctx context.Context) error {

	insert := "INSERT INTO events (name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)"
	stmt, err := connection.DB.PrepareContext(ctx, insert)
	helper.PanicIfError(err)
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, e.Name, e.Description, e.Location, e.Date, e.UserIds)
	helper.PanicIfError(err)

	lastInsertId, err := result.LastInsertId()
	helper.PanicIfError(err)

	e.Id = lastInsertId

	return nil
}

func GetAllEvents(ctx context.Context) []Event {
	query := "SELECT * FROM events"
	rows, err := connection.DB.QueryContext(ctx, query)
	helper.PanicIfError(err)
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds)
		helper.PanicIfError(err)
		events = append(events, event)
	}

	return events
}
