package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var webhookSecret = "webhook-secret-key"

type Order struct {
	ID     int     `json:"id"`
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type WebhookRegistration struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Events string `json:"events"`
	Active bool   `json:"active"`
}

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Ping()
	createTables()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Order endpoints
	r.Post("/orders", createOrder)
	r.Get("/orders", listOrders)
	r.Put("/orders/{id}/status", updateOrderStatus)

	// Webhook registration
	r.Post("/webhooks", registerWebhook)
	r.Get("/webhooks", listWebhooks)
	r.Delete("/webhooks/{id}", deleteWebhook)

	log.Println("Sender (Order Service) starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		item TEXT NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		status TEXT DEFAULT 'pending'
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS webhooks (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		events TEXT NOT NULL DEFAULT '*',
		active BOOLEAN DEFAULT TRUE
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS webhook_logs (
		id SERIAL PRIMARY KEY,
		webhook_id INT REFERENCES webhooks(id),
		event TEXT NOT NULL,
		payload TEXT NOT NULL,
		status_code INT,
		response TEXT,
		attempt INT DEFAULT 1,
		sent_at TIMESTAMP DEFAULT NOW()
	)`)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Item   string  `json:"item"`
		Amount float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	if input.Item == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Item is required"})
		return
	}

	var order Order
	db.QueryRow("INSERT INTO orders (item, amount) VALUES ($1, $2) RETURNING id, item, amount, status",
		input.Item, input.Amount).Scan(&order.ID, &order.Item, &order.Amount, &order.Status)

	go sendWebhooks("order.created", order)
	writeJSON(w, http.StatusCreated, order)
}

func listOrders(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, item, amount, status FROM orders ORDER BY id")
	defer rows.Close()
	orders := []Order{}
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.Item, &o.Amount, &o.Status)
		orders = append(orders, o)
	}
	writeJSON(w, http.StatusOK, orders)
}

func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var input struct {
		Status string `json:"status"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	validStatuses := map[string]bool{"pending": true, "confirmed": true, "shipped": true, "delivered": true, "cancelled": true}
	if !validStatuses[input.Status] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid status"})
		return
	}

	var order Order
	err := db.QueryRow("UPDATE orders SET status=$1 WHERE id=$2 RETURNING id, item, amount, status",
		input.Status, id).Scan(&order.ID, &order.Item, &order.Amount, &order.Status)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	}

	go sendWebhooks("order."+input.Status, order)
	writeJSON(w, http.StatusOK, order)
}

func registerWebhook(w http.ResponseWriter, r *http.Request) {
	var input struct {
		URL    string `json:"url"`
		Events string `json:"events"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	if input.URL == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "URL is required"})
		return
	}
	if input.Events == "" {
		input.Events = "*"
	}

	var wh WebhookRegistration
	db.QueryRow("INSERT INTO webhooks (url, events) VALUES ($1, $2) RETURNING id, url, events, active",
		input.URL, input.Events).Scan(&wh.ID, &wh.URL, &wh.Events, &wh.Active)
	writeJSON(w, http.StatusCreated, wh)
}

func listWebhooks(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, url, events, active FROM webhooks ORDER BY id")
	defer rows.Close()
	webhooks := []WebhookRegistration{}
	for rows.Next() {
		var wh WebhookRegistration
		rows.Scan(&wh.ID, &wh.URL, &wh.Events, &wh.Active)
		webhooks = append(webhooks, wh)
	}
	writeJSON(w, http.StatusOK, webhooks)
}

func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	db.Exec("DELETE FROM webhooks WHERE id = $1", id)
	w.WriteHeader(http.StatusNoContent)
}

func sendWebhooks(event string, data interface{}) {
	rows, _ := db.Query("SELECT id, url FROM webhooks WHERE active = TRUE")
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		rows.Scan(&id, &url)

		payload := WebhookPayload{
			Event:     event,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Data:      data,
		}

		go deliverWebhook(id, url, event, payload)
	}
}

func deliverWebhook(webhookID int, url string, event string, payload WebhookPayload) {
	body, _ := json.Marshal(payload)

	// Create HMAC signature
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Webhook-Event", event)
		req.Header.Set("X-Webhook-Signature", signature)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)

		statusCode := 0
		respBody := ""
		if err != nil {
			respBody = err.Error()
		} else {
			statusCode = resp.StatusCode
			resp.Body.Close()
		}

		db.Exec("INSERT INTO webhook_logs (webhook_id, event, payload, status_code, response, attempt) VALUES ($1, $2, $3, $4, $5, $6)",
			webhookID, event, string(body), statusCode, respBody, attempt)

		if err == nil && statusCode >= 200 && statusCode < 300 {
			log.Printf("Webhook delivered: %s -> %s (attempt %d)", event, url, attempt)
			return
		}

		log.Printf("Webhook failed: %s -> %s (attempt %d, status: %d)", event, url, attempt, statusCode)
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
