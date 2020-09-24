package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/amqp/mock"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	startDate, _ = time.Parse("2006-01-02", "2020-01-01")
	interval, _  = time.ParseDuration("24h")
)

func TestPublish(t *testing.T) {
	q := mock.NewFakeQueue(5)

	st := memorystorage.NewTestStorage(startDate, 10)

	scheduler := New(q, st, interval)
	scheduler.publishEvents(context.Background(), startDate)
	q.Close()

	var list []models.MQNotification

	msgCh, err := q.Consume()
	require.NoError(t, err)

	for data := range msgCh {
		list = append(list, data)
	}

	require.Len(t, list, 1)
}

func TestClearOldEvents(t *testing.T) {
	now := time.Now()
	start := now.AddDate(-1, -1, 0)
	st := memorystorage.NewTestStorage(start, 10)

	list := st.Events().ListBeforeDate(context.Background(), now)
	require.Len(t, list, 10)

	scheduler := New(nil, st, interval)
	scheduler.clearEvents(context.Background(), now)

	list = st.Events().ListBeforeDate(context.Background(), now)
	require.Len(t, list, 0)
}
