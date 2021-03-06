package router

import (
	"encoding/json"
	"errors"
	"net/http"

	a "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type eventHandler struct {
	app a.Application
}

func newEventHandler(app a.Application) eventHandler {
	return eventHandler{app}
}

/*
	GET 	/:id	// get event
	POST 	/ 		// create event
	PUT  	/:id	// update event
	DELETE 	/:id 	// delete event
*/

func (h *eventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetEvent(w, r)
	case http.MethodPost:
		h.handleCreateEvent(w, r)
	case http.MethodPut:
		h.handleUpdateEvent(w, r)
	case http.MethodDelete:
		h.handleDeleteEvent(w, r)
	default:
		clientError(w, r, http.StatusMethodNotAllowed, ErrWrongMethod)
	}
}

func (h *eventHandler) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	id, _, err := parseInt(r.URL.Path)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, err)
		return
	}

	event, err := h.app.GetEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrEvent404) {
			clientError(w, r, http.StatusNotFound, err)
			return
		}
		serverError(w, err)
		return
	}

	respond(w, http.StatusOK, event)
}

func (h *eventHandler) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			serverError(w, err)
		}
	}()
	var incomingEvent server.IncomingEvent
	err := json.NewDecoder(r.Body).Decode(&incomingEvent)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, err)
		return
	}

	respEvent, err := h.app.CreateEvent(r.Context(), incomingEvent)
	if err != nil {
		if errors.Is(err, storage.ErrEventExists) {
			clientError(w, r, http.StatusBadRequest, err)
			return
		}
		serverError(w, err)
		return
	}

	respond(w, http.StatusCreated, respEvent)
}

func (h *eventHandler) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			serverError(w, err)
		}
	}()
	id, _, err := parseInt(r.URL.Path)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, err)
		return
	}

	var incomingEvent server.IncomingEvent
	err = json.NewDecoder(r.Body).Decode(&incomingEvent)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, err)
		return
	}

	respEvent, err := h.app.UpdateEvent(r.Context(), id, incomingEvent)
	if err != nil {
		if errors.Is(err, storage.ErrEvent404) {
			clientError(w, r, http.StatusNotFound, err)
			return
		}
		serverError(w, err)
		return
	}

	respond(w, http.StatusOK, respEvent)
}

func (h *eventHandler) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, _, err := parseInt(r.URL.Path)
	if err != nil {
		clientError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.app.DeleteEvent(r.Context(), id); err != nil {
		if errors.Is(err, storage.ErrEvent404) {
			clientError(w, r, http.StatusNotFound, err)
			return
		}
		serverError(w, err)
		return
	}
}
