package config

import (
	"abbassmortazavi/go-microservice/pkg/env"
	"time"
)

type GlobalConfig struct {
	JWT_SECRET     string        `config:"JWT_SECRET"`
	Host           string        `config:"DB_HOST"`
	Port           string        `config:"DB_PORT"`
	Debug          string        `config:"APP_DEBUG"`
	Username       string        `config:"DB_USERNAME"`
	Password       string        `config:"DB_PASSWORD"`
	Name           string        `config:"DB_NAME"`
	MaxIdle        int           `config:"DB_MAX_IDLE"`
	MaxConn        int           `config:"DB_MAX_CONN"`
	MaxIdleTimeout time.Duration `config:"DB_MAX_IDLE_TIMEOUT"`
}

func Load() *GlobalConfig {
	cfg := &GlobalConfig{
		JWT_SECRET:     env.GetString("JWT_SECRET", ""),
		Host:           env.GetString("DB_HOST", "postgres-service"),
		Port:           env.GetString("DB_PORT", "5432"),
		Debug:          env.GetString("APP_DEBUG", "false"),
		Username:       env.GetString("DB_USERNAME", "root"),
		Password:       env.GetString("DB_PASSWORD", "root"),
		Name:           env.GetString("DB_NAME", "microservice_db"),
		MaxIdle:        10,
		MaxConn:        10,
		MaxIdleTimeout: 10 * time.Second,
	}
	return cfg
}
