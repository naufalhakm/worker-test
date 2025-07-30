package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-worker/internal/commons"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type VerifyPayload struct {
	UserID uint64 `json:"user_id"`
}

type Consumer struct {
	queueName string
	worker    int
	handler   func(ctx context.Context, payload VerifyPayload) error
}

func NewConsumer(queueName string, worker int, handler func(ctx context.Context, payload VerifyPayload) error) *Consumer {
	return &Consumer{
		queueName: queueName,
		worker:    worker,
		handler:   handler,
	}
}

// func (c *Consumer) Start(ctx context.Context) error {
// 	start := time.Now()
// 	log.Printf("Consumer started at: %v", start)

// 	defer func() {
// 		duration := time.Since(start)
// 		log.Printf("Consumer stopped after: %v", duration)
// 	}()

// 	conn, err := ConnectRabbitMQ()
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return fmt.Errorf("failed to open a channel: %w", err)
// 	}
// 	defer ch.Close()

// 	msgs, err := ch.Consume(
// 		c.queueName,
// 		"",    // consumer tag
// 		true,  // auto-ack
// 		false, // not exclusive
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to register consumer: %w", err)
// 	}

// 	log.Printf("Consumer started with %d workers", c.worker)

// 	workerPool := make(chan struct{}, c.worker)
// 	var wg sync.WaitGroup

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Println("Consumer shutting down gracefully")
// 			wg.Wait()
// 			return nil
// 		case d, ok := <-msgs:
// 			if !ok {
// 				log.Println("Message channel closed")
// 				wg.Wait()
// 				return nil
// 			}

// 			workerPool <- struct{}{}
// 			wg.Add(1)

// 			go func(d amqp091.Delivery) {
// 				defer func() {
// 					wg.Done()
// 					<-workerPool
// 				}()

// 				var payload VerifyPayload
// 				if err := json.Unmarshal(d.Body, &payload); err != nil {
// 					log.Printf("Failed to unmarshal payload: %v", err)
// 					return
// 				}

// 				// log.Printf("Processing user_id: %d", payload.UserID)
// 				if err := c.handler(ctx, payload); err != nil {
// 					log.Printf("Error handling user_id %d: %v", payload.UserID, err)
// 				}
// 			}(d)
// 		}
// 	}
// }

func (c *Consumer) Start(ctx context.Context) error {
	start := time.Now()
	log.Printf("Consumer started at: %v", start)

	defer func() {
		duration := time.Since(start)
		log.Printf("Consumer stopped after: %v", duration)
	}()

	conn, err := ConnectRabbitMQ()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	err = ch.Qos(
		60,    // prefetch count - limit messages per consumer
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("📋 Consumer started with %d workers, initial goroutines: %d", c.worker, runtime.NumGoroutine())

	workerPool := make(chan struct{}, c.worker)
	var wg sync.WaitGroup

	for i := 0; i < c.worker; i++ {
		workerPool <- struct{}{}
	}
	messageCount := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutting down gracefully")
			wg.Wait()
			return nil
		case d, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				wg.Wait()
				return nil
			}
			messageCount++

			<-workerPool
			wg.Add(1)

			go func(delivery amqp091.Delivery, msgNum int) {
				defer func() {
					workerPool <- struct{}{}
					wg.Done()
				}()
				currentGoroutines := runtime.NumGoroutine()
				log.Printf("Worker %d started (msg #%d), goroutines: %d", msgNum, msgNum, currentGoroutines)

				if currentGoroutines > 200 {
					log.Printf("HIGH GOROUTINE COUNT: %d", currentGoroutines)
				}

				c.processMessage(ctx, delivery, msgNum)
			}(d, messageCount)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, d amqp091.Delivery, msgNum int) {
	start := time.Now()
	defer func() {
		log.Printf("Worker %d finished in %v, goroutines: %d", msgNum, time.Since(start), runtime.NumGoroutine())
	}()

	var payload VerifyPayload
	if err := json.Unmarshal(d.Body, &payload); err != nil {
		log.Printf("Failed to unmarshal payload: %v", err)
		d.Nack(false, false)
		return
	}

	log.Printf("Processing user_id: %d", payload.UserID)

	if err := c.handler(ctx, payload); err != nil {
		duration := time.Since(start)
		log.Printf("Error handling user_id %d after %v: %v", payload.UserID, duration, err)
		var perr *commons.PermanentError
		if errors.As(err, &perr) {
			log.Printf("Permanent error for user_id %d: %v (no requeue)", payload.UserID, perr)
			d.Nack(false, false) // tidak requeue
		} else {
			log.Printf("Temporary error for user_id %d: %v (will requeue)", payload.UserID, err)
			d.Nack(false, true) // requeue
		}
		return
	}

	duration := time.Since(start)
	log.Printf("Successfully processed user_id %d in %v", payload.UserID, duration)
	d.Ack(false)
}

// func (c *Consumer) Start(ctx context.Context) error {
// 	conn, err := ConnectRabbitMQ()
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return fmt.Errorf("failed to open a channel: %w", err)
// 	}
// 	defer ch.Close()

// 	msgs, err := ch.Consume(
// 		c.queueName,
// 		"",
// 		true,  // auto-ack
// 		false, // not exclusive
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to register consumer: %w", err)
// 	}
// 	workerID := 1
// 	// var wg sync.WaitGroup
// 	// for i := 0; i < c.worker; i++ {
// 	// 	wg.Add(1)
// 	// 	go func(workerID int) {
// 	// 		defer wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Printf("[Worker %d] Stopping consumer", workerID)
// 			// return
// 		case d, ok := <-msgs:
// 			if !ok {
// 				log.Printf("[Worker %d] Channel closed", workerID)
// 				// return
// 			}

// 			var payload VerifyPayload
// 			if err := json.Unmarshal(d.Body, &payload); err != nil {
// 				log.Printf("[Worker %d] Failed to unmarshal payload: %v", workerID, err)
// 				continue
// 			}

// 			log.Printf("[Worker %d] Processing user_id: %d", workerID, payload.UserID)

// 			if err := c.handler(ctx, payload); err != nil {
// 				log.Printf("[Worker %d] Error handling message: %v", workerID, err)
// 			}
// 		}
// 	}
// 	// 	}(i + 1)
// 	// }

// 	// wg.Wait()
// 	return nil
// }
