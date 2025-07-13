package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-rest-api/model"
	"time"

	"github.com/google/uuid"
)

type WaitlistRepository interface {
	AddUserToWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*model.WaitlistEntry, error)
	RemoveUserFromWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error
	GetWaitlistForEvent(ctx context.Context, eventID uuid.UUID) ([]model.WaitlistEntry, error)
	GetNextUserFromWaitlist(ctx context.Context, eventID uuid.UUID) (*model.WaitlistEntry, error)
	IsUserOnWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (bool, error)
}

type sqliteWaitlistRepository struct {
	db *sql.DB
}

func NewWaitlistRepository(db *sql.DB) WaitlistRepository {
	return &sqliteWaitlistRepository{db: db}
}

func (r *sqliteWaitlistRepository) AddUserToWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*model.WaitlistEntry, error) {
	entry := &model.WaitlistEntry{
		Id:      uuid.New(),
		EventID: eventID,
		UserID:  userID,
	}
	query := `
		INSERT INTO waitlist_entries (id, event_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, entry.Id, eventID, userID, now).Scan(&entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add user to waitlist: %w", err)
	}
	return entry, nil
}

func (r *sqliteWaitlistRepository) RemoveUserFromWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error {
	query := "DELETE FROM waitlist_entries WHERE event_id = $1 AND user_id = $2"
	res, err := r.db.ExecContext(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from waitlist: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after removing from waitlist: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found on waitlist or already removed") // Or sql.ErrNoRows if preferred
	}
	return nil
}

func (r *sqliteWaitlistRepository) GetWaitlistForEvent(ctx context.Context, eventID uuid.UUID) ([]model.WaitlistEntry, error) {
	query := `
		SELECT id, event_id, user_id, created_at
		FROM waitlist_entries
		WHERE event_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get waitlist for event: %w", err)
	}
	defer rows.Close()

	var entries []model.WaitlistEntry
	for rows.Next() {
		var entry model.WaitlistEntry
		if err := rows.Scan(&entry.Id, &entry.EventID, &entry.UserID, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan waitlist entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating waitlist entries: %w", err)
	}
	return entries, nil
}

func (r *sqliteWaitlistRepository) GetNextUserFromWaitlist(ctx context.Context, eventID uuid.UUID) (*model.WaitlistEntry, error) {
	query := `
		SELECT id, event_id, user_id, created_at
		FROM waitlist_entries
		WHERE event_id = $1
		ORDER BY created_at ASC
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, eventID)
	var entry model.WaitlistEntry
	err := row.Scan(&entry.Id, &entry.EventID, &entry.UserID, &entry.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No one on the waitlist
		}
		return nil, fmt.Errorf("failed to get next user from waitlist: %w", err)
	}
	return &entry, nil
}

func (r *sqliteWaitlistRepository) IsUserOnWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM waitlist_entries WHERE event_id = $1 AND user_id = $2)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, eventID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is on waitlist: %w", err)
	}
	return exists, nil
}
