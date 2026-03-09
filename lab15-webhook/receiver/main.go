package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var webhookSecret = "webhook-secret-key"

type ReceivedEvent struct {
	Event     string          `json:"event"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
	Signature string          `json:"signature"`
	Valid     bool            `json:"valid"`
}

var (
	mu     sync.RWMutex
	events []ReceivedEvent
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("Receiver starting on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Verify signature
	signature := r.Header.Get("X-Webhook-Signature")
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	valid := hmac.Equal([]byte(signature), []byte(expectedSig))

	event := r.Header.Get("X-Webhook-Event")

	var payload struct {
		Event     string          `json:"event"`
		Timestamp string          `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	}
	json.Unmarshal(body, &payload)

	received := ReceivedEvent{
		Event:     event,
		Timestamp: payload.Timestamp,
		Data:      payload.Data,
		Signature: signature,
		Valid:     valid,
	}

	mu.Lock()
	events = append(events, received)
	mu.Unlock()

	if valid {
		log.Printf("Received valid webhook: %s", event)
		fmt.Fprintf(w, `{"status":"accepted"}`)
	} else {
		log.Printf("Received webhook with INVALID signature: %s", event)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"status":"invalid signature"}`)
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
