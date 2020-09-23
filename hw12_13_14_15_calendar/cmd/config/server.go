package config

import "time"

type Server struct {
	Host            string        `koanf:"host"`
	Port            int           `koanf:"port"`
	ShutdownTimeout time.Duration `koanf:"shutdown_timeout_ms"`
	Grpc            bool          `koanf:"grpc"`
}
