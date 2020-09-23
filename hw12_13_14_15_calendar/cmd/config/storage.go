package config

type Storage struct {
	SQL bool   `koanf:"sql"`
	Dsn string `koanf:"dsn"`
}
