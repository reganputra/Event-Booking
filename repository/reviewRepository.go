package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-rest-api/model"
)

type ReviewRepository interface {
	SaveReview(ctx context.Context, review *model.Review) error
	GetReviewsByEventID(ctx context.Context, eventID int64) ([]model.Review, error)
	GetReviewByEventAndUser(ctx context.Context, eventID int64, userID int64) (*model.Review, error)
	// Potentially: UpdateReview, DeleteReview if needed later
}

type sqliteReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &sqliteReviewRepository{db: db}
}

func (r *sqliteReviewRepository) SaveReview(ctx context.Context, review *model.Review) error {
	query := `
		INSERT INTO reviews (event_id, user_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for save review: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, review.EventID, review.UserID, review.Rating, review.Comment).Scan(&review.Id, &review.CreatedAt)
	if err != nil {
		// Consider checking for unique constraint violation error here if your DB driver supports it well
		return fmt.Errorf("failed to execute statement and scan for save review: %w", err)
	}
	return nil
}

func (r *sqliteReviewRepository) GetReviewsByEventID(ctx context.Context, eventID int64) ([]model.Review, error) {
	query := `
		SELECT id, event_id, user_id, rating, comment, created_at
		FROM reviews
		WHERE event_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews by event ID: %w", err)
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var review model.Review
		if err := rows.Scan(&review.Id, &review.EventID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan review row: %w", err)
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating review rows: %w", err)
	}
	return reviews, nil
}

func (r *sqliteReviewRepository) GetReviewByEventAndUser(ctx context.Context, eventID int64, userID int64) (*model.Review, error) {
	query := `
		SELECT id, event_id, user_id, rating, comment, created_at
		FROM reviews
		WHERE event_id = $1 AND user_id = $2
	`
	row := r.db.QueryRowContext(ctx, query, eventID, userID)
	var review model.Review
	err := row.Scan(&review.Id, &review.EventID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No review found is not an error in this specific query's context
		}
		return nil, fmt.Errorf("failed to scan review by event and user: %w", err)
	}
	return &review, nil
}
