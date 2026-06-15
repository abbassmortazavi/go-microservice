package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Consumer{conn: conn, channel: ch}, nil
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *Consumer) Setup(exchange, exchangeType, queueName, routingKey string) (<-chan amqp.Delivery, error) {
	if err := c.channel.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := c.channel.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
		return nil, err
	}

	return c.channel.Consume(q.Name, "", false, false, false, false, nil)
}
