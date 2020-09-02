package storage

import (
	"context"
	"time"
)

type Repo interface {
	Events() EventsRepo
	Open(ctx context.Context) error
	Close() error
}

type EventsRepo interface {
	Create(event Event) (Event, error)
	Get(id string) (Event, error)
	Update(id string, event Event) error
	Delete(id string)
	ListForDay(date time.Time) []Event
	ListForWeek(date time.Time) []Event
	ListForMonth(date time.Time) []Event
}
