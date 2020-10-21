package memorystorage

import (
	"context"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type eventStorage struct {
	st *Storage
}

func newEventStorage(st *Storage) *eventStorage {
	return &eventStorage{st}
}

var id = 0

func nextID() int {
	id++
	return id
}

func (s *eventStorage) Create(_ context.Context, event models.Event) (models.Event, error) {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	for _, ev := range s.st.data {
		if ev.StartAt.UnixNano() < event.StartAt.UnixNano() &&
			ev.EndAt.UnixNano() > event.StartAt.UnixNano() {
			return models.Event{}, storage.ErrEventExists
		}
	}

	event.ID = nextID()
	s.st.data[event.ID] = event

	return event, nil
}

func (s *eventStorage) Get(_ context.Context, id int) (models.Event, error) {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	if _, ok := s.st.data[id]; !ok {
		return models.Event{}, storage.ErrEvent404
	}
	return s.st.data[id], nil
}

func (s *eventStorage) Update(_ context.Context, id int, event models.Event) error {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	_, ok := s.st.data[id]
	if !ok {
		return storage.ErrEvent404
	}

	s.st.data[id] = event

	return nil
}

func (s *eventStorage) Delete(_ context.Context, id int) {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	delete(s.st.data, id)
}

func (s *eventStorage) ListForDay(_ context.Context, date time.Time) models.EventsList {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		return newTime.YearDay() == exTime.YearDay()
	})
}

func (s *eventStorage) ListForWeek(_ context.Context, date time.Time) models.EventsList {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		origYear, origWeek := newTime.ISOWeek()
		destYear, destWeek := exTime.ISOWeek()
		return origYear == destYear && origWeek == destWeek
	})
}

func (s *eventStorage) ListForMonth(_ context.Context, date time.Time) models.EventsList {
	return s.filterEvents(date, func(newTime time.Time, exTime time.Time) bool {
		return newTime.Month() == exTime.Month()
	})
}

func (s *eventStorage) ListBeforeDate(_ context.Context, date time.Time) []models.Event {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	var list []models.Event
	for _, ev := range s.st.data {
		if ev.StartAt.UnixNano() < date.UnixNano() {
			list = append(list, ev)
		}
	}

	return list
}

func (s *eventStorage) ListByReminderBetweenDates(_ context.Context, startDate, endDate time.Time) []models.Event {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	var list []models.Event
	for _, ev := range s.st.data {
		remindDate := ev.RemindAt
		if remindDate.UnixNano() > startDate.UnixNano() && remindDate.UnixNano() < endDate.UnixNano() {
			list = append(list, ev)
		}
	}

	return list
}

func (s *eventStorage) filterEvents(
	date time.Time,
	cmp func(newTime time.Time, exTime time.Time) bool,
) models.EventsList {
	s.st.mu.Lock()
	defer s.st.mu.Unlock()

	list := make([]models.Event, 0)

	for _, ev := range s.st.data {
		if cmp(date, ev.StartAt) {
			list = append(list, ev)
		}
	}

	return models.NewEventsList(list)
}
