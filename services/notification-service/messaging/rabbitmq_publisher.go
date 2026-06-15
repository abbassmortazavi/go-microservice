package messaging

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch       *amqp091.Channel
	exchange string
}

func NewPublisher(ch *amqp091.Channel, exchange string) *Publisher {
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}
}

func (p *Publisher) Publish(ctx context.Context, routekey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(
		ctx,
		p.exchange,
		routekey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
