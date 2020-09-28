package mq

import (
	"context"
	"fmt"
	"io"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	a "github.com/streadway/amqp"
)

type Message struct {
	Data models.MQNotification
	Err  error
}

type Publisher interface {
	Channel() Channel
	Publish(ctx context.Context, body []byte) error
	io.Closer
	fmt.Stringer
}

type Consumer interface {
	Channel() Channel
	Consume(ctx context.Context) (<-chan Message, error)
	io.Closer
	fmt.Stringer
}

type Channel interface {
	Open() error
	Get() *a.Channel
	Close() error
}
