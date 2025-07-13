package model

import (
	"time"

	"github.com/google/uuid"
)

type WaitlistEntry struct {
	Id        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id" binding:"required"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// You might also want to extend the Event model to include Capacity
// If Event model is in a different file, ensure it's updated accordingly.
// For example, in model/event.go:
// type Event struct {
//     ...
//     Capacity          int       `json:"capacity" binding:"gte=0"` // Maximum number of attendees
//     RegisteredCount   int       `json:"registered_count"` // Current number of registered attendees (calculated)
//     ...
// }
