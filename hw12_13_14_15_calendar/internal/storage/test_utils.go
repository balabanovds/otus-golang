package storage

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
)

var id = 0

func NewTestEvent(start time.Time) models.Event {
	id++
	return models.Event{
		ID:          id,
		Title:       utils.RandString(10),
		StartAt:     start,
		EndAt:       start.Add(1 * time.Hour),
		Description: utils.RandString(30),
		UserID:      1,
		RemindAt:    start.Add(-1 * time.Hour),
	}
}
