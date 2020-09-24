package amqp

import (
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type EventPublisher struct {
	cfg  config.Rmq
	conn *amqp.Connection
}

func NewPublisher(cfg config.Rmq) (Publisher, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	zap.L().Info("connected to amqp")
	return &EventPublisher{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (p *EventPublisher) Publish(body []byte) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer utils.Close(channel)

	if err := exchangeDeclare(channel, p.cfg.ExchangeName, p.cfg.ExchangeType); err != nil {
		return err
	}

	if err := channel.Publish(
		p.cfg.ExchangeName,
		p.cfg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			Body:            body,
			DeliveryMode:    amqp.Persistent,
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
