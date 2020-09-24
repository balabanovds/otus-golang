package memorystorage

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

// NewTestStorage crates a storage of num Events with interval of 24h betwee,
// duration is 1 hour and ReminderDuration is 1 hour as well
func NewTestStorage(start time.Time, num int) storage.IStorage {
	return &Storage{
		data: generateTestData(start, num),
	}
}

// creates since start time num of events each num day with duration of 1 hour.
func generateTestData(start time.Time, num int) map[int]models.Event {
	m := make(map[int]models.Event)

	var dayNum = 0
	for i := 0; i < num; i++ {
		dayDuration, _ := time.ParseDuration("24h")
		ev := storage.NewTestEvent(start.Add(time.Duration(dayNum) * dayDuration))
		m[ev.ID] = ev
		dayNum++
	}
	return m
}
