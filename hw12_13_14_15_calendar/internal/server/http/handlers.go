package internalhttp

import (
	"fmt"
	"net/http"
)

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
		w.WriteHeader(200)
	}
}
