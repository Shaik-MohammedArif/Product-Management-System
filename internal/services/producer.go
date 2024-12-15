package services

import (
	"log"
	"database/sql"
	"github.com/streadway/amqp"
	_ "github.com/lib/pq"
)

// Producer sends product image URLs to RabbitMQ
func Producer() {
	// Connect to the PostgreSQL database
	connStr := "user=postgres  password=Shaik@2003 dbname=assignment sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query product details
	rows, err := db.Query("SELECT id, product_images FROM products WHERE product_images IS NOT NULL")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"image_processing_queue",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Process each product image URL and send it to the queue
	for rows.Next() {
		var id int
		var imageUrl string
		err := rows.Scan(&id, &imageUrl)
		if err != nil {
			log.Fatal(err)
		}

		err = ch.Publish(
			"", q.Name, false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(imageUrl),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Sent image URL: %s", imageUrl)
	}
}
