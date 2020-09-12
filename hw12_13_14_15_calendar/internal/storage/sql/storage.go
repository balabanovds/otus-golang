package sqlstorage

import (
	"context"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint:golint
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	config     storage.Config
	db         *sqlx.DB
	eventsRepo storage.IEventStorage
}

func New(config storage.Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.config.Dsn)

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
