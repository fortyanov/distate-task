package config

type Database struct {
	DSN string `envconfig:"dsn"`
}
