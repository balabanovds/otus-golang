package app

import (
	"context"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type ctxKey uint8

const (
	CtxKeyUserID ctxKey = iota
)

type App struct {
	storage storage.IStorage
}

func New(storage storage.IStorage) *App {
	return &App{storage}
}

func (a *App) CreateEvent(ctx context.Context, in storage.IncomingEvent) error {
	ctxID, ok := ctx.Value(CtxKeyUserID).(int)
	if !ok {
		return ErrAppGeneral
	}
	var ev storage.Event
	ev.CopyFromIncoming(in)
	ev.UserID = ctxID

	_, err := a.storage.Events().Create(ctx, ev)

	return err
}

func (a *App) Get(ctx context.Context, id int) (storage.Event, error) {
	return a.storage.Events().Get(ctx, id)
}

func (a *App) Update(ctx context.Context, id int, in storage.IncomingEvent) error {
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

func (a *App) ListForDay(ctx context.Context, date time.Time) []storage.Event {
	return a.storage.Events().ListForDay(ctx, date)
}

func (a *App) ListForWeek(ctx context.Context, date time.Time) []storage.Event {
	return a.storage.Events().ListForWeek(ctx, date)
}

func (a *App) ListForMonth(ctx context.Context, date time.Time) []storage.Event {
	return a.storage.Events().ListForMonth(ctx, date)
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
