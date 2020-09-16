package memorystorage

import (
	"context"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"sync"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu           sync.RWMutex
	data         map[int]models.Event
	eventStorage storage.IEventStorage
}

func New() storage.IStorage {
	return &Storage{
		data: make(map[int]models.Event),
	}
}

func (s *Storage) Events() storage.IEventStorage {
	if s.eventStorage == nil {
		s.eventStorage = newEventStorage(s)
	}
	return s.eventStorage
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}
