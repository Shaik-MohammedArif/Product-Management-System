package main

import (
	"log"
	"net/http"

	"assignment/internal/cache"
	"assignment/internal/handlers"
	"assignment/internal/queue"
	"assignment/internal/storage"
	"assignment/internal/services"
	"github.com/gorilla/mux"
	
)

func main() {
	// Initialize PostgreSQL connection
	db, err := storage.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis cache
	redisClient := cache.InitRedis()
	defer redisClient.Close()

	// Initialize RabbitMQ
	rabbitMQConn, rabbitMQChannel, err := queue.InitRabbitMQ() // Capture error here
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err) // Handle the error
	}
	defer rabbitMQConn.Close()
	defer rabbitMQChannel.Close()

	// Set up HTTP router
	router := mux.NewRouter()
	handlers.RegisterRoutes(router, db, redisClient, rabbitMQChannel)

	//start the consumer
	go services.Consumer()

	// Start Producer for asynchronous image processing
	go services.Producer()

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	services.Producer()
	

}
