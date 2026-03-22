package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

type contextKey string

const versionKey contextKey = "api-version"

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

	// Discovery endpoint
	r.Get("/", discoveryHandler)

	// Content-negotiation versioned routes
	r.Route("/api/products", func(r chi.Router) {
		r.Use(contentNegotiationMiddleware)
		r.Get("/", listProducts)
		r.Get("/{id}", getProduct)
		r.Post("/", createProduct)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// contentNegotiationMiddleware extracts the API version from the Accept header.
//
// Supported media types:
//   - application/vnd.workshop.v1+json  -> version "1"
//   - application/vnd.workshop.v2+json  -> version "2"
//   - application/json or */*           -> defaults to version "1"
//
// Any other Accept value returns 406 Not Acceptable.
func contentNegotiationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		version := "1" // default

		if strings.Contains(accept, "application/vnd.workshop.v2+json") {
			version = "2"
			w.Header().Set("Content-Type", "application/vnd.workshop.v2+json")
		} else if strings.Contains(accept, "application/vnd.workshop.v1+json") {
			version = "1"
			w.Header().Set("Content-Type", "application/vnd.workshop.v1+json")
		} else if accept != "" && accept != "*/*" && !strings.Contains(accept, "application/json") {
			// Unsupported media type
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Not Acceptable",
				"message": "Supported media types: application/vnd.workshop.v1+json, application/vnd.workshop.v2+json, application/json",
			})
			return
		}

		w.Header().Set("Vary", "Accept")
		ctx := context.WithValue(r.Context(), versionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Dispatcher handlers — route to v1 or v2 based on negotiated version

func listProducts(w http.ResponseWriter, r *http.Request) {
	version := r.Context().Value(versionKey).(string)
	switch version {
	case "2":
		v2ListProducts(w, r)
	default:
		v1ListProducts(w, r)
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	version := r.Context().Value(versionKey).(string)
	switch version {
	case "2":
		v2GetProduct(w, r)
	default:
		v1GetProduct(w, r)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	version := r.Context().Value(versionKey).(string)
	switch version {
	case "2":
		v2CreateProduct(w, r)
	default:
		v1CreateProduct(w, r)
	}
}

func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	discovery := map[string]interface{}{
		"service": "Content Negotiation Versioning Demo",
		"description": "Use the Accept header with vendor media types to select the API version",
		"media_types": map[string]interface{}{
			"v1": "application/vnd.workshop.v1+json",
			"v2": "application/vnd.workshop.v2+json",
		},
		"default": "v1 (when Accept is application/json or */*)",
		"examples": []map[string]string{
			{
				"description": "Request V1 explicitly",
				"curl":        `curl -H "Accept: application/vnd.workshop.v1+json" http://localhost:8080/api/products`,
			},
			{
				"description": "Request V2 explicitly",
				"curl":        `curl -H "Accept: application/vnd.workshop.v2+json" http://localhost:8080/api/products`,
			},
			{
				"description": "Default to V1 (standard JSON)",
				"curl":        `curl -H "Accept: application/json" http://localhost:8080/api/products`,
			},
			{
				"description": "Unsupported media type (returns 406)",
				"curl":        `curl -H "Accept: text/xml" http://localhost:8080/api/products`,
			},
		},
	}
	writeJSON(w, http.StatusOK, discovery)
}

func createTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL,
		description TEXT DEFAULT '',
		tags TEXT[] DEFAULT '{}'
	)`)

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO products (name, price, category, description, tags) VALUES
			('Laptop', 999.99, 'electronics', 'A powerful laptop for developers', '{portable,computing}'),
			('Go Book', 39.99, 'books', 'Learn Go programming', '{programming,education}'),
			('T-Shirt', 19.99, 'clothing', 'Comfortable cotton t-shirt', '{casual,cotton}')`)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
