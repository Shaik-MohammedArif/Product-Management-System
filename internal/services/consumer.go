package services

import (
	"log"

	"github.com/streadway/amqp"
)

func connectToRabbitMQ() (*amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // RabbitMQ default connection string
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return ch, conn, nil
}

func Consumer() {
	// Connect to RabbitMQ
	ch, conn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer ch.Close()

	// Declare the queue in RabbitMQ
	_, err = ch.QueueDeclare(
		"image_processing_queue", // Queue name
		true,                    // Durable
		false,                    // Auto delete
		false,                    // Exclusive
		false,                    // No wait
		nil,                      // Arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start consuming messages from the queue
	msgs, err := ch.Consume(
		"image_processing_queue", // Queue name
		"",                       // Consumer tag
		true,                     // Auto-acknowledge
		false,                    // Exclusive
		false,                    // No-local
		false,                    // No-wait
		nil,                      // Arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Go routine to handle messages
	go func() {
		for msg := range msgs {
			log.Printf("Received image URL: %s", msg.Body)
			// Here you can add your image processing logic like downloading, resizing, etc.
		}
	}()

	// Block forever to keep consuming messages
	select {}
}