package grpcsrv

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	a "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc -I=./schema --go_out=plugins=grpc:. ./schema/event_service.proto

type Server struct {
	app a.Application
	cfg config.Server
	srv *grpc.Server
}

func NewServer(app a.Application, cfg config.Server) *Server {
	return &Server{
		app: app,
		cfg: cfg,
		srv: grpc.NewServer(grpc.UnaryInterceptor(logInterceptor)),
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		zap.L().Error("start grpc server", zap.Error(err))
		return err
	}
	RegisterEventsServiceServer(s.srv, s)
	zap.L().Info("start grpc server",
		zap.String("host", s.cfg.Host),
		zap.Int("port", s.cfg.Port))

	return s.srv.Serve(lsn)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}

func (s *Server) CreateEvent(ctx context.Context, req *CreateEventRequest) (*Event, error) {
	if req == nil || req.GetEvent() == nil || req.GetUserID() == 0 {
		zap.L().Warn("invalid argument")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	ctx = context.WithValue(ctx, a.CtxKeyUserID, int(req.UserID))

	incomingEvent, err := copyProtoToIncoming(req.Event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event, err := s.app.CreateEvent(ctx, *incomingEvent)
	if err != nil {
		if errors.Is(err, storage.ErrEventExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return copyEventToProto(event)
}

func (s *Server) GetEvent(ctx context.Context, req *GetEventRequest) (*Event, error) {
	if req == nil || req.GetID() == 0 || req.GetUserID() == 0 {
		zap.L().Warn("invalid argument")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	ctx = context.WithValue(ctx, a.CtxKeyUserID, int(req.UserID))
	event, err := s.app.GetEvent(ctx, int(req.ID))
	if err != nil {
		if errors.Is(err, storage.ErrEvent404) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return copyEventToProto(event)
}

func (s *Server) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*Event, error) {
	if req == nil || req.GetID() == 0 || req.GetUserID() == 0 || req.Event == nil {
		zap.L().Warn("invalid argument")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	ctx = context.WithValue(ctx, a.CtxKeyUserID, int(req.UserID))

	incomingEvent, err := copyProtoToIncoming(req.Event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event, err := s.app.UpdateEvent(ctx, int(req.ID), *incomingEvent)
	if err != nil {
		if errors.Is(err, storage.ErrEvent404) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return copyEventToProto(event)
}

func (s *Server) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*empty.Empty, error) {
	if req == nil || req.GetID() == 0 || req.GetUserID() == 0 {
		zap.L().Warn("invalid argument")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	ctx = context.WithValue(ctx, a.CtxKeyUserID, int(req.UserID))

	return &empty.Empty{}, s.app.DeleteEvent(ctx, int(req.ID))
}

func (s *Server) EventListForDay(ctx context.Context, req *EventListRequest) (*EventList, error) {
	return s.runEventListFunc(ctx, req, s.app.EventListForDay)
}

func (s *Server) EventListForWeek(ctx context.Context, req *EventListRequest) (*EventList, error) {
	return s.runEventListFunc(ctx, req, s.app.EventListForWeek)
}

func (s *Server) EventListForMonth(ctx context.Context, req *EventListRequest) (*EventList, error) {
	return s.runEventListFunc(ctx, req, s.app.EventListForMonth)
}

func (s *Server) runEventListFunc(ctx context.Context, req *EventListRequest, fn a.ListFunc) (*EventList, error) {
	if req == nil || req.GetUserID() == 0 || req.GetValue() == 0 || req.GetYear() == 0 {
		zap.L().Warn("invalid argument")
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	ctx = context.WithValue(ctx, a.CtxKeyUserID, int(req.UserID))

	return copyEventList(fn(ctx, int(req.GetYear()), int(req.GetValue())))
}

func copyEventToProto(event models.Event) (*Event, error) {
	startTime, err := ptypes.TimestampProto(event.StartTime)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:             int64(event.ID),
		Title:          event.Title,
		StartTime:      startTime,
		Duration:       ptypes.DurationProto(event.Duration),
		Description:    event.Description,
		UserID:         int64(event.UserID),
		RemindDuration: ptypes.DurationProto(event.RemindDuration),
	}, nil
}

func copyProtoToIncoming(from *IncomingEvent) (*models.IncomingEvent, error) {
	startTime, err := ptypes.Timestamp(from.StartTime)
	if err != nil {
		return nil, err
	}
	duration, err := ptypes.Duration(from.Duration)
	if err != nil {
		return nil, err
	}
	remDuration, err := ptypes.Duration(from.RemindDuration)
	if err != nil {
		return nil, err
	}
	return &models.IncomingEvent{
		Title:          from.Title,
		StartTime:      startTime,
		Duration:       duration,
		Description:    from.Description,
		RemindDuration: remDuration,
	}, nil
}

func copyEventList(list models.EventsList, err error) (*EventList, error) {
	if err != nil {
		return nil, err
	}

	protoList := make([]*Event, list.Len)
	for i := range list.List {
		pe, err := copyEventToProto(list.List[i])
		if err != nil {
			return nil, err
		}
		protoList[i] = pe
	}

	return &EventList{
		List: protoList,
		Time: list.Time,
		Len:  int32(list.Len),
	}, nil
}
