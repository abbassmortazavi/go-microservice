package config

import (
	"abbassmortazavi/go-microservice/pkg/env"
)

type Config struct {
	HTTP_ADDR     string `config:"HTTP_ADDR"`
	RABBITMQ_URL  string `config:"RABBITMQ_URL"`
	USER_EXCHANGE string `config:"USER_EXCHANGE"`
}

func Load() *Config {
	cfg := &Config{
		HTTP_ADDR:     env.GetString("HTTP_ADDR", ":9092"),
		RABBITMQ_URL:  env.GetString("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		USER_EXCHANGE: env.GetString("USER_EXCHANGE", "user.events"),
	}
	return cfg
}
