package server

import "time"

type Config struct {
	Host            string        `koanf:"host"`
	Port            string        `koanf:"port"`
	ShutdownTimeout time.Duration `koanf:"shutdown_timeout_ms"`
}
