package router

import (
	"net/http"

	a "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/justinas/alice"
)

type Router struct {
	rootHandler rootHandler
}

func New(app a.Application) *Router {
	return &Router{rootHandler: newRootHandler(app)}
}

func (rtr *Router) RootHandler() http.Handler {
	return alice.New(recoverPanic, jsonContent, logRequest, authorize).Then(rtr.rootHandler)
}
