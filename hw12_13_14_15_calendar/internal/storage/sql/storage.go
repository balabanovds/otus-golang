package sqlstorage

import (
	"context"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
}

func New(config storage.Config) *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
