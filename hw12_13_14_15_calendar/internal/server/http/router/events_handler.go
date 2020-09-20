package router

import (
	"net/http"

	a "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
)

type eventsHandler struct {
	app a.Application
}

func newEventsHandler(app a.Application) eventsHandler {
	return eventsHandler{app}
}

/*
	/year/:year_num/day/:day_in_year_number
	/year/:year_num/week/:week_number
	/year/:year_num/month/:month_number
*/

func (h *eventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		clientError(w, r, http.StatusMethodNotAllowed, ErrWrongMethod)
		return
	}

	var head string
	head, r.URL.Path = splitPath(r.URL.Path)

	if head != "year" {
		clientError(w, r, http.StatusBadRequest, ErrWrongPath)
		return
	}

	var (
		year int
		err  error
	)

	year, r.URL.Path, err = parseInt(r.URL.Path)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, ErrWrongPath)
		return
	}

	head, r.URL.Path = splitPath(r.URL.Path)
	num, _, err := parseInt(r.URL.Path)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, ErrWrongPath)
		return
	}

	respondFunc := func(fn a.ListFunc) {
		list, err := fn(r.Context(), year, num)
		if err != nil {
			clientError(w, r, http.StatusBadRequest, ErrWrongPath)
			return
		}
		respond(w, http.StatusOK, list)
	}

	switch head {
	case "day":
		respondFunc(h.app.EventListForDay)
	case "week":
		respondFunc(h.app.EventListForWeek)
	case "month":
		respondFunc(h.app.EventListForMonth)
	default:
		http.NotFound(w, r)
	}
}
