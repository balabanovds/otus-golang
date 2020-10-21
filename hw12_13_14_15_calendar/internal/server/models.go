package server

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
)

type IncomingEvent struct {
	Title          string        `json:"title"`
	StartAt        time.Time     `db:"start_at" json:"start_time"`
	Duration       time.Duration `json:"duration"`
	Description    string        `json:"description"`
	RemindDuration time.Duration `json:"remind_duration"`
}

func NewTestIncomingEvent(start time.Time) IncomingEvent {
	return IncomingEvent{
		Title:          utils.RandString(10),
		StartAt:        start,
		Duration:       1 * time.Hour,
		Description:    utils.RandString(30),
		RemindDuration: 0,
	}
}
