package config

import "time"

type HTTP struct {
	Host            string        `koanf:"host"`
	Port            int           `koanf:"port"`
	ShutdownTimeout time.Duration `koanf:"shutdown_timeout_ms"`
}

type GRPC struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}
