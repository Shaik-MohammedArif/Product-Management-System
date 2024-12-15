package queue
import (
    "log"
    "github.com/rabbitmq/amqp091-go"
)

// InitRabbitMQ initializes a connection to RabbitMQ and returns the connection and channel
func InitRabbitMQ() (*amqp091.Connection, *amqp091.Channel, error) {
    conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
        return nil, nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
        return nil, nil, err
    }

    return conn,ch,nil
}