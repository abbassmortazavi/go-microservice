package config

import "abbassmortazavi/go-microservice/pkg/env"

type Config struct {
	RabbitMQURL  string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() *Config {
	return &Config{
		RabbitMQURL:  env.GetString("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		SMTPHost:     env.GetString("SMTP_HOST", "sandbox.smtp.mailtrap.io"),
		SMTPPort:     env.GetString("SMTP_PORT", "587"),
		SMTPUsername: env.GetString("SMTP_USERNAME", "ceddf8b2b5de1d"),
		SMTPPassword: env.GetString("SMTP_PASSWORD", "9b1825ca7c514f"),
		SMTPFrom:     env.GetString("SMTP_FROM", "abbassmortazavi@gmail.com"),
	}

}
