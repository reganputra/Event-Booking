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
	GetEventsByCriteria(ctx context.Context, keyword string, startDate string, endDate string) ([]model.Event, error)
	UpdateAverageRating(ctx context.Context, eventID int64, avgRating float64) error
	Update(ctx context.Context, event *model.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	RegisterEvent(ctx context.Context, eventID, userID int64) error
	GetRegistrationCount(ctx context.Context, eventID int64) (int, error)
	IsUserRegistered(ctx context.Context, eventID int64, userID int64) (bool, error)
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

	// Include capacity in the INSERT statement
	insert := "INSERT INTO events (name, description, location, dateTime, category, user_id, capacity) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	stmt, err := r.db.PrepareContext(ctx, insert)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for event save: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, event.Name, event.Description, event.Location, event.Date, event.Category, event.UserIds, event.Capacity)
	err = row.Scan(&event.Id)
	if err != nil {
		return fmt.Errorf("failed to execute statement and scan ID for event save: %w", err)
	}

	return nil
}

func (r *sqliteEventRepository) GetRegistrationCount(ctx context.Context, eventID int64) (int, error) {
	query := "SELECT COUNT(*) FROM registrations WHERE event_id = $1"
	var count int
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get registration count for event %d: %w", eventID, err)
	}
	return count, nil
}

func (r *sqliteEventRepository) IsUserRegistered(ctx context.Context, eventID int64, userID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM registrations WHERE event_id = $1 AND user_id = $2)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, eventID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user %d is registered for event %d: %w", userID, eventID, err)
	}
	return exists, nil
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
	// Ensure all columns including 'capacity' are selected
	query := "SELECT id, name, description, location, dateTime, user_id, category, average_rating, capacity FROM events"
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
		// Add capacity to scan
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category, &event.AverageRating, &event.Capacity)
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
	// Ensure all columns including 'capacity' are selected
	query := "SELECT id, name, description, location, dateTime, user_id, category, average_rating, capacity FROM events WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var event model.Event
	// Add capacity to scan
	err := row.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category, &event.AverageRating, &event.Capacity)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *sqliteEventRepository) Update(ctx context.Context, event *model.Event) error {
	// Note: average_rating is not updated here, it's handled by UpdateAverageRating
	// Include capacity in the UPDATE statement
	updateQuery := "UPDATE events SET name = $1, description = $2, location = $3, dateTime = $4, category = $5, capacity = $6 WHERE id = $7"
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, event.Name, event.Description, event.Location, event.Date, event.Category, event.Capacity, event.Id)
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
			e.category,
			e.average_rating,
			e.capacity -- Add capacity
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
		// Add capacity to scan
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category, &event.AverageRating, &event.Capacity)
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
	// Ensure all columns including 'capacity' are selected
	query := "SELECT id, name, description, location, dateTime, user_id, category, average_rating, capacity FROM events WHERE category = $1"
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
		// Add capacity to scan
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category, &event.AverageRating, &event.Capacity)
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

func (r *sqliteEventRepository) GetEventsByCriteria(ctx context.Context, keyword string, startDate string, endDate string) ([]model.Event, error) {
	// Ensure all columns including 'capacity' are selected
	query := "SELECT id, name, description, location, dateTime, user_id, category, average_rating, capacity FROM events WHERE 1=1"
	args := []interface{}{}
	argId := 1

	if keyword != "" {
		// PostgreSQL uses ILIKE for case-insensitive search. For SQLite, use LIKE and ensure keywords are lower/upper cased if needed, or use COLLATE NOCASE on schema.
		// Assuming PostgreSQL or similar that supports ILIKE.
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argId, argId)
		args = append(args, "%"+keyword+"%")
		argId++
	}

	if startDate != "" {
		query += fmt.Sprintf(" AND dateTime >= $%d", argId)
		args = append(args, startDate)
		argId++
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND dateTime <= $%d", argId)
		args = append(args, endDate)
		argId++
	}

	log.Printf("Executing query: %s with args: %v", query, args)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing search query: %v", err)
		return nil, fmt.Errorf("failed to query events by criteria: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	log.Println("Scanning search result rows...")
	for rows.Next() {
		var event model.Event
		// Add capacity to scan
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserIds, &event.Category, &event.AverageRating, &event.Capacity)
		if err != nil {
			log.Printf("Error scanning search row: %v", err)
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		log.Printf("Scanned search event: %+v", event)
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating search rows: %v", err)
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}
	log.Printf("Successfully retrieved %d events from search", len(events))
	return events, nil
}

func (r *sqliteEventRepository) UpdateAverageRating(ctx context.Context, eventID int64, avgRating float64) error {
	query := "UPDATE events SET average_rating = $1 WHERE id = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for update average rating: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, avgRating, eventID)
	if err != nil {
		return fmt.Errorf("failed to execute statement for update average rating: %w", err)
	}
	return nil
}
