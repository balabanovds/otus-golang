package config

type Rmq struct {
	Host         string `koanf:"host"`
	Port         int    `koanf:"port"`
	User         string `koanf:"user"`
	Password     string `koanf:"password"`
	ExchangeName string `koanf:"exchange_name"`
	ExchangeType string `koanf:"exchange_type"`
	RoutingKey   string `koanf:"routing_key"`
	QueueName    string `koanf:"queue_name"`
	Tag          string `koanf:"consumer_tag"`
}
