package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `binding:"required,min=5"`
	Description   string    `binding:"required,min=10"`
	Location      string    `binding:"required"`
	Date          time.Time `binding:"required"`
	Category      string    `binding:"required"`
	UserIds       uuid.UUID `json:"user_id"`
	AverageRating float64   `json:"average_rating,omitempty"`
	Capacity      int       `json:"capacity,omitempty" binding:"gte=0"`
	// RegisteredCount will be dynamically determined, not stored directly in DB or handled by simple binding.
}
