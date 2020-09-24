package storage

import (
	"context"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
)

type IStorage interface {
	Events() IEventStorage
	Connect(ctx context.Context) error
	utils.CloseStringer
}

type IEventStorage interface {
	Create(ctx context.Context, event models.Event) (models.Event, error)
	Get(ctx context.Context, id int) (models.Event, error)
	Update(ctx context.Context, id int, event models.Event) error
	Delete(ctx context.Context, id int)
	ListForDay(ctx context.Context, date time.Time) models.EventsList
	ListForWeek(ctx context.Context, date time.Time) models.EventsList
	ListForMonth(ctx context.Context, date time.Time) models.EventsList
	ListBeforeDate(ctx context.Context, date time.Time) []models.Event
	ListByReminderBetweenDates(ctx context.Context, startDate, endDate time.Time) []models.Event
}
