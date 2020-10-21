package grpcsrv_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/server/grpcsrv"
	memorystorage "github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var (
	start     = time.Date(2020, time.January, 1, 0, 0, 0, 1, time.Local)
	createReq = &grpcsrv.CreateEventRequest{
		Event:  grpcsrv.NewTestIncomingEvent(start.AddDate(0, 0, 256)),
		UserID: 1,
	}
)

func TestServer_CreateEvent(t *testing.T) {
	tests := []struct {
		name    string
		req     *grpcsrv.CreateEventRequest
		errCode codes.Code
	}{
		{
			name: "test nil request",
			req: &grpcsrv.CreateEventRequest{
				Event:  nil,
				UserID: 0,
			},
			errCode: codes.InvalidArgument,
		},
		{
			name: "empty incoming event inside request",
			req: &grpcsrv.CreateEventRequest{
				Event:  &grpcsrv.IncomingEvent{},
				UserID: 1,
			},
			errCode: codes.Internal,
		},
		{
			name: "already exists",
			req: &grpcsrv.CreateEventRequest{
				Event:  grpcsrv.NewTestIncomingEvent(start),
				UserID: 1,
			},
			errCode: codes.AlreadyExists,
		},
		{
			name: "ok",
			req:  createReq,
		},
	}

	ctx := context.Background()

	client, cleanUp := prepareClient(ctx, t)
	defer cleanUp()

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			resp, err := client.CreateEvent(ctx, tst.req)
			if err != nil {
				er, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tst.errCode, er.Code())
				return
			}
			require.NotNil(t, resp)
		})
	}
}

func TestServer_GetEvent(t *testing.T) {
	tests := []struct {
		name      string
		req       *grpcsrv.GetEventRequest
		createReq *grpcsrv.CreateEventRequest
		errCode   codes.Code
	}{
		{
			name:    "empty request",
			req:     &grpcsrv.GetEventRequest{},
			errCode: codes.InvalidArgument,
		},
		{
			name: "not found",
			req: &grpcsrv.GetEventRequest{
				ID:     999,
				UserID: 1,
			},
			errCode: codes.NotFound,
		},
		{
			name: "ok",
			req: &grpcsrv.GetEventRequest{
				ID:     999,
				UserID: 1,
			},
			createReq: createReq,
		},
	}

	ctx := context.Background()

	client, cleanUp := prepareClient(ctx, t)
	defer cleanUp()

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			var createdEvent *grpcsrv.Event
			if tst.createReq != nil {
				var err error
				createdEvent, err = client.CreateEvent(ctx, tst.createReq)
				require.NoError(t, err)
				tst.req.ID = createdEvent.GetID()
			}
			resp, err := client.GetEvent(ctx, tst.req)
			if err != nil {
				er, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tst.errCode, er.Code())
				return
			}
			require.Equal(t, createdEvent, resp)
		})
	}
}

func TestServer_UpdateEvent(t *testing.T) {
	tests := []struct {
		name      string
		createReq *grpcsrv.CreateEventRequest
		req       *grpcsrv.UpdateEventRequest
		errCode   codes.Code
	}{
		{
			name:      "test nil request",
			createReq: createReq,
			req: &grpcsrv.UpdateEventRequest{
				ID:     0,
				UserID: 0,
				Event:  nil,
			},
			errCode: codes.InvalidArgument,
		},
		{
			name:      "empty incoming event inside request",
			createReq: createReq,
			req: &grpcsrv.UpdateEventRequest{
				Event:  &grpcsrv.IncomingEvent{},
				UserID: 1,
			},
			errCode: codes.Internal,
		},
		{
			name: "not found",
			req: &grpcsrv.UpdateEventRequest{
				Event:  grpcsrv.NewTestIncomingEvent(start.AddDate(0, 0, 255)),
				ID:     1,
				UserID: 1,
			},
			errCode: codes.NotFound,
		},
		{
			name:      "ok",
			createReq: createReq,
			req: &grpcsrv.UpdateEventRequest{
				Event:  grpcsrv.NewTestIncomingEvent(start.AddDate(0, 0, 255)),
				UserID: 1,
			},
		},
	}

	ctx := context.Background()

	client, cleanUp := prepareClient(ctx, t)
	defer cleanUp()

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			if tst.createReq != nil {
				var err error
				createdEvent, err := client.CreateEvent(ctx, tst.createReq)
				require.NoError(t, err)
				tst.req.ID = createdEvent.GetID()
			}
			resp, err := client.UpdateEvent(ctx, tst.req)
			if err != nil {
				er, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tst.errCode, er.Code())
				return
			}
			require.Equal(t, tst.req.Event.GetStartAt().Nanos, resp.StartAt.Nanos)
		})
	}
}

