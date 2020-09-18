package app_test

import (
	"context"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestApp_CreateEvent(t *testing.T) {
	a := app.New(memorystorage.NewTestStorage(time.Now(), 0))

	ctx := context.WithValue(context.Background(), app.CtxKeyUserID, 1)
	inEvent := models.IncomingEvent{
		Title:          "test title",
		StartTime:      time.Now(),
		Duration:       300,
		Description:    "",
		RemindDuration: 0,
	}
	createdEvent, err := a.CreateEvent(ctx, inEvent)
	require.NoError(t, err)

	gotEvent, err := a.Get(ctx, createdEvent.ID)
	require.NoError(t, err)
	require.Equal(t, createdEvent, gotEvent)
}

func TestApp_List(t *testing.T) {
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local)
	a := app.New(memorystorage.NewTestStorage(start, 20))

	tests := []struct {
		name string
		fn   app.ListFunc
		year int
		val  int
		err  error
		len  int
	}{
		{"day", a.ListForDay, 2020, 3, nil, 1},
		{"day empty", a.ListForDay, 2020, 300, nil, 0},
		{"day error", a.ListForDay, 2020, 400, app.ErrDateFormat, 0},
		{"week", a.ListForWeek, 2020, 2, nil, 7},
		{"week empty", a.ListForWeek, 2020, 20, nil, 0},
		{"week error", a.ListForWeek, 2020, 200, app.ErrDateFormat, 0},
		{"month", a.ListForMonth, 2020, 1, nil, 20},
		{"month empty", a.ListForMonth, 2020, 10, nil, 0},
		{"month error", a.ListForMonth, 2020, 100, app.ErrDateFormat, 0},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			list, err := tst.fn(context.Background(), tst.year, tst.val)
			if tst.err != nil {
				require.Error(t, err)
				require.EqualError(t, tst.err, err.Error())
				return
			}
			require.NoError(t, err)
			require.Len(t, list.List, tst.len)
		})
	}
}
