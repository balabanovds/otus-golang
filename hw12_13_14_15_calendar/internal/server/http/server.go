package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	config server.Config
	srv    *http.Server
	app    Application
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.IncomingEvent) error
	Get(ctx context.Context, id int) (storage.Event, error)
	Update(ctx context.Context, id int, event storage.IncomingEvent) error
	Delete(ctx context.Context, id int) error
	ListForDay(ctx context.Context, date time.Time) []storage.Event
	ListForWeek(ctx context.Context, date time.Time) []storage.Event
	ListForMonth(ctx context.Context, date time.Time) []storage.Event
}

func NewServer(app Application, config server.Config) *Server {
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
	return &Server{config, srv, app}
}

func (s *Server) Start() error {
	s.srv.Handler = s.routes()
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
