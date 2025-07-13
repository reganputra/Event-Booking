package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	Id        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating" binding:"required,gte=1,lte=5"`
	Comment   string    `json:"comment" binding:"required,min=10"`
	CreatedAt time.Time `json:"created_at"`
}
