// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"

	"github.com/stretchr/testify/require"
)

var start = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local)

func TestListEvents(t *testing.T) {
	err, cancel := s.initDB()
	require.NoError(t, err)
	defer cancel()
	defer utils.Close(s)

	tests := []struct {
		path     string
		length   int
		status   int
		checkLen bool
	}{
		{"/yea/2020/day/3", 0, http.StatusBadRequest, false},
		{"/year/2020a/day/3", 0, http.StatusBadRequest, false},
		{"/year/2020/day/3q", 0, http.StatusBadRequest, false},
		{"/year/2020/day/400", 0, http.StatusBadRequest, false},
		{"/year/2020/week/100", 0, http.StatusBadRequest, false},
		{"/year/2020/month/100", 0, http.StatusBadRequest, false},
		{"/year/2020/day/272", 2, http.StatusOK, true},
		{"/year/2020/day/300", 0, http.StatusOK, true},
		{"/year/2020/week/40", 3, http.StatusOK, true},
		{"/year/2020/week/20", 0, http.StatusOK, true},
		{"/year/2020/month/9", 3, http.StatusOK, true},
		{"/year/2020/month/11", 0, http.StatusOK, true},
	}

	for _, tst := range tests {
		t.Run(tst.path, func(t *testing.T) {
			path := "/events" + tst.path
			resp, err := http.Get(s.url() + path)
			require.NoError(t, err)
			require.Equal(t, tst.status, resp.StatusCode)
			if tst.checkLen {
				require.Len(t, decodeEventsList(t, resp.Body), tst.length)
			}
		})
	}
}

func TestCreateEvent(t *testing.T) {
	err, cancel := s.initDB()
	require.NoError(t, err)
	defer cancel()
	defer utils.Close(s)

	respEvent := createEvent(t, s.url(), start)

	checkEventsForMonth(t, s.url(), 2020, 1, respEvent)
}

func TestGetEvent(t *testing.T) {
	err, cancel := s.initDB()
	require.NoError(t, err)
	defer cancel()
	defer utils.Close(s)

	createdEvent := createEvent(t, s.url(), start.AddDate(0, 0, 1))

	resp, err := http.Get(fmt.Sprintf("%s/event/%d", s.url(), createdEvent.ID))
	require.NoError(t, err)

	require.Equal(t, createdEvent, decodeEvent(t, resp.Body))
}

func TestUpdateEvent(t *testing.T) {
	err, cancel := s.initDB()
	require.NoError(t, err)
	defer cancel()
	defer utils.Close(s)

	createdEvent := createEvent(t, s.url(), start.AddDate(0, 0, 1))
	updEvent := server.IncomingEvent{
		Title: "upd",
	}

	data, err := json.Marshal(&updEvent)
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/event/%d", s.url(), createdEvent.ID),
		bytes.NewReader(data),
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	createdEvent.Title = "upd"
	require.Equal(t, createdEvent, decodeEvent(t, resp.Body))
}

func TestDeleteEvent(t *testing.T) {
	err, cancel := s.initDB()
	require.NoError(t, err)
	defer cancel()
	defer utils.Close(s)

	createdEvent := createEvent(t, s.url(), start.AddDate(0, 0, 1))
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/event/%d", s.url(), createdEvent.ID),
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	checkEventsForMonth(t, s.url(), 2020, 1, []models.Event{}...)
}

func createEvent(t *testing.T, url string, time time.Time) models.Event {
	incomingEvent := server.NewTestIncomingEvent(time)
	data, err := json.Marshal(&incomingEvent)
	require.NoError(t, err)

	resp, err := http.Post(url+"/event", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return decodeEvent(t, resp.Body)
}

func checkEventsForMonth(t *testing.T, url string, year, month int, events ...models.Event) {
	resp, err := http.Get(fmt.Sprintf("%s/events/year/%d/month/%d", url, year, month))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, events, decodeEventsList(t, resp.Body))
}

func decodeEventsList(t *testing.T, r io.Reader) []models.Event {
	var list models.EventsList
	err := json.NewDecoder(r).Decode(&list)
	require.NoError(t, err)

	events := []models.Event{}
	if list.Len > 0 {
		events = list.List
	}

	return events
}

func decodeEvent(t *testing.T, r io.Reader) models.Event {
	var respEvent models.Event
	err := json.NewDecoder(r).Decode(&respEvent)
	require.NoError(t, err)
	return respEvent
}
