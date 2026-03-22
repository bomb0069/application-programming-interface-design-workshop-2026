package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitConn *amqp.Connection
var rabbitCh *amqp.Channel

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

	var err error
	for i := 0; i < 30; i++ {
		rabbitConn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Waiting for RabbitMQ... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer rabbitCh.Close()

	// Declare exchange
	err = rabbitCh.ExchangeDeclare(
		"orders", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	http.HandleFunc("/orders", createOrderHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("Publisher (Order API) starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Item   string  `json:"item"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	order := Order{
		ID:        fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		Item:      input.Item,
		Amount:    input.Amount,
		Status:    "created",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	body, _ := json.Marshal(order)

	err := rabbitCh.Publish(
		"orders",        // exchange
		"order.created", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to publish message"})
		return
	}

	log.Printf("Published order: %s", order.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order accepted for processing",
		"order":   order,
	})
}
