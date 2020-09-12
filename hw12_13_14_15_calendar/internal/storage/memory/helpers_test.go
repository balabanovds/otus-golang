package memorystorage

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

// creates since start time num of events each num day with duration of 1 hour.
func generateTestData(start time.Time, num int) map[int]storage.Event {
	m := make(map[int]storage.Event)

	var dayNum = 0
	for i := 0; i < num; i++ {
		dayDuration, _ := time.ParseDuration("24h")
		ev := storage.NewTestEvent(start.Add(time.Duration(dayNum) * dayDuration))
		m[ev.ID] = ev
		dayNum++
	}
	return m
}
