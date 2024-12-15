package worker

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

func ProcessImages(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // Queue
		"",        // Consumer
		true,      // Auto-acknowledge
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			imageURL := string(d.Body)
			compressedImage, err := compressImage(imageURL)
			if err != nil {
				log.Printf("Failed to process image: %v", err)
				continue
			}

			// Save compressed image to local file (for now; can upload to S3)
			err = os.WriteFile("compressed_image.jpg", compressedImage, 0644)
			if err != nil {
				log.Printf("Failed to save compressed image: %v", err)
			} else {
				log.Println("Compressed image saved successfully")
			}
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}

func compressImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 50}) // Compress to 50% quality
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
