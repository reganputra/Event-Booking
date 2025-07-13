package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id            uuid.UUID  `json:"id"`
	Name          *string    `json:"name,omitempty" binding:"omitempty,min=5"`
	Description   *string    `json:"description,omitempty" binding:"omitempty,min=10"`
	Location      *string    `json:"location,omitempty"`
	Date          *time.Time `json:"date,omitempty"`
	Category      *string    `json:"category,omitempty"`
	UserIds       uuid.UUID  `json:"user_id"`
	AverageRating float64    `json:"average_rating,omitempty"`
	Capacity      *int       `json:"capacity,omitempty" binding:"omitempty,gte=0"`
}
