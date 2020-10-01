package mock

import (
	"context"
	"encoding/json"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/mq"
)

type FakeQueue struct {
	channel chan models.MQNotification
}

func NewFakeQueue(l int) *FakeQueue {
	return &FakeQueue{
		channel: make(chan models.MQNotification, l),
	}
}

func (f *FakeQueue) Channel() mq.Channel {
	return nil
}

func (f *FakeQueue) Publish(_ context.Context, data []byte) error {
	var msg models.MQNotification
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	f.channel <- msg
	return nil
}

func (f *FakeQueue) Consume(_ context.Context) (<-chan models.MQNotification, <-chan error) {
	return f.channel, nil
}

func (f *FakeQueue) Close() error {
	close(f.channel)
	return nil
}

func (f *FakeQueue) String() string {
	return "fake publisher"
}
