package amqp

import (
	"context"
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	a "github.com/streadway/amqp"
	"go.uber.org/zap"
)

type EventPublisher struct {
	cfg     config.Rmq
	conn    *a.Connection
	channel Channel
}

func NewPublisher(cfg config.Rmq) (Publisher, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := a.Dial(uri)
	if err != nil {
		return nil, err
	}
	zap.L().Info("publisher: connected to amqp")
	return &EventPublisher{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (p *EventPublisher) Channel() Channel {
	if p.conn == nil {
		panic("connection is nil")
	}
	if p.channel == nil {
		p.channel = newChannel(p.cfg, p.conn)
	}
	return p.channel
}

func (p *EventPublisher) Publish(ctx context.Context, body []byte) error {
	if p.channel == nil {
		return ErrChannelNil
	}

	if err := p.channel.Get().Publish(
		p.cfg.ExchangeName,
		p.cfg.RoutingKey,
		false,
		false,
		a.Publishing{
			Headers:         a.Table{},
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			Body:            body,
			DeliveryMode:    a.Persistent,
			Priority:        0,
		},
	); err != nil {
		return err
	}

	return nil
}

func (p *EventPublisher) Close() error {
	zap.L().Info("closing publisher channel")
	return p.conn.Close()
}

func (p *EventPublisher) String() string {
	return "amqp publisher"
}
