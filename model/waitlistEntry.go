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

// TODO feature: Extend the Event model to include a Capacity field
