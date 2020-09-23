package sqlstorage

import (
	"context"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" // import pgx
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	dsn        string
	db         *sqlx.DB
	eventsRepo storage.IEventStorage
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.dsn)

	return
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Events() storage.IEventStorage {
	if s.eventsRepo == nil {
		s.eventsRepo = newEventStorage(s)
	}
	return s.eventsRepo
}
