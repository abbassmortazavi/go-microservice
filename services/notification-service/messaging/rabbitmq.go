package messaging

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(url string) (*amqp091.Connection, *amqp091.Channel) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return conn, ch
}
