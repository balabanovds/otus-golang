package internalhttp

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHome())

	return alice.New(recoverPanic, jsonContent, logRequest).Then(mux)
}
