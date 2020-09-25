package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http/router"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	start = time.Date(2020, time.January, 1, 0, 0, 0, 1, time.Local)
)

func TestNotFound(t *testing.T) {
	srv := initSrv(0)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestListEvents(t *testing.T) {
	srv := initSrv(20)
	defer srv.Close()

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
		{"/year/2020/day/3", 1, http.StatusOK, true},
		{"/year/2020/day/300", 0, http.StatusOK, true},
		{"/year/2020/week/2", 7, http.StatusOK, true},
		{"/year/2020/week/20", 0, http.StatusOK, true},
		{"/year/2020/month/1", 20, http.StatusOK, true},
		{"/year/2020/month/11", 0, http.StatusOK, true},
	}

	for _, tst := range tests {
		t.Run(tst.path, func(t *testing.T) {
			path := "/events" + tst.path
			resp, err := http.Get(srv.URL + path)
			require.NoError(t, err)
			require.Equal(t, tst.status, resp.StatusCode)
			if tst.checkLen {
				require.Len(t, decodeEventsList(t, resp.Body), tst.length)
			}
		})
	}
}

func TestCreateEvent(t *testing.T) {
	srv := initSrv(0)
	defer srv.Close()

	respEvent := createEvent(t, srv.URL, start.AddDate(0, 0, 1))

	checkEventsForMonth(t, srv.URL, 2020, 1, respEvent)
}

func TestGetEvent(t *testing.T) {
	srv := initSrv(0)
	defer srv.Close()

	createdEvent := createEvent(t, srv.URL, start.AddDate(0, 0, 1))

	resp, err := http.Get(fmt.Sprintf("%s/event/%d", srv.URL, createdEvent.ID))
	require.NoError(t, err)

	require.Equal(t, createdEvent, decodeEvent(t, resp.Body))
}

func TestUpdateEvent(t *testing.T) {
	srv := initSrv(0)
	defer srv.Close()

	createdEvent := createEvent(t, srv.URL, start.AddDate(0, 0, 1))
	updEvent := models.IncomingEvent{
		Title: "upd",
	}

	data, err := json.Marshal(&updEvent)
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/event/%d", srv.URL, createdEvent.ID),
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
	srv := initSrv(0)
	defer srv.Close()

	createdEvent := createEvent(t, srv.URL, start.AddDate(0, 0, 1))
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/event/%d", srv.URL, createdEvent.ID),
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	checkEventsForMonth(t, srv.URL, 2020, 1, []models.Event{}...)
}

func initSrv(elements int) *httptest.Server {
	a := app.New(memorystorage.NewTestStorage(start, elements))

	r := router.New(a)
	return httptest.NewServer(r.RootHandler())
}

func createEvent(t *testing.T, url string, time time.Time) models.Event {
	incomingEvent := storage.NewTestIncomingEvent(time)
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

	return list.List
}

func decodeEvent(t *testing.T, r io.Reader) models.Event {
	var respEvent models.Event
	err := json.NewDecoder(r).Decode(&respEvent)
	require.NoError(t, err)
	return respEvent
}
