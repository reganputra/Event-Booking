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
	UserIds     int64
}
