package mock

import (
	"encoding/json"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
)

type FakeQueue struct {
	channel chan models.MQNotification
}

func NewFakeQueue(l int) *FakeQueue {
	return &FakeQueue{
		channel: make(chan models.MQNotification, l),
	}
}

func (f *FakeQueue) Publish(data []byte) error {
	var msg models.MQNotification
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	f.channel <- msg
	return nil
}

func (f *FakeQueue) Consume() (<-chan models.MQNotification, error) {
	return f.channel, nil
}

func (f *FakeQueue) Close() error {
	close(f.channel)
	return nil
}

func (f *FakeQueue) String() string {
	return "fake publisher"
}
