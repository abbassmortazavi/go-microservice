package worker

import (
	"abbassmortazavi/go-microservice/services/notification-service/config"
	"encoding/json"
	"log"
	"net/smtp"

	"github.com/rabbitmq/amqp091-go"
)

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailWorker struct {
	cfg *config.Config
}

func NewEmailWorker(cfg *config.Config) *EmailWorker {
	return &EmailWorker{
		cfg: cfg,
	}
}

func (w *EmailWorker) Start() {
	conn, err := amqp091.Dial(w.cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare("email_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
	log.Printf("Email worker started, waiting for messages...")

	for msg := range msgs {
		var emailMessage EmailMessage
		err := json.Unmarshal(msg.Body, &emailMessage)
		if err != nil {
			log.Fatalf("Failed to unmarshal email message: %s", err)
			msg.Nack(false, false)
		} else {
			log.Printf("Received a message: %s", emailMessage.Subject)
			msg.Ack(false)
		}
	}
}
func (w *EmailWorker) sendEmail(msg EmailMessage) error {
	auth := smtp.PlainAuth("", w.cfg.SMTPUsername, w.cfg.SMTPPassword, w.cfg.SMTPHost)
	body := "Subject: " + msg.Subject + "\n" + "Body: " + msg.Body
	return smtp.SendMail(
		w.cfg.SMTPHost+":"+w.cfg.SMTPPort,
		auth,
		w.cfg.SMTPFrom,
		[]string{msg.To},
		[]byte(body),
	)
}
