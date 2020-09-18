package router

import (
	"net/http"

	a "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
)

type rootHandler struct {
	eventsHandler eventsHandler
	eventHandler  eventHandler
}

func newRootHandler(app a.Application) rootHandler {
	return rootHandler{
		eventsHandler: newEventsHandler(app),
		eventHandler:  newEventHandler(app),
	}
}

func (h rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = splitPath(r.URL.Path)

	switch head {
	case "events":
		h.eventsHandler.ServeHTTP(w, r)
	case "event":
		h.eventHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
