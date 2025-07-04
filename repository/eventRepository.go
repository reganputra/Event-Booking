package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-rest-api/model"
	"log"
)

type EventRepository interface {
	Save(ctx context.Context, event *model.Event) error
	GetAllEvents(ctx context.Context) ([]model.Event, error)
	GetEventById(ctx context.Context, id int64) (*model.Event, error)
	GetEventsByCategory(ctx context.Context, category string) ([]model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	RegisterEvent(ctx context.Context, eventID, userID int64) error
	CancelRegistration(ctx context.Context, eventID, userID int64) error
	GetRegisteredEventByUserId(ctx context.Context, userId int64) ([]model.Event, error)
}

type sqliteEventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &sqliteEventRepository{db: db}
}

func (r *sqliteEventRepository) Save(ctx context.Context, event *model.Event) error {

	insert := "INSERT INTO events (name, description, location, dateTime, category, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	stmt, err := r.db.PrepareContext(ctx, insert)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for event save: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, event.Name, event.Description, event.Location, event.Date, event.Category, event.UserIds)
	err = row.Scan(&event.Id) //
	if err != nil {
		return fmt.Errorf("failed to execute statement and scan ID for event save: %w", err)
	}

	return nil
}

func (r *sqliteEventRepository) GetAllEvents(ctx context.Context) ([]model.Event, error) {
	log.Println("Getting all events from database")

	// First, check the table schema to debug
	schemaQuery := "SELECT column_name FROM information_schema.columns WHERE table_name = 'events' ORDER BY ordinal_position"
	schemaRows, err := r.db.QueryContext(ctx, schemaQuery)
	if err != nil {
		log.Printf("Error querying schema: %v", err)
	} else {
		defer schemaRows.Close()
		log.Println("Events table columns:")
		var columnName string
		for schemaRows.Next() {
			schemaRows.Scan(&columnName)
			log.Printf("- %s", columnName)
		}
	}

	query := "SELECT * FROM events"
	log.Printf("Executing query: %s", query)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	log.Println("Scanning rows...")
	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		log.Printf("Scanned event: %+v", event)
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}
	log.Printf("Successfully retrieved %d events", len(events))
	return events, nil
}

func (r *sqliteEventRepository) GetEventById(ctx context.Context, id int64) (*model.Event, error) {
	query := "SELECT * FROM events WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var event model.Event
	err := row.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *sqliteEventRepository) Update(ctx context.Context, event *model.Event) error {
	update := "UPDATE events SET name = $1, description = $2, location = $3, dateTime = $4, category = $5 WHERE id = $6"
	stmt, err := r.db.PrepareContext(ctx, update)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, event.Name, event.Description, event.Location, event.Date, event.Category, event.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqliteEventRepository) DeleteEvent(ctx context.Context, id int64) error {
	query := "DELETE FROM events WHERE id = $1"
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
	insert := "INSERT INTO registrations (event_id, user_id) VALUES ($1, $2)"
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
	query := "DELETE FROM registrations WHERE event_id = $1 AND user_id = $2"
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

func (r *sqliteEventRepository) GetRegisteredEventByUserId(ctx context.Context, userId int64) ([]model.Event, error) {
	query := `
		SELECT
			e.id,
			e.name,
			e.description,
			e.location,
			e.dateTime,
			e.user_id,
			e.category
		FROM events AS e
		JOIN registrations AS r ON e.id = r.event_id
		WHERE r.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query registered events: %w", err)
	}
	defer rows.Close()

	events := make([]model.Event, 0)

	log.Printf("Querying registered events for user ID: %d", userId)
	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category)
		if err != nil {
			log.Printf("Error scanning registered event row: %v", err) // Log jika scan gagal
			return nil, fmt.Errorf("failed to scan registered event row: %w", err)
		}
		events = append(events, event)
		log.Printf("Scanned event: %+v", event) // Log jika event berhasil di-scan
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating registered event rows: %w", err)
	}
	log.Printf("Finished scanning. Total events found: %d", len(events)) // Log jumlah event yang ditemukan
	return events, nil
}

func (r *sqliteEventRepository) GetEventsByCategory(ctx context.Context, category string) ([]model.Event, error) {
	log.Printf("Getting events with category: %s", category)

	query := "SELECT * FROM events WHERE category = $1"
	log.Printf("Executing query: %s with category=%s", query, category)
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		log.Printf("Error executing category query: %v", err)
		return nil, fmt.Errorf("failed to query events by category: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	log.Println("Scanning category rows...")
	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category)
		if err != nil {
			log.Printf("Error scanning category row: %v", err)
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		log.Printf("Scanned category event: %+v", event)
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating category rows: %v", err)
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}
	log.Printf("Successfully retrieved %d events with category %s", len(events), category)
	return events, nil
}
