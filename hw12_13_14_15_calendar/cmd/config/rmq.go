package config

type Rmq struct {
	ExchangeName string `koanf:"exchange_name"`
	ExchangeType string `koanf:"exchange_type"`
	RoutingKey   string `koanf:"routing_key"`
	QueueName    string `koanf:"queue_name"`
	Login        string `koanf:"login"`
	Password     string `koanf:"password"`
}
