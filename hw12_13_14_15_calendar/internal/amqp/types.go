package amqp

import (
	"fmt"
	"io"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
)

type Publisher interface {
	Publish(body []byte) error
	io.Closer
	fmt.Stringer
}

type Consumer interface {
	Consume() (<-chan models.MQNotification, error)
	io.Closer
	fmt.Stringer
}
