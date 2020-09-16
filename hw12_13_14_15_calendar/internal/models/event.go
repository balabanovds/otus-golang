package models

import (
	"time"
)

type Event struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	StartTime      time.Time     `db:"start_at" json:"start_time"`
	Duration       time.Duration `json:"duration"`
	Description    string        `json:"description"`
	UserID         int           `db:"user_id" json:"user_id"`
	RemindDuration time.Duration `json:"remind_duration"`
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