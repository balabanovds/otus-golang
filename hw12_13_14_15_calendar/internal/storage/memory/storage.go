package memorystorage

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	// TODO
	mu sync.RWMutex
}


func New() storage.IStorage {
	return &Storage{}
}


func (s *Storage) CreateEvent(event storage.Event) error {
	panic("implement me")
}

func (s *Storage) UpdateEvent(ID string, event storage.Event) error {
	panic("implement me")
}

func (s *Storage) DeleteEvent(ID string) {
	panic("implement me")
}

func (s *Storage) ListEventsForDay(date time.Time) []storage.Event {
	panic("implement me")
}

func (s *Storage) ListEventsForWeek(date time.Time) []storage.Event {
	panic("implement me")
}

func (s *Storage) ListEventsForMonth(date time.Time) []storage.Event {
	panic("implement me")
}
