package router_test

import (
	"encoding/json"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http/router"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNotFound(t *testing.T) {
	r := router.New(nil)
	srv := httptest.NewServer(r.RootHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestEventsHandler(t *testing.T) {
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local)
	a := app.New(memorystorage.NewTestStorage(start, 20))

	r := router.New(a)
	srv := httptest.NewServer(r.RootHandler())
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
				var list models.EventsList
				err := json.NewDecoder(resp.Body).Decode(&list)
				require.NoError(t, err)
				require.Len(t, list.List, tst.length)
			}
		})
	}
}
