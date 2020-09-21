package storage

import (
	"time"
)

type Event struct {
	ID             int
	Title          string
	StartTime      time.Time `db:"start_at"`
	Duration       time.Duration
	Description    string
	UserID         int `db:"user_id"`
	RemindDuration time.Duration
}

func (e Event) CopyFromIncoming(in IncomingEvent) Event {
	if !IsZeroValue(in.Title) {
		e.Title = in.Title
	}
	if !IsZeroValue(in.StartTime) {
		e.StartTime = in.StartTime
	}
	if !IsZeroValue(in.Duration) {
		e.Duration = in.Duration
	}
	if !IsZeroValue(in.Description) {
		e.Description = in.Description
	}
	if !IsZeroValue(in.RemindDuration) {
		e.RemindDuration = in.RemindDuration
	}

	return e
}

type IncomingEvent struct {
	Title          string
	StartTime      time.Time `db:"start_at"`
	Duration       time.Duration
	Description    string
	RemindDuration time.Duration
}
