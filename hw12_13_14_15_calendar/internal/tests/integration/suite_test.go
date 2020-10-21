// +build integration

package integration

import (
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"

	"github.com/jmoiron/sqlx"
)

type suite struct {
	cfg     config
	db      *sqlx.DB
	storage storage.IStorage
	events  []models.Event // events present in db
}

func newSuite(cfg config, storage storage.IStorage) *suite {
	return &suite{
		cfg:     cfg,
		storage: storage,
	}
}

func (s *suite) url() string {
	return fmt.Sprintf("http://%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)
}
