package eventing

import "time"

type Event struct {
	CreatedAt   time.Time
	Description string
}

func NewEvent(description string) *Event {
	return &Event{
		CreatedAt:   time.Now(),
		Description: description,
	}
}
