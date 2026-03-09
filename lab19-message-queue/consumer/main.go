package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	ID        string  `json:"id"`
	Item      string  `json:"item"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	var conn *amqp.Connection
	var err error
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Waiting for RabbitMQ... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// Declare exchange
	ch.ExchangeDeclare("orders", "topic", true, false, false, false, nil)

	// Declare queue
	q, err := ch.QueueDeclare(
		"order_processing", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,    // queue name
		"order.*", // routing key pattern
		"orders",  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to bind queue:", err)
	}

	// Set prefetch
	ch.Qos(1, 0, false)

	msgs, err := ch.Consume(
		q.Name,           // queue
		"order-consumer", // consumer tag
		false,            // auto-ack (false = manual ack)
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	for msg := range msgs {
		var order Order
		if err := json.Unmarshal(msg.Body, &order); err != nil {
			log.Printf("Error parsing message: %v", err)
			msg.Nack(false, false)
			continue
		}

		log.Printf("Processing order: %s - %s ($%.2f)", order.ID, order.Item, order.Amount)

		// Simulate processing time
		time.Sleep(2 * time.Second)

		log.Printf("Order processed: %s", order.ID)
		msg.Ack(false)
	}
}
