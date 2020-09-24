package sqlstorage

import (
	"context"
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" // import pgx
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Storage struct {
	dsn        string
	db         *sqlx.DB
	eventsRepo storage.IEventStorage
}

func New(cfg config.Storage) storage.IStorage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", s.dsn)
	if err != nil {
		return err
	}

	s.db = db
	zap.L().Info("connected to sql db")

	return nil
}

func (s *Storage) Close() error {
	zap.L().Info("close sql db")

	return s.db.Close()
}

func (s *Storage) String() string {
	return "sql storage"
}

func (s *Storage) Events() storage.IEventStorage {
	if s.eventsRepo == nil {
		s.eventsRepo = newEventStorage(s)
	}
	return s.eventsRepo
}
