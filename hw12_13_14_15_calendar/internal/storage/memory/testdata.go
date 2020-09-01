package memorystorage

import (
	"math/rand"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

// creates since start time num of events each num day with duration of 1 hour.
func generateTestData(start time.Time, num int) map[string]storage.Event {
	m := make(map[string]storage.Event)

	var dayNum = 0
	for i := 0; i < num; i++ {
		dayDuration, _ := time.ParseDuration("24h")
		ev := newTestEvent(start.Add(time.Duration(dayNum) * dayDuration))
		m[ev.ID] = ev
		dayNum++
	}
	return m
}

func newTestEvent(start time.Time) storage.Event {
	return storage.Event{
		ID:             uuid.New().String(),
		Title:          randString(10),
		StartTime:      start,
		Duration:       1 * time.Hour,
		Description:    randString(30),
		UserID:         "user1",
		RemindDuration: 0,
	}
}

func randString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 _")

	b := make([]rune, length)
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
