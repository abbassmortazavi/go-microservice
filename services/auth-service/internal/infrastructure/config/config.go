package config

import "os"

type Config struct {
	DB_URL     string `envconfig:"DB_URL"`
	JWT_SECRET string `envconfig:"JWT_SECRET"`
}

// "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
func Load() *Config {
	cfg := &Config{
		DB_URL:     os.Getenv("DB_URL"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}
	return cfg
}
