package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ")
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open a channel")
	}
	log.Info().Msg("Successfully connected to RabbitMQ")
}

// CloseRabbitMQ closes the RabbitMQ channel and connection
func CloseRabbitMQ() {
	if Channel != nil {
		Channel.Close()
	}
	if Conn != nil {
		Conn.Close()
	}
	log.Info().Msg("RabbitMQ connection closed")
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
	log.Info().Msgf(" [x] Sent %s to %s", body, queueName)
	return nil
}

// ConsumeMessages consumes messages from the specified queue
func ConsumeMessages(queueName string, handler func([]byte)) {
	if Channel == nil {
		log.Fatal().Err(fmt.Errorf("RabbitMQ channel is not initialized")).Msg("RabbitMQ channel is not initialized")
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
		log.Fatal().Err(err).Msg("Failed to declare a queue")
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
		log.Fatal().Err(err).Msg("Failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Info().Msgf("Received a message from %s: %s", queueName, d.Body)
			handler(d.Body)
			d.Ack(false)
		}
	}()

	log.Info().Msgf(" [*] Waiting for messages in %s. To exit press CTRL+C", queueName)
	<-forever
}

// SetupConsumers sets up all necessary message consumers
func SetupConsumers() {
	// Example: Report generation consumer
	go ConsumeMessages("report_generation_queue", func(body []byte) {
		log.Info().Msgf("Processing report generation task: %s", string(body))
		var report models.Report
		if err := json.Unmarshal(body, &report); err != nil {
			log.Error().Err(err).Msg("Error unmarshalling report details")
			return
		}

		// Simulate report generation and save to file
		filePath, err := reportgen.GenerateReportFile(report)
		if err != nil {
			log.Error().Err(err).Msg("Error generating report file")
			return
		}

		// Update report entry in DB with file path
		// This requires a DB connection, which is not directly available in consumer goroutine.
		// In a real app, you'd pass the DB connection or use a service layer.
		// For now, we'll just log it.
		log.Info().Msgf("Report %s generated and saved to: %s", report.ID, filePath)

		// TODO: Update report entry in database with file_path
		// Example: db.Exec(`UPDATE reports SET file_path = $1 WHERE id = $2`, filePath, report.ID)

		log.Info().Msgf("Report generation task completed for: %s", string(body))
	})

	// Email sending consumer
	go ConsumeMessages("email_sending_queue", func(body []byte) {
		log.Info().Msgf("Processing email sending task: %s", string(body))
		var emailDetails struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}
		if err := json.Unmarshal(body, &emailDetails); err != nil {
			log.Error().Err(err).Msg("Error unmarshalling email details")
			return
		}

		cfg := config.LoadConfig()
		err := utils.SendEmail(cfg, emailDetails.To, emailDetails.Subject, emailDetails.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error sending email")
			// TODO: Implement retry logic or move to DLQ
			return
		}
		log.Info().Msgf("Email sent successfully to: %s", emailDetails.To)
	})

	// TODO: Add more consumers for other background tasks (e.g., health checks, notifications)
}


