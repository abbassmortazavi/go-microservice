package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWT_SECRET     string `config:"JWT_SECRET"`
	Host           string `config:"DB_HOST"`
	Port           string `config:"DB_PORT"`
	Debug          string `config:"APP_DEBUG"`
	Username       string `config:"DB_USERNAME"`
	Password       string `config:"DB_PASSWORD"`
	Name           string `config:"DB_NAME"`
	MaxIdle        int    `config:"DB_MAX_IDLE"`
	MaxConn        int    `config:"DB_MAX_CONN"`
	MaxIdleTimeout string `config:"DB_MAX_IDLE_TIMEOUT"`
	AppPort        string `config:"APP_PORT"`
}

// "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
func Load() *Config {
	cfg := &Config{
		JWT_SECRET:     os.Getenv("JWT_SECRET"),
		Host:           os.Getenv("DB_HOST"),
		Port:           os.Getenv("DB_PORT"),
		Debug:          os.Getenv("APP_DEBUG"),
		Username:       os.Getenv("DB_USERNAME"),
		Password:       os.Getenv("DB_PASSWORD"),
		Name:           os.Getenv("DB_NAME"),
		MaxIdle:        10,
		MaxConn:        10,
		MaxIdleTimeout: strconv.Itoa(10),
		AppPort:        os.Getenv("APP_PORT"),
	}
	return cfg
}
