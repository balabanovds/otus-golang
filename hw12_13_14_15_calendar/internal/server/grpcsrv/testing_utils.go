package grpcsrv

import (
	"time"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"

	"github.com/golang/protobuf/ptypes"
)

func NewTestIncomingEvent(tm time.Time) *IncomingEvent {
	t, _ := ptypes.TimestampProto(tm)
	return &IncomingEvent{
		Title:          utils.RandString(10),
		StartAt:        t,
		Duration:       ptypes.DurationProto(1 * time.Hour),
		Description:    utils.RandString(20),
		RemindDuration: ptypes.DurationProto(0 * time.Hour),
	}
}
