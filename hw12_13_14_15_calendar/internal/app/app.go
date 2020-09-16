package app

import (
	"context"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type ctxKey uint8

const (
	CtxKeyUserID ctxKey = iota
)

type ListFunc func(ctx context.Context, year, day int) (models.EventsList, error)

type Application interface {
	CreateEvent(ctx context.Context, event models.IncomingEvent) error
	Get(ctx context.Context, id int) (models.Event, error)
	Update(ctx context.Context, id int, event models.IncomingEvent) error
	Delete(ctx context.Context, id int) error
	ListForDay(ctx context.Context, year, day int) (models.EventsList, error)
	ListForWeek(ctx context.Context, year, week int) (models.EventsList, error)
	ListForMonth(ctx context.Context, year, month int) (models.EventsList, error)
}

type App struct {
	storage storage.IStorage
}

func New(storage storage.IStorage) *App {
	return &App{storage}
}

func (a *App) CreateEvent(ctx context.Context, in models.IncomingEvent) error {
	ctxID, ok := ctx.Value(CtxKeyUserID).(int)
	if !ok {
		return ErrAppGeneral
	}
	var ev models.Event
	ev.CopyFromIncoming(in)
	ev.UserID = ctxID

	_, err := a.storage.Events().Create(ctx, ev)

	return err
}

func (a *App) Get(ctx context.Context, id int) (models.Event, error) {
	return a.storage.Events().Get(ctx, id)
}

func (a *App) Update(ctx context.Context, id int, in models.IncomingEvent) error {
	ev, err := a.Get(ctx, id)
	if err != nil {
		return err
	}
	ev.CopyFromIncoming(in)
	return a.storage.Events().Update(ctx, id, ev)
}

func (a *App) Delete(ctx context.Context, id int) error {
	if err := a.checkUserID(ctx, id); err != nil {
		return err
	}
	a.storage.Events().Delete(ctx, id)

	return nil
}

func (a *App) ListForDay(ctx context.Context, year, day int) (models.EventsList, error) {
	yearEnd := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)
	if day > yearEnd.YearDay() {
		return models.NewEventsList(nil), ErrDateFormat
	}
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	date := yearStart.AddDate(0, 0, day)

	return a.storage.Events().ListForDay(ctx, date), nil
}

func (a *App) ListForWeek(ctx context.Context, year, week int) (models.EventsList, error) {
	yearEnd := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)
	_, maxWeek := yearEnd.ISOWeek()
	if week > maxWeek {
		return models.NewEventsList(nil), ErrDateFormat
	}
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	date := yearStart.Add(time.Duration(week*7*24) * time.Hour)

	return a.storage.Events().ListForWeek(ctx, date), nil
}

func (a *App) ListForMonth(ctx context.Context, year, month int) (models.EventsList, error) {
	if month > 12 {
		return models.NewEventsList(nil), ErrDateFormat
	}
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	date := yearStart.AddDate(0, month-1, 0)

	return a.storage.Events().ListForMonth(ctx, date), nil
}

func (a *App) checkUserID(ctx context.Context, eventID int) error {
	ev, err := a.Get(ctx, eventID)
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
