package model

import (
	"time"
)

type Review struct {
	Id        int64     `json:"id"`
	EventID   int64     `json:"event_id" binding:"required"`
	UserID    int64     `json:"user_id"` // Should be set from authenticated user
	Rating    int       `json:"rating" binding:"required,gte=1,lte=5"` // Rating from 1 to 5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
