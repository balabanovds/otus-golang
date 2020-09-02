package memorystorage

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type eventsRepo struct {
	st *Storage
}

func newEventStore(st *Storage) *eventsRepo {
	return &eventsRepo{st}
}

func (s *eventsRepo) Create(event storage.Event) (storage.Event, error) {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	for _, ev := range s.st.data {
		if ev.StartTime.UnixNano() < event.StartTime.UnixNano() &&
			ev.StartTime.Add(ev.Duration).UnixNano() > event.StartTime.UnixNano() {
			return storage.Event{}, storage.ErrEventExists
		}
	}

	s.st.data[event.ID] = event

	return event, nil
}

func (s *eventsRepo) Get(id string) (storage.Event, error) {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	if _, ok := s.st.data[id]; !ok {
		return storage.Event{}, storage.ErrEvent404
	}
	return s.st.data[id], nil
}

func (s *eventsRepo) Update(id string, event storage.Event) error {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	if _, ok := s.st.data[id]; !ok {
		return storage.ErrEvent404
	}
	s.st.data[id] = event

	return nil
}

func (s *eventsRepo) Delete(id string) {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	delete(s.st.data, id)
}

func (s *eventsRepo) ListForDay(date time.Time) []storage.Event {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		return newTime.YearDay() == exTime.YearDay()
	})
}

func (s *eventsRepo) ListForWeek(date time.Time) []storage.Event {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		origYear, origWeek := newTime.ISOWeek()
		destYear, destWeek := exTime.ISOWeek()
		return origYear == destYear && origWeek == destWeek
	})
}

func (s *eventsRepo) ListForMonth(date time.Time) []storage.Event {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		return newTime.Month() == exTime.Month()
	})
}

func (s *eventsRepo) filterEvents(date time.Time,
	cmp func(newTime time.Time, exTime time.Time) bool) []storage.Event {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	list := make([]storage.Event, 0)

	for _, ev := range s.st.data {
		if cmp(date, ev.StartTime) {
			list = append(list, ev)
		}
	}

	return list
}
