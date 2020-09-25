package models_test

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEvent_CopyFromIncoming(t *testing.T) {
	event := models.Event{
		Title: "title",
	}
	in := models.IncomingEvent{
		StartTime: time.Now(),
	}

	event.CopyFromIncoming(in)
	require.Equal(t, event.Title, "title")
}
