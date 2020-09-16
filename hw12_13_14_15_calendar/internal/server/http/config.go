package internalhttp

import "time"

type Config struct {
	Host            string        `koanf:"host"`
	Port            int           `koanf:"port"`
	ShutdownTimeout time.Duration `koanf:"shutdown_timeout_ms"`
}