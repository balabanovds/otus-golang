package config

type Logger struct {
	Level   string `koanf:"level"`
	LogFile string `koanf:"log_file"`
}
