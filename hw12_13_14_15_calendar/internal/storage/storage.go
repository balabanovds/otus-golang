package storage

import "time"

type IStorage interface {
	CreateEvent(event Event) (Event, error)
	GetEvent(id string) (Event, error)
	UpdateEvent(id string, event Event) error
	DeleteEvent(id string)
	ListEventsForDay(date time.Time) []Event
	ListEventsForWeek(date time.Time) []Event
	ListEventsForMonth(date time.Time) []Event
}
