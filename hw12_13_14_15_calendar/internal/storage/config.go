package storage

type Config struct {
	Sql bool   `koanf:"sql"`
	Dsn string `koanf:"dsn"`
}
