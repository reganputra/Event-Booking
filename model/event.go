package model

import "time"

type Event struct {
	Id          int
	Name        string    `bind:"required"`
	Description string    `bind:"required"`
	Location    string    `bind:"required"`
	Date        time.Time `bind:"required"`
	UserIds     int
}

var events []Event

func (e *Event) Save() {

	events = append(events, *e)
}

func GetAllEvents() []Event {
	return events
}
