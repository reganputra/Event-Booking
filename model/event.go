package model

import (
	"context"
	"go-rest-api/connection"
	"time"
)

type Event struct {
	Id          int64
	Name        string    `bind:"required"`
	Description string    `bind:"required"`
	Location    string    `bind:"required"`
	Date        time.Time `bind:"required"`
	UserIds     int64
}

var events []Event

func (e *Event) Save(ctx context.Context) error {

	insert := "INSERT INTO events (name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)"
	stmt, err := connection.DB.PrepareContext(ctx, insert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, e.Name, e.Description, e.Location, e.Date, e.UserIds)
	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.Id = lastInsertId
	return nil
}

func GetAllEvents(ctx context.Context) ([]Event, error) {

	query := "SELECT * FROM events"
	rows, err := connection.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func GetEventById(ctx context.Context, id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := connection.DB.QueryRowContext(ctx, query, id)

	var event Event
	err := row.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *Event) UpdateEvent(ctx context.Context) error {
	update := "UPDATE events SET name = ?, description = ?, location = ?, dateTime = ? WHERE id = ?"
	stmt, err := connection.DB.PrepareContext(ctx, update)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e.Name, e.Description, e.Location, e.Date, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) DeleteEvent(ctx context.Context) error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := connection.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, e.Id)
	if err != nil {
		return err
	}

	return nil
}
