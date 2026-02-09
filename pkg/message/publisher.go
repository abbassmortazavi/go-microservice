package message

import (
	"abbassmortazavi/go-microservice/pkg/env"
	"abbassmortazavi/go-microservice/services/auth-service/config"
	messaging2 "abbassmortazavi/go-microservice/services/auth-service/messaging"
	"log"
)

var Publisher *messaging2.Publisher

func Init() {
	rabbitmqURL := env.GetString("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	cfg := config.Load()
	conn, ch := messaging2.NewRabbitMQ(rabbitmqURL)
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare(
		cfg.UserExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	Publisher = messaging2.NewPublisher(ch, cfg.UserExchange)
}
