package repository

import (
	"context"
	"database/sql"
	"go-rest-api/model"
)

type EventRepository interface {
	Save(ctx context.Context, event *model.Event) error
	GetAllEvents(ctx context.Context) ([]model.Event, error)
	GetEventById(ctx context.Context, id int64) (*model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	RegisterEvent(ctx context.Context, eventID, userID int64) error
	CancelRegistration(ctx context.Context, eventID, userID int64) error
}

type sqliteEventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &sqliteEventRepository{db: db}
}

func (r *sqliteEventRepository) Save(ctx context.Context, event *model.Event) error {

	insert := "INSERT INTO events (name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)"
	stmt, err := r.db.PrepareContext(ctx, insert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, event.Name, event.Description, event.Location, event.Date, event.UserIds)
	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.Id = lastInsertId
	return nil
}

func (r *sqliteEventRepository) GetAllEvents(ctx context.Context) ([]model.Event, error) {

	query := "SELECT * FROM events"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
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

func (r *sqliteEventRepository) GetEventById(ctx context.Context, id int64) (*model.Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var event model.Event
	err := row.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *sqliteEventRepository) Update(ctx context.Context, event *model.Event) error {
	update := "UPDATE events SET name = ?, description = ?, location = ?, dateTime = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, update)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, event.Name, event.Description, event.Location, event.Date, event.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqliteEventRepository) DeleteEvent(ctx context.Context, id int64) error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqliteEventRepository) RegisterEvent(ctx context.Context, eventId, userId int64) error {
	insert := "INSERT INTO registrations (event_id, user_id) VALUES (?, ?)"
	stmt, err := r.db.PrepareContext(ctx, insert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, eventId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *sqliteEventRepository) CancelRegistration(ctx context.Context, eventId, userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, eventId, userId)
	if err != nil {
		return err
	}
	return nil
}
