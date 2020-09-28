package rabbitmq

import "github.com/streadway/amqp"

func exchangeDeclare(ch *amqp.Channel, exchName, exchType string) error {
	return ch.ExchangeDeclare(
		exchName,
		exchType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	)
}
