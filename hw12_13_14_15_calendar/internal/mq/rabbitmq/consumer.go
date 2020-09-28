package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/mq"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type EventConsumer struct {
	cfg     config.Rmq
	conn    *amqp.Connection
	channel mq.Channel
	done    chan error
}

func NewConsumer(cfg config.Rmq) (mq.Consumer, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("consume: dial: %w", err)
	}
	infoLog("connected to amqp")
	return &EventConsumer{
		conn: conn,
		cfg:  cfg,
		done: make(chan error),
	}, nil
}

func (c *EventConsumer) Channel() mq.Channel {
	if c.conn == nil {
		panic("connection is nil")
	}
	if c.channel == nil {
		c.channel = newChannel(c.cfg, c.conn)
	}
	return c.channel
}

func (c *EventConsumer) Consume(ctx context.Context) (<-chan mq.Message, error) {
	if c.channel == nil {
		return nil, mq.ErrChannelNil
	}

	go func() {
		<-ctx.Done()
		if err := c.channel.Close(); err != nil {
			zap.L().Error("close channel", zap.Error(err))
		}
	}()

	msgCh := make(chan mq.Message)

	deliveries, err := c.consumeQueue(c.channel.Get(), c.cfg.QueueName)
	if err != nil {
		return nil, fmt.Errorf("consumer: queue consume: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(msgCh)
				return
			case d := <-deliveries:
				var notif models.MQNotification
				err := json.Unmarshal(d.Body, &notif)

				msg := mq.Message{
					Data: notif,
					Err:  err,
				}

				if err == nil {
					if err := d.Ack(false); err != nil {
						zap.L().Error("delivery ack", zap.Error(err))
						continue
					}
				}

				select {
				case <-ctx.Done():
					close(msgCh)
					return
				case msgCh <- msg:
				}
			}
		}
	}()

	return msgCh, nil
}

func (c *EventConsumer) Close() error {
	return c.conn.Close()
}

func (c *EventConsumer) String() string {
	return "event consumer"
}

func (c *EventConsumer) consumeQueue(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	queue, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("consumer: declare queue: %w", err)
	}
	infoLog("queue created")

	if err := channel.Qos(1, 0, false); err != nil {
		return nil, err
	}

	if err = channel.QueueBind(
		queue.Name,         // name of the queue
		c.cfg.RoutingKey,   // bindingKey
		c.cfg.ExchangeName, // sourceExchange
		false,              // noWait
		nil,                // arguments
	); err != nil {
		return nil, fmt.Errorf("consumer: queue bind: %w", err)
	}
	infoLog("queue binded")

	deliveries, err := channel.Consume(
		queue.Name, // name
		c.cfg.Tag,  // consumerTag,
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("consumer: queue consume: %w", err)
	}

	return deliveries, nil
}

func infoLog(msg string) {
	zap.L().Info("consumer", zap.String("msg", msg))
}
