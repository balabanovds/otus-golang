package amqp

import (
	"context"
	"fmt"
	"io"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
)

type Publisher interface {
	Publish(ctx context.Context, body []byte) error
	io.Closer
	fmt.Stringer
}

type Consumer interface {
	Consume(ctx context.Context) (<-chan models.MQNotification, error)
	io.Closer
	fmt.Stringer
}
