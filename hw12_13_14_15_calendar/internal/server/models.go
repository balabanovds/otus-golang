package server

import (
	"math/rand"
	"time"
)

type IncomingEvent struct {
	Title          string        `json:"title"`
	StartTime      time.Time     `db:"start_at" json:"start_time"`
	Duration       time.Duration `json:"duration"`
	Description    string        `json:"description"`
	RemindDuration time.Duration `json:"remind_duration"`
}

func NewTestIncomingEvent(start time.Time) IncomingEvent {
	return IncomingEvent{
		Title:          RandString(10),
		StartTime:      start,
		Duration:       1 * time.Hour,
		Description:    RandString(30),
		RemindDuration: 0,
	}
}

func RandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 _")

	b := make([]rune, length)
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
