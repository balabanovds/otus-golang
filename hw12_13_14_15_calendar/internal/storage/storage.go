package storage

import (
	"context"
	"time"
)

type IStorage interface {
	Events() IEventStorage
	Connect(ctx context.Context) error
	Close() error
}

type IEventStorage interface {
	Create(ctx context.Context, event Event) (Event, error)
	Get(ctx context.Context, id int) (Event, error)
	Update(ctx context.Context, id int, event Event) error
	Delete(ctx context.Context, id int)
	ListForDay(ctx context.Context, date time.Time) []Event
	ListForWeek(ctx context.Context, date time.Time) []Event
	ListForMonth(ctx context.Context, date time.Time) []Event
}
