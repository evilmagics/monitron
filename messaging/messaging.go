package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"monitron-server/config"
	"monitron-server/internal/reportgen"
	"monitron-server/models"
	"monitron-server/utils"
)

// Channel is the RabbitMQ channel
var Channel *amqp.Channel

// Conn is the RabbitMQ connection
var Conn *amqp.Connection

// InitRabbitMQ initializes the RabbitMQ connection and channel
func InitRabbitMQ(cfg *config.Config) {
	var err error
	Conn, err = amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	log.Println("Successfully connected to RabbitMQ")
}

// CloseRabbitMQ closes the RabbitMQ channel and connection
func CloseRabbitMQ() {
	if Channel != nil {
		Channel.Close()
	}
	if Conn != nil {
		Conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}

// PublishMessage publishes a message to the specified queue
func PublishMessage(queueName string, body []byte) error {
	if Channel == nil {
		return fmt.Errorf("RabbitMQ channel is not initialized")
	}

	q, err := Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = Channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("Failed to publish a message: %w", err)
	}
	log.Printf(" [x] Sent %s to %s", body, queueName)
	return nil
}

// ConsumeMessages consumes messages from the specified queue
func ConsumeMessages(queueName string, handler func([]byte)) {
	if Channel == nil {
		log.Fatalf("RabbitMQ channel is not initialized")
	}

	q, err := Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (we'll manually ack)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message from %s: %s", queueName, d.Body)
			handler(d.Body)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages in %s. To exit press CTRL+C", queueName)
	<-forever
}

// SetupConsumers sets up all necessary message consumers
func SetupConsumers() {
	// Example: Report generation consumer
	go ConsumeMessages("report_generation_queue", func(body []byte) {
		log.Printf("Processing report generation task: %s", string(body))
		var report models.Report
		if err := json.Unmarshal(body, &report); err != nil {
			log.Printf("Error unmarshalling report details: %v", err)
			return
		}

		// Simulate report generation and save to file
		filePath, err := reportgen.GenerateReportFile(report)
		if err != nil {
			log.Printf("Error generating report file: %v", err)
			return
		}

		// Update report entry in DB with file path
		// This requires a DB connection, which is not directly available in consumer goroutine.
		// In a real app, you'd pass the DB connection or use a service layer.
		// For now, we'll just log it.
		log.Printf("Report %s generated and saved to: %s", report.ID, filePath)

		// TODO: Update report entry in database with file_path
		// Example: db.Exec(`UPDATE reports SET file_path = $1 WHERE id = $2`, filePath, report.ID)

		log.Printf("Report generation task completed for: %s", string(body))
	})

	// Email sending consumer
	go ConsumeMessages("email_sending_queue", func(body []byte) {
		log.Printf("Processing email sending task: %s", string(body))
		var emailDetails struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}
		if err := json.Unmarshal(body, &emailDetails); err != nil {
			log.Printf("Error unmarshalling email details: %v", err)
			return
		}

		cfg := config.LoadConfig()
		err := utils.SendEmail(cfg, emailDetails.To, emailDetails.Subject, emailDetails.Body)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			// TODO: Implement retry logic or move to DLQ
			return
		}
		log.Printf("Email sent successfully to: %s", emailDetails.To)
	})

	// TODO: Add more consumers for other background tasks (e.g., health checks, notifications)
}


