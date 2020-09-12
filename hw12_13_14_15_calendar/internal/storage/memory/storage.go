package memorystorage

import (
	"context"
	"sync"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu           sync.RWMutex
	data         map[int]storage.Event
	eventStorage storage.IEventStorage
}

func New() storage.IStorage {
	return &Storage{
		data: make(map[int]storage.Event),
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
