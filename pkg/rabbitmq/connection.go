package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection

func ConnectRabbitMQ() (*amqp091.Connection, error) {
	if conn != nil && !conn.IsClosed() {
		return conn, nil
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Retrying to connect to RabbitMQ (%d/5)...", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}
