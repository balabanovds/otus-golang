package app

import (
	"context"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type ctxKey uint8

const (
	CtxKeyUserID ctxKey = iota
)

type ListFunc func(ctx context.Context, year, day int) (models.EventsList, error)

type Application interface {
	CreateEvent(ctx context.Context, event models.IncomingEvent) (models.Event, error)
	GetEvent(ctx context.Context, id int) (models.Event, error)
	UpdateEvent(ctx context.Context, id int, event models.IncomingEvent) (models.Event, error)
	DeleteEvent(ctx context.Context, id int) error
	EventListForDay(ctx context.Context, year, day int) (models.EventsList, error)
	EventListForWeek(ctx context.Context, year, week int) (models.EventsList, error)
	EventListForMonth(ctx context.Context, year, month int) (models.EventsList, error)
}

type App struct {
	storage storage.IStorage
}

func New(storage storage.IStorage) *App {
	return &App{storage}
}

func (a *App) CreateEvent(ctx context.Context, in models.IncomingEvent) (models.Event, error) {
	ctxID, ok := ctx.Value(CtxKeyUserID).(int)
	if !ok {
		return models.Event{}, ErrAppGeneral
	}
	var ev models.Event
	ev.CopyFromIncoming(in)
	ev.UserID = ctxID

	return a.storage.Events().Create(ctx, ev)
}

func (a *App) GetEvent(ctx context.Context, id int) (models.Event, error) {
	return a.storage.Events().Get(ctx, id)
}

func (a *App) UpdateEvent(ctx context.Context, id int, in models.IncomingEvent) (models.Event, error) {
	if err := a.checkUserID(ctx, id); err != nil {
		return models.Event{}, err
	}
	ev, err := a.GetEvent(ctx, id)
	if err != nil {
		return models.Event{}, err
	}
	ev.CopyFromIncoming(in)
	err = a.storage.Events().Update(ctx, id, ev)
	return ev, err
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	if err := a.checkUserID(ctx, id); err != nil {
		return err
	}
	a.storage.Events().Delete(ctx, id)

	return nil
}

func (a *App) EventListForDay(ctx context.Context, year, day int) (models.EventsList, error) {
	yearEnd := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)
	if day > yearEnd.YearDay() {
		return models.NewEventsList(nil), ErrDateFormat
	}
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	date := yearStart.AddDate(0, 0, day)

	return a.storage.Events().ListForDay(ctx, date), nil
}

func (a *App) EventListForWeek(ctx context.Context, year, week int) (models.EventsList, error) {
	yearEnd := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)
	_, maxWeek := yearEnd.ISOWeek()
	if week > maxWeek {
		return models.NewEventsList(nil), ErrDateFormat
	}
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	date := yearStart.Add(time.Duration(week*7*24) * time.Hour)

	return a.storage.Events().ListForWeek(ctx, date), nil
}

func (a *App) EventListForMonth(ctx context.Context, year, month int) (models.EventsList, error) {
	if month > 12 {
		return models.NewEventsList(nil), ErrDateFormat
	}
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	date := yearStart.AddDate(0, month-1, 0)

	return a.storage.Events().ListForMonth(ctx, date), nil
}

func (a *App) checkUserID(ctx context.Context, eventID int) error {
	ev, err := a.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}

	ctxID, ok := ctx.Value(CtxKeyUserID).(int)
	if !ok {
		return ErrAppGeneral
	}
	if ev.UserID != ctxID {
		return ErrForbidden
	}

	return nil
}
