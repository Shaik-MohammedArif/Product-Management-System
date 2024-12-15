package handlers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes sets up all the application routes
func RegisterRoutes(router *mux.Router, db *sql.DB, redisClient *redis.Client, rabbitMQChannel *amqp091.Channel) {
	// Route for getting all products
	router.HandleFunc("/products", GetAllProductsHandler(db)).Methods("GET")

	// Route for getting product by ID
	router.HandleFunc("/products/{id}", GetProductByIDHandler(db)).Methods("GET")

	// Route for creating a product
	router.HandleFunc("/products", CreateProductHandler(db)).Methods("POST")

	// Route for updating a product by ID
	router.HandleFunc("/products/{id}", UpdateProductHandler(db)).Methods("PUT")

	// Route for deleting a product by ID
	router.HandleFunc("/products/{id}", DeleteProductHandler(db)).Methods("DELETE")
}



