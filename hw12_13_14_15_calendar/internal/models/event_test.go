package models_test

import (
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/stretchr/testify/require"
)

func TestEvent_CopyFromIncoming(t *testing.T) {
	event := models.Event{
		Title: "title",
	}
	in := server.IncomingEvent{
		StartAt: time.Now(),
	}

	event.CopyFromIncoming(in)
	require.Equal(t, event.Title, "title")
}
