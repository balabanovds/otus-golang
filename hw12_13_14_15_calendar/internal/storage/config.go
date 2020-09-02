package storage

type Config struct {
	SQL bool   `koanf:"sql"`
	Dsn string `koanf:"dsn"`
}
