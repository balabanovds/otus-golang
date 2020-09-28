package main

import (
	"context"
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/mq/mock"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	startDate, _ = time.Parse("2006-01-02", "2020-01-01")
	interval, _  = time.ParseDuration("24h")
)

func TestPublish(t *testing.T) {
	pub := mock.NewFakeQueue(5)

	st := memorystorage.NewTestStorage(startDate, 10)

	scheduler := new(pub, st, interval)
	scheduler.publishEvents(context.Background(), startDate)
	pub.Close()

	var list []models.MQNotification

	msgCh, _ := pub.Consume(nil)

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

	scheduler := new(nil, st, interval)
	scheduler.clearEvents(context.Background(), now)

	list = st.Events().ListBeforeDate(context.Background(), now)
	require.Len(t, list, 0)
}
