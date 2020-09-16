package router

import (
	"net/http"

	a "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
)

type rootHandler struct {
	eventsHandler eventsHandler
}

func newRootHandler(app a.Application) rootHandler {
	return rootHandler{eventsHandler: newEventsHandler(app)}
}

func (h rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = splitPath(r.URL.Path)

	switch head {
	case "events":
		h.eventsHandler.ServeHTTP(w, r)
	case "event":
	default:
		http.NotFound(w, r)
	}
}
