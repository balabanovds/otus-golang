package storage

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"math/rand"
	"time"
)

var id = 0

func NewTestEvent(start time.Time) models.Event {
	id++
	return models.Event{
		ID:             id,
		Title:          randString(10),
		StartTime:      start,
		Duration:       1 * time.Hour,
		Description:    randString(30),
		UserID:         1,
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
