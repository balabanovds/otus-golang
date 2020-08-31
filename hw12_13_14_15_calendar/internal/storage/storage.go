package storage

import "time"

type IStorage interface {
	CreateEvent(event Event) error
	UpdateEvent(ID string, event Event) error
	DeleteEvent(ID string)
	ListEventsForDay(date time.Time) []Event
	ListEventsForWeek(date time.Time) []Event
	ListEventsForMonth(date time.Time) []Event
}
