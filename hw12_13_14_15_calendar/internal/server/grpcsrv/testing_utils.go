package grpcsrv

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes"
)

func NewTestIncomingEvent(tm time.Time) *IncomingEvent {
	t, _ := ptypes.TimestampProto(tm)
	return &IncomingEvent{
		Title:          storage.RandString(10),
		StartTime:      t,
		Duration:       ptypes.DurationProto(1 * time.Hour),
		Description:    storage.RandString(20),
		RemindDuration: ptypes.DurationProto(0 * time.Hour),
	}
}
