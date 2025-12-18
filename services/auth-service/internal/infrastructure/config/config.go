package config

import (
	"abbassmortazavi/go-microservice/pkg/env"
)

type Config struct {
	HTTP_ADDR string `config:"HTTP_ADDR"`
}

func Load() *Config {
	cfg := &Config{
		HTTP_ADDR: env.GetString("HTTP_ADDR", ":9091"),
	}
	return cfg
}
