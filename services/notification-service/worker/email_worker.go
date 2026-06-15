package worker

import (
	"abbassmortazavi/go-microservice/pkg/events"
	"abbassmortazavi/go-microservice/services/notification-service/config"
	"abbassmortazavi/go-microservice/services/notification-service/messaging"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
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
	consumer, err := messaging.NewConsumer(w.cfg.RabbitMQURL)
	log.Println("Consumer created")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer consumer.Close()

	msgs, err := consumer.Setup(
		w.cfg.RabbitmqExchange, // باید با exchange در auth-service یکسان باشد
		"topic",                // چون Publisher با "direct" declare می‌کند
		"email_queue",          // اسم queue مخصوص ایمیل
		"user.registered",      // routing key همان چیزی که auth-service پابلیش می‌کند
	)
	if err != nil {
		log.Fatalf("Failed to setup consumer: %s", err)
	}

	log.Printf("Email worker started, waiting for messages...")

	for msg := range msgs {
		var event events.UserRegistered
		//var event eventpb.UserRegistered
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to unmarshal event: %s", err)
			msg.Nack(false, false) // پیام خراب است، دوباره صف نشود
			continue
		}

		log.Printf("Received UserRegistered event for: %s", event.Email)

		if err := w.sendEmail(event.Email, event.Name); err != nil {
			log.Printf("Failed to send email to %s: %s", event.Email, err)
			msg.Nack(false, true) // دوباره صف شود تا retry شود
			continue
		}

		log.Printf("Welcome email sent to %s", event.Email)
		msg.Ack(false)
	}
}

func (w *EmailWorker) sendEmail(to, name string) error {
	auth := smtp.PlainAuth("", w.cfg.SMTPUsername, w.cfg.SMTPPassword, w.cfg.SMTPHost)

	subject := "Welcome to our platform!"
	body := fmt.Sprintf("Hi %s,\r\n\r\nThank you for registering with us!", name)

	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)

	return smtp.SendMail(
		w.cfg.SMTPHost+":"+w.cfg.SMTPPort,
		auth,
		w.cfg.SMTPFrom,
		[]string{to},
		msg,
	)
}
