package model

import (
	"time"
)

type Event struct {
	Id          int64
	Name        string    `bind:"required"`
	Description string    `bind:"required"`
	Location    string    `bind:"required"`
	Date        time.Time `bind:"required"`
	Category    string    `bind:"required"`
	UserIds     int64
	AverageRating float64 `json:"average_rating,omitempty"` // omitempty so it doesn't show if 0
	Capacity    int       `json:"capacity,omitempty" binding:"gte=0"` // omitempty so not required on update if not changing
	// RegisteredCount will be dynamically determined, not stored directly in DB or handled by simple binding.
}
