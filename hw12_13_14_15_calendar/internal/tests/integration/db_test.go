// +build integration

package integration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
	_ "github.com/jackc/pgx/v4/stdlib" // import pgx

	"github.com/jmoiron/sqlx"
)

func (s *suite) initDB() (error, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := s.open(ctx); err != nil {
		return err, nil
	}
	if err := s.seedDB(ctx); err != nil {
		return err, nil
	}

	return nil, cancel
}

func (s *suite) open(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.cfg.Storage.Host, s.cfg.Storage.Port, s.cfg.Storage.User, s.cfg.Storage.Password, s.cfg.Storage.DBName,
	)
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *suite) seedDB(ctx context.Context) error {
	events := []models.Event{
		mustNewEvent("2020-09-28 12:00", 2, 1),
		mustNewEvent("2020-09-28 18:00", 1, 1),
		mustNewEvent("2020-09-30 12:00", 3, 1),
	}

	for _, ev := range events {
		e, err := s.storage.Events().Create(ctx, ev)
		if err != nil {
			return err
		}
		s.events = append(s.events, e)
	}

	return nil
}

func (s *suite) truncate(tables ...string) error {
	_, err := s.db.Exec(fmt.Sprintf("truncate table %s cascade", strings.Join(tables, ", ")))
	return err
}

func (s *suite) Close() error {
	if err := s.truncate("events"); err != nil {
		return err
	}
	return s.db.Close()
}

// mustNewEvent accepts startTime in format '2006-01-02 15:04'
func mustNewEvent(startTime string, durHrs, remindDurHrs int) models.Event {
	t, err := time.Parse("2006-01-02 15:04", startTime)
	if err != nil {
		panic(err)
	}
	return newEvent(t, durHrs, remindDurHrs)
}

func newEvent(startTime time.Time, durHrs, remindDurHrs int) models.Event {
	return models.Event{
		Title:       utils.RandString(10),
		StartAt:     startTime,
		EndAt:       startTime.Add(time.Duration(durHrs) * time.Hour),
		Description: utils.RandString(30),
		UserID:      1,
		RemindAt:    startTime.Add(time.Duration(-remindDurHrs) * time.Hour),
	}
}
