package model

import (
	"time"
)

type Event struct {
	Id            int64
	Name          string    `binding:"required,min=5"`
	Description   string    `binding:"required,min=10"`
	Location      string    `binding:"required"`
	Date          time.Time `binding:"required"`
	Category      string    `binding:"required"`
	UserIds       int64
	AverageRating float64 `json:"average_rating,omitempty"`
	Capacity      int     `json:"capacity,omitempty" binding:"gte=0"`
	// RegisteredCount will be dynamically determined, not stored directly in DB or handled by simple binding.
}
