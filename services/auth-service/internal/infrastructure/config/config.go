package config

import (
	"abbassmortazavi/go-microservice/pkg/env"
)

type Config struct {
	HttpAddr     string `config:"HTTP_ADDR"`
	RabbitmqUrl  string `config:"RABBITMQ_URL"`
	UserExchange string `config:"USER_EXCHANGE"`
}

func Load() *Config {
	cfg := &Config{
		HttpAddr:     env.GetString("HTTP_ADDR", ":9092"),
		RabbitmqUrl:  env.GetString("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		UserExchange: env.GetString("USER_EXCHANGE", "user.events"),
	}
	return cfg
}
