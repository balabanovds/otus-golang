package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/jinzhu/now"
	"go.uber.org/zap"
)

type eventStorage struct {
	s *Storage
}

func newEventStorage(s *Storage) storage.IEventStorage {
	return &eventStorage{s}
}

func (e *eventStorage) Create(ctx context.Context, ev models.Event) (models.Event, error) {
	var cntr int
	err := e.s.db.GetContext(ctx, &cntr, "select count(*) from events where start_at < $1 and end_at > $1", ev.StartAt)
	if err != nil {
		return models.Event{}, err
	}
	if cntr > 0 {
		return models.Event{}, storage.ErrEventExists
	}

	var id int64
	err = e.s.db.QueryRowContext(ctx, "insert into events (title, start_at, end_at, description, user_id, remind_at) "+
		"values ($1, $2, $3, $4, $5, $6) returning id",
		ev.Title, ev.StartAt, ev.EndAt,
		ev.Description, ev.UserID, ev.RemindAt).Scan(&id)
	if err != nil {
		return models.Event{}, err
	}

	ev.ID = int(id)

	return ev, nil
}

func (e *eventStorage) Get(ctx context.Context, id int) (models.Event, error) {
	event := models.Event{}
	err := e.s.db.GetContext(ctx, &event, "select id, title, start_at, end_at, description, user_id, remind_at from events where id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Event{}, storage.ErrEvent404
		}
		return models.Event{}, err
	}
	return event, nil
}

func (e *eventStorage) Update(ctx context.Context, id int, event models.Event) error {
	_, err := e.s.db.NamedQueryContext(ctx, "update events set title = :title, start_at = :start_at, "+ //nolint:sqlclosecheck
		"end_at = :end_at, description = :descr, remind_at = :remind_at where id = :id",
		map[string]interface{}{
			"id":        id,
			"title":     event.Title,
			"start_at":  event.StartAt,
			"end_at":    event.EndAt,
			"descr":     event.Description,
			"remind_at": event.RemindAt,
		})
	return err
}

func (e *eventStorage) Delete(ctx context.Context, id int) {
	e.s.db.QueryRowxContext(ctx, "delete from events where id = $1", id)
}

func (e *eventStorage) ListForDay(ctx context.Context, date time.Time) models.EventsList {
	n := now.New(date)

	return e.filterEvents(ctx, n.BeginningOfDay(), n.EndOfDay())
}

func (e *eventStorage) ListForWeek(ctx context.Context, date time.Time) models.EventsList {
	n := now.New(date)

	return e.filterEvents(ctx, n.BeginningOfWeek(), n.EndOfWeek())
}

func (e *eventStorage) ListForMonth(ctx context.Context, date time.Time) models.EventsList {
	n := now.New(date)

	return e.filterEvents(ctx, n.BeginningOfMonth(), n.EndOfMonth())
}

func (e *eventStorage) ListBeforeDate(ctx context.Context, date time.Time) []models.Event {
	var events []models.Event
	err := e.s.db.SelectContext(ctx, &events, "select id, title, start_at, end_at, description, user_id, remind_at from events where start_at < $1", date)
	if err != nil {
		zap.L().Error("db: failed to get list of events", zap.Error(err))
	}

	return events
}

func (e *eventStorage) ListByReminderBetweenDates(ctx context.Context, startDate, endDate time.Time) []models.Event {
	var events []models.Event
	err := e.s.db.SelectContext(ctx, &events, "select id, title, start_at, end_at, description, user_id, remind_at from events where remind_at between $1 and $2", startDate, endDate)
	if err != nil {
		zap.L().Error("db: failed to get list of events", zap.Error(err))
	}

	return events
}

func (e *eventStorage) filterEvents(ctx context.Context, start, end time.Time) models.EventsList {
	var events []models.Event
	err := e.s.db.SelectContext(ctx, &events, "select id, title, start_at, end_at, description, user_id, remind_at from events "+
		"where start_at >= $1 and start_at < $2", start, end)
	if err != nil {
		zap.L().Error("db: failed to get list of events", zap.Error(err))
	}
	return models.NewEventsList(events)
}
