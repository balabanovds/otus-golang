package internalhttp

import (
	"context"
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http/router"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	config Config
	srv    *http.Server
	a      app.Application
}

func NewServer(app app.Application, config Config) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router.New(app).RootHandler(),
	}
	return &Server{config, srv, app}
}

func (s *Server) Start() error {
	zap.L().Info("start server",
		zap.String("host", s.config.Host),
		zap.Int("port", s.config.Port),
	)

	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