func TestServer_DeleteEvent(t *testing.T) {
	tests := []struct {
		name      string
		createReq *grpcsrv.CreateEventRequest
		req       *grpcsrv.DeleteEventRequest
		errCode   codes.Code
	}{
		{
			name: "zero value inside request",
			req: &grpcsrv.DeleteEventRequest{
				UserID: 1,
			},
			errCode: codes.InvalidArgument,
		},
		{
			name:      "ok",
			createReq: createReq,
			req: &grpcsrv.DeleteEventRequest{
				UserID: 1,
			},
		},
	}

	ctx := context.Background()

	client, cleanUp := prepareClient(ctx, t)
	defer cleanUp()

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			var createdEvent *grpcsrv.Event
			if tst.createReq != nil {
				var err error
				createdEvent, err = client.CreateEvent(ctx, tst.createReq)
				require.NoError(t, err)
				tst.req.ID = createdEvent.GetID()
			}
			_, err := client.DeleteEvent(ctx, tst.req)
			if err != nil {
				er, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tst.errCode, er.Code())
				return
			}
			getReq := &grpcsrv.GetEventRequest{
				ID:     createdEvent.ID,
				UserID: createdEvent.UserID,
			}
			_, err = client.GetEvent(ctx, getReq)
			require.Error(t, err)
			er, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, codes.NotFound, er.Code())
		})
	}
}

func TestServer_EventList(t *testing.T) {
	ctx := context.Background()

	client, cleanUp := prepareClient(ctx, t)
	defer cleanUp()

	tests := []struct {
		name     string
		listFunc func(ctx context.Context, in *grpcsrv.EventListRequest, opts ...grpc.CallOption) (*grpcsrv.EventList, error)
		value    uint32
		length   int
	}{
		{
			name:     "day",
			listFunc: client.EventListForDay,
			value:    2,
			length:   1,
		},
		{
			name:     "day empty list",
			listFunc: client.EventListForDay,
			value:    200,
			length:   0,
		},
		{
			name:     "week",
			listFunc: client.EventListForWeek,
			value:    2,
			length:   7,
		},
		{
			name:     "week empty list",
			listFunc: client.EventListForWeek,
			value:    20,
			length:   0,
		},
		{
			name:     "month",
			listFunc: client.EventListForMonth,
			value:    1,
			length:   20,
		},
		{
			name:     "month empty list",
			listFunc: client.EventListForMonth,
			value:    10,
			length:   0,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			req := &grpcsrv.EventListRequest{
				UserID: 1,
				Year:   2020,
				Value:  tst.value,
			}
			res, err := tst.listFunc(ctx, req)
			require.NoError(t, err)
			require.Len(t, res.List, tst.length)
		})
	}
}

func prepareClient(ctx context.Context, t *testing.T) (client grpcsrv.EventsServiceClient, cleanUp func()) {
	conn, err := grpc.DialContext(
		ctx, "",
		grpc.WithInsecure(),
		grpc.WithContextDialer(dialer(t)),
	)
	require.NoError(t, err)

	return grpcsrv.NewEventsServiceClient(conn), func() {
		err := conn.Close()
		require.NoError(t, err)
	}
}

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	lsn := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	a := app.New(memorystorage.NewTestStorage(start, 20))

	grpcsrv.RegisterEventsServiceServer(
		srv,
		grpcsrv.NewEventsServiceServer(a, config.GRPC{}),
	)

	go func() {
		err := srv.Serve(lsn)
		require.NoError(t, err)
	}()

	return func(context.Context, string) (net.Conn, error) {
		return lsn.Dial()
	}
}
