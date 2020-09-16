package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"time"

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
	err := e.s.db.GetContext(ctx, &cntr, "select count(*) from events where start_at < $1 and end_at > $1", ev.StartTime)
	if err != nil {
		return models.Event{}, err
	}
	if cntr > 0 {
		return models.Event{}, storage.ErrEventExists
	}

	res, err := e.s.db.NamedExecContext(ctx, "insert into events (title, start_at, end_at, description, user_id, remind_at) "+
		"values (:title, :start_at, :end_at, :descr, :uid, :remind_at) returning id",
		map[string]interface{}{
			"title":     ev.Title,
			"start_at":  ev.StartTime,
			"end_at":    ev.StartTime.Add(ev.Duration),
			"descr":     ev.Description,
			"uid":       ev.UserID,
			"remind_at": ev.StartTime.Add(-ev.RemindDuration),
		})
	if err != nil {
		return models.Event{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ev, err
	}

	ev.ID = int(id)

	return ev, nil
}

func (e *eventStorage) Get(ctx context.Context, id int) (models.Event, error) {
	event := models.Event{}
	err := e.s.db.GetContext(ctx, &event, "select * from events where id = $1", id)
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
			"start_at":  event.StartTime,
			"end_at":    event.StartTime.Add(event.Duration),
			"descr":     event.Description,
			"remind_at": event.StartTime.Add(-event.RemindDuration),
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

func (e *eventStorage) filterEvents(ctx context.Context, start, end time.Time) models.EventsList {
	var events []models.Event
	err := e.s.db.SelectContext(ctx, &events, "select * from events where start_at > $1 and start_at < $2", start, end)
	if err != nil {
		zap.L().Error("db: failed to get list of events")
	}
	return models.NewEventsList(events)
}
