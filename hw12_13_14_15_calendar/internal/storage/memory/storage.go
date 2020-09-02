package memorystorage

import (
	"context"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	mu         sync.RWMutex
	data       map[string]storage.Event
	eventsRepo *eventsRepo
}

func New() storage.Repo {
	return &Storage{
		data: make(map[string]storage.Event),
	}
}

func (s *Storage) Events() storage.EventsRepo {
	if s.eventsRepo == nil {
		s.eventsRepo = newEventStore(s)
	}
	return s.eventsRepo
}

func (s *Storage) Open(ctx context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}
