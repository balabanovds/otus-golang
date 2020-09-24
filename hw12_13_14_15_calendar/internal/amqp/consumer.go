package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/cmd/config"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/internal/models"
	"github.com/balabanovds/otus-golang/hw12_13_14_15_calendar/pkg/utils"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type EventConsumer struct {
	cfg  config.Rmq
	conn *amqp.Connection
	done chan error
}

func NewConsumer(cfg config.Rmq) (Consumer, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("consume: dial: %w", err)
	}
	zap.L().Info("consumer: connected to amqp")
	return &EventConsumer{
		conn: conn,
		cfg:  cfg,
		done: make(chan error),
	}, nil
}

func (c *EventConsumer) Consume(ctx context.Context) (<-chan models.MQNotification, error) {
	msgCh := make(chan models.MQNotification)

	channel, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("consumer: create channel: %w", err)
	}
	defer utils.Close(channel)

	go func() {
		<-ctx.Done()
		utils.Close(channel)
	}()

	if err := exchangeDeclare(channel, c.cfg.ExchangeName, c.cfg.ExchangeType); err != nil {
		return nil, fmt.Errorf("consumer: exchange declare: %w", err)
	}

	deliveries, err := c.consumeQueue(channel, c.cfg.QueueName)
	if err != nil {
		return nil, fmt.Errorf("consumer: queue consume: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-deliveries:
				var msg models.MQNotification
				if err := json.Unmarshal(d.Body, &msg); err != nil {
					zap.L().Error("consumer: unmarshal message", zap.Error(err))
					continue
				}
				select {
				case <-ctx.Done():
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

	if err = channel.QueueBind(
		queue.Name,         // name of the queue
		c.cfg.RoutingKey,   // bindingKey
		c.cfg.ExchangeName, // sourceExchange
		false,              // noWait
		nil,                // arguments
	); err != nil {
		return nil, fmt.Errorf("consumer: queue bind: %w", err)
	}

	deliveries, err := channel.Consume(
		queue.Name, // name
		c.cfg.Tag,  // consumerTag,
		true,       // autoAck
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
