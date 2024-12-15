# Product Management System with Asynchronous Image Processing

## Project Overview

This project is a backend system built in Go for a Product Management application, focusing on asynchronous image processing, caching, logging, and scalability. The application provides an API to manage products, with functionality for handling product data, image compression, and caching of product details to improve performance.

The system architecture is designed for high scalability and performance, using PostgreSQL for data storage, Redis for caching, and RabbitMQ (or Kafka) for asynchronous image processing.

## Features

- **RESTful API** for managing product data:
  - `POST /products`: Accepts product details with multiple image URLs for processing
  - `GET /products/:id`: Retrieves product details by ID, including compressed image data
  - `GET /products`: Retrieves all products for a user, with optional filtering by price and product name
- **Asynchronous Image Processing**: Processes image URLs using RabbitMQ or Kafka for offloading the image compression task
- **Caching**: Product data is cached in Redis to reduce load on the database for frequently requested products
- **Logging**: Structured logging is implemented to track API requests, image processing events, and errors

## System Architecture

The architecture follows a modular approach, separating the components for API handling, image processing, caching, and logging. The system consists of the following key components:

- **API Service**: A Go-based server providing the RESTful API.
- **Image Processing Service**: A background service that processes image URLs sent through a message queue.
- **Caching Service**: Redis is used for caching product data and reducing database load.
- **Logging**: Structured logging using libraries like `logrus` or `zap` to capture relevant events and metrics.

### Flow of Data

1. **Product Creation**: When a user adds a new product, the product data, including image URLs, is stored in the PostgreSQL database. These image URLs are then pushed to a message queue (RabbitMQ or Kafka).
2. **Image Processing**: A consumer service listens to the queue, downloads the images, compresses them, and stores the compressed images (e.g., in S3). Once the processing is complete, the `compressed_product_images` field in the database is updated.
3. **Caching**: Product details are cached in Redis. When a product is requested, the system first checks if the data is available in the cache, reducing database load.
4. **Logging**: All actions, including API requests, image processing tasks, and failures, are logged in a structured format to help with monitoring and debugging.

## API Endpoints

### `POST /products`
Creates a new product with the given data, including image URLs.

**Request Body:**
```json
{
  "user_id": "12345",
  "product_name": "Wireless Mouse",
  "product_description": "A high-quality wireless mouse",
  "product_images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "product_price": 25.99
}

# Response:

### POST /products
Creates a new product with the given data, including image URLs.

**Response:**
```json
{
  "status": "success",
  "message": "Product created successfully"
}
GET /products/:id
Retrieves product details by ID, including compressed images if available.

Response:
{
  "product_id": "12345",
  "user_id": "12345",
  "product_name": "Wireless Mouse",
  "product_description": "A high-quality wireless mouse",
  "product_images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "compressed_product_images": [
    "https://s3.amazonaws.com/compressed/image1.jpg",
    "https://s3.amazonaws.com/compressed/image2.jpg"
  ],
  "product_price": 25.99
}

GET /products
Retrieves all products for a specific user_id, with optional filters by price_range and product_name.

Query Parameters:

user_id: ID of the user to fetch products for
price_range: Optional filter by price range (e.g., price_range=20-50)
product_name: Optional filter by product name
Response:

[
  {
    "product_id": "12345",
    "user_id": "12345",
    "product_name": "Wireless Mouse",
    "product_description": "A high-quality wireless mouse",
    "product_images": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ],
    "product_price": 25.99
  },
  ...
]

Data Storage
PostgreSQL is used to store user and product information. The products table includes fields for product data such as name, description, price, and image URLs, as well as a compressed_product_images column to store the links to the processed images.

Example PostgreSQL Schema
CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL
);

CREATE TABLE products (
  product_id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(user_id),
  product_name VARCHAR(255),
  product_description TEXT,
  product_images TEXT[],
  compressed_product_images TEXT[],
  product_price DECIMAL
);
## Asynchronous Image Processing

### Workflow
- **Producer**: When a new product is created, the image URLs are published to a message queue (RabbitMQ or Kafka).
- **Consumer**: A separate service consumes the image URLs from the queue, processes them (e.g., downloading and compressing images), and stores the compressed images (e.g., in Amazon S3).
- Once the processing is complete, the `compressed_product_images` field in the database is updated with the URLs of the processed images.

### Message Queue
RabbitMQ or Kafka are used as message queues to decouple the image processing task from the main API flow. The system pushes image URLs to the queue, and the image processing service consumes them asynchronously.

## Caching
Redis is used to cache product data for faster retrieval.

- When a product's details are requested, the system first checks if the data exists in the Redis cache. If it does, the cached data is returned. Otherwise, the data is fetched from the PostgreSQL database, cached, and returned to the user.

## Logging
Structured logging is implemented using **Logrus** or **Zap**.

- Logs include details about API requests, image processing events, errors, and performance metrics.
- Logs are stored in a central log management system (optional for production).

## Error Handling
- **Asynchronous Failures**: If the image processing task fails (e.g., image download or compression failure), the system implements retries or places the task in a dead-letter queue for further inspection.
- **API Errors**: Proper HTTP error responses are provided in case of invalid requests or server errors.

## Testing
- **Unit Tests**: Each API endpoint and core function is covered with unit tests.
- **Integration Tests**: Tests are written to ensure end-to-end functionality, especially for the image processing workflow and cache functionality.
- **Benchmarking**: Tests are written to measure the performance of the `GET /products/:id` endpoint, comparing the response time with and without cache hits.

## Setup and Configuration

### Environment Setup
- **PostgreSQL**: Set up a PostgreSQL instance and create the `users` and `products` tables as shown above.
- **Redis**: Install and configure Redis for caching.
- **RabbitMQ**: Install and configure RabbitMQ (or Kafka) for message queuing.
- **Image Storage**: Set up a storage system (e.g., Amazon S3) for storing compressed images.

### Configuration Files
- **.env**: Contains configuration values for PostgreSQL, Redis, RabbitMQ, and image storage.

#### Example `.env` file:
```bash
POSTGRES_URL=postgres://user:password@localhost:5432/dbname
REDIS_URL=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
IMAGE_STORAGE_URL=s3://bucket-name/
