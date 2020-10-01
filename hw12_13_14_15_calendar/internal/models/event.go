package models

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
)

type Event struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	StartTime      time.Time     `db:"start_at" json:"start_time"`
	EndTime        time.Time     `db:"end_at"`
	Duration       time.Duration `json:"duration"`
	Description    string        `json:"description"`
	UserID         int           `db:"user_id" json:"user_id"`
	RemindDuration time.Duration `json:"remind_duration"`
	RemindAt       time.Time     `db:"remind_at"`
}

func (e *Event) CopyFromIncoming(in server.IncomingEvent) *Event {
	if in.Title != "" {
		e.Title = in.Title
	}
	if !in.StartTime.IsZero() {
		e.StartTime = in.StartTime
	}
	if in.Duration != 0 {
		e.Duration = in.Duration
	}
	if in.Description != "" {
		e.Description = in.Description
	}
	if in.RemindDuration != 0 {
		e.RemindDuration = in.RemindDuration
	}

	return e
}
