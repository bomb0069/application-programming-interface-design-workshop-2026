package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type contextKey string

const versionKey contextKey = "api-version"

var db *sql.DB

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
	r.Use(headerVersionMiddleware)

	// Single URL — version is determined by the X-Api-Version request header
	r.Get("/api/products", listProducts)
	r.Get("/api/products/{id}", getProduct)
	r.Post("/api/products", createProduct)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// headerVersionMiddleware extracts X-Api-Version from the request header,
// sets confirmation and caching headers on the response, and stores the
// version in context for downstream handlers.
func headerVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get("X-Api-Version")
		if version == "" {
			version = "1" // default to v1 when header is missing
		}

		// Confirm the resolved version back to the caller
		w.Header().Set("X-Api-Version", version)
		// Instruct caches that the response varies by this header
		w.Header().Set("Vary", "X-Api-Version")

		ctx := context.WithValue(r.Context(), versionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getVersion retrieves the API version stored in the request context.
func getVersion(r *http.Request) string {
	if v, ok := r.Context().Value(versionKey).(string); ok {
		return v
	}
	return "1"
}

// ---------------------------------------------------------------------------
// Dispatcher handlers — route to v1 or v2 logic based on the resolved version
// ---------------------------------------------------------------------------

func listProducts(w http.ResponseWriter, r *http.Request) {
	switch getVersion(r) {
	case "2":
		v2ListProducts(w, r)
	default:
		v1ListProducts(w, r)
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	switch getVersion(r) {
	case "2":
		v2GetProduct(w, r)
	default:
		v1GetProduct(w, r)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	switch getVersion(r) {
	case "2":
		v2CreateProduct(w, r)
	default:
		v1CreateProduct(w, r)
	}
}

// ---------------------------------------------------------------------------
// Database helpers
// ---------------------------------------------------------------------------

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// parseID extracts and validates the {id} URL parameter.
func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return 0, false
	}
	return id, true
}
