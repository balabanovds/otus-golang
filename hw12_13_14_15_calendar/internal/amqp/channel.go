package amqp

import (
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	a "github.com/streadway/amqp"
	"go.uber.org/zap"
)

type channel struct {
	cfg  config.Rmq
	conn *a.Connection
	ch   *a.Channel
}

func newChannel(cfg config.Rmq, conn *a.Connection) Channel {
	return &channel{
		cfg:  cfg,
		conn: conn,
	}
}

func (c *channel) Open() error {
	var err error
	c.ch, err = c.conn.Channel()
	if err != nil {
		return err
	}

	if err := exchangeDeclare(c.ch, c.cfg.ExchangeName, c.cfg.ExchangeType); err != nil {
		return err
	}
	zap.L().Info("exchange declared",
		zap.String("name", c.cfg.ExchangeName),
		zap.String("type", c.cfg.ExchangeType),
	)
	return nil
}

func (c *channel) Get() *a.Channel {
	return c.ch
}

func (c *channel) Close() error {
	return c.ch.Close()
}
