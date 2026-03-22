package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

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
	r.Use(queryVersionMiddleware)

	// Single set of routes — version is resolved from ?api-version= query param
	r.Get("/api/products", listProducts)
	r.Get("/api/products/{id}", getProduct)
	r.Post("/api/products", createProduct)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// queryVersionMiddleware extracts api-version from query params and stores it in request context.
// Defaults to "1" when the parameter is omitted.
func queryVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.URL.Query().Get("api-version")
		if version == "" {
			version = "1" // default to v1
		}
		ctx := context.WithValue(r.Context(), versionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getVersion retrieves the API version string from request context.
func getVersion(r *http.Request) string {
	if v, ok := r.Context().Value(versionKey).(string); ok {
		return v
	}
	return "1"
}

// listProducts dispatches to v1 or v2 handler based on the version in context.
func listProducts(w http.ResponseWriter, r *http.Request) {
	switch getVersion(r) {
	case "2":
		v2ListProducts(w, r)
	default:
		v1ListProducts(w, r)
	}
}

// getProduct dispatches to v1 or v2 handler based on the version in context.
func getProduct(w http.ResponseWriter, r *http.Request) {
	switch getVersion(r) {
	case "2":
		v2GetProduct(w, r)
	default:
		v1GetProduct(w, r)
	}
}

// createProduct dispatches to v1 or v2 handler based on the version in context.
func createProduct(w http.ResponseWriter, r *http.Request) {
	switch getVersion(r) {
	case "2":
		v2CreateProduct(w, r)
	default:
		v1CreateProduct(w, r)
	}
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
