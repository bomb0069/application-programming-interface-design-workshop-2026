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

	// URL path versioned routes (highest priority — version from URL)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(setVersion("1"))
		r.Get("/products", listProducts)
		r.Get("/products/{id}", getProduct)
		r.Post("/products", createProduct)
	})
	r.Route("/api/v2", func(r chi.Router) {
		r.Use(setVersion("2"))
		r.Get("/products", listProducts)
		r.Get("/products/{id}", getProduct)
		r.Post("/products", createProduct)
	})

	// Non-URL-versioned routes — fall back to query param, then header, then default
	r.Route("/api/products", func(r chi.Router) {
		r.Use(combinedVersionMiddleware)
		r.Get("/", listProducts)
		r.Get("/{id}", getProduct)
		r.Post("/", createProduct)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"service": "Combining Versioning Strategies Demo",
			"resolution_priority": []string{
				"1. URL path: /api/v1/products or /api/v2/products",
				"2. Query parameter: /api/products?api-version=1",
				"3. Header: X-Api-Version: 1",
				"4. Default: v1",
			},
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// setVersion is middleware that sets a specific version (used for URL path routes).
// Because the version is embedded in the URL, it takes the highest priority.
func setVersion(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Version-Source", "url-path")
			w.Header().Set("X-Api-Version", version)
			ctx := context.WithValue(r.Context(), versionKey, version)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// combinedVersionMiddleware resolves the API version from query param > header > default.
// This middleware is used for routes without a version segment in the URL path.
func combinedVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := ""
		source := "default"

		// Priority 2: Query parameter
		if qv := r.URL.Query().Get("api-version"); qv != "" {
			version = qv
			source = "query-parameter"
		}

		// Priority 3: Header (lower priority, only if query not set)
		if version == "" {
			if hv := r.Header.Get("X-Api-Version"); hv != "" {
				version = hv
				source = "header"
			}
		}

		// Default
		if version == "" {
			version = "1"
		}

		w.Header().Set("X-Version-Source", source)
		w.Header().Set("X-Api-Version", version)
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
