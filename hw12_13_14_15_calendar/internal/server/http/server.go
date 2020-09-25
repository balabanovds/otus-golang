package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/http/router"
	"go.uber.org/zap"
)

type Server struct {
	config config.HTTP
	srv    *http.Server
	a      app.Application
}

func New(app app.Application, config config.HTTP) server.IServer {
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

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	zap.L().Info("close http server")

	return s.srv.Shutdown(ctx)
}

func (s *Server) String() string {
	return "http server"
}
