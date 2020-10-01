package models

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
)

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	StartAt     time.Time `db:"start_at" json:"start_time"`
	EndAt       time.Time `db:"end_at"`
	Description string    `json:"description"`
	UserID      int       `db:"user_id" json:"user_id"`
	RemindAt    time.Time `db:"remind_at"`
}

func (e *Event) CopyFromIncoming(in server.IncomingEvent) *Event {
	if in.Title != "" {
		e.Title = in.Title
	}
	if !in.StartAt.IsZero() {
		e.StartAt = in.StartAt
	}
	if in.Duration != 0 {
		e.EndAt = e.StartAt.Add(in.Duration)
	}
	if in.Description != "" {
		e.Description = in.Description
	}
	e.RemindAt = e.StartAt.Add(-in.RemindDuration)

	return e
}
