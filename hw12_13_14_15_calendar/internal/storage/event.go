package storage

import "time"

type Event struct {
	ID             string
	Title          string
	StartTime      time.Time
	Duration       time.Duration
	Description    string
	UserID         string
	RemindDuration time.Duration
}
