package memorystorage

import (
	"sync"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	data map[string]storage.Event
	mu   sync.RWMutex
}

func New() storage.IStorage {
	return &Storage{
		data: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ev := range s.data {
		if ev.StartTime.UnixNano() < event.StartTime.UnixNano() &&
			ev.StartTime.Add(ev.Duration).UnixNano() > event.StartTime.UnixNano() {
			return storage.Event{}, storage.ErrEventExists
		}
	}

	s.data[event.ID] = event

	return event, nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return storage.Event{}, storage.ErrEvent404
	}
	return s.data[id], nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return storage.ErrEvent404
	}
	s.data[id] = event

	return nil
}

func (s *Storage) DeleteEvent(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
}

func (s *Storage) ListEventsForDay(date time.Time) []storage.Event {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		return newTime.YearDay() == exTime.YearDay()
	})
}

func (s *Storage) ListEventsForWeek(date time.Time) []storage.Event {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		origYear, origWeek := newTime.ISOWeek()
		destYear, destWeek := exTime.ISOWeek()
		return origYear == destYear && origWeek == destWeek
	})
}

func (s *Storage) ListEventsForMonth(date time.Time) []storage.Event {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		return newTime.Month() == exTime.Month()
	})
}

func (s *Storage) filterEvents(date time.Time,
	cmp func(newTime time.Time, exTime time.Time) bool) []storage.Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := make([]storage.Event, 0)

	for _, ev := range s.data {
		if cmp(date, ev.StartTime) {
			list = append(list, ev)
		}
	}

	return list
}
