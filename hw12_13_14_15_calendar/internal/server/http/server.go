package internalhttp

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
)

type Server struct {
	config server.Config
}

type Application interface {
	// TODO
}

func NewServer(app Application, config server.Config) *Server {
	return &Server{}
}

func (s *Server) Start() error {
	// TODO
	return nil
}

func (s *Server) Stop() error {
	//ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	// TODO
	return nil
}

// TODO
