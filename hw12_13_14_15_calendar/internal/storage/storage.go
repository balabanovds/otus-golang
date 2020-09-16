package storage

import (
	"context"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"time"
)

type IStorage interface {
	Events() IEventStorage
	Connect(ctx context.Context) error
	Close() error
}

type IEventStorage interface {
	Create(ctx context.Context, event models.Event) (models.Event, error)
	Get(ctx context.Context, id int) (models.Event, error)
	Update(ctx context.Context, id int, event models.Event) error
	Delete(ctx context.Context, id int)
	ListForDay(ctx context.Context, date time.Time) models.EventsList
	ListForWeek(ctx context.Context, date time.Time) models.EventsList
	ListForMonth(ctx context.Context, date time.Time) models.EventsList
}
