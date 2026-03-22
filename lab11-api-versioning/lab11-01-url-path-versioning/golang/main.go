package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

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

	// Discovery endpoint — lists all versioning variations
	r.Get("/", discoveryHandler)

	// Variation 1: Version as root prefix — /v1/api/products
	r.Route("/v1/api", func(r chi.Router) {
		r.Get("/products", v1ListProducts)
		r.Get("/products/{id}", v1GetProduct)
		r.Post("/products", v1CreateProduct)
	})
	r.Route("/v2/api", func(r chi.Router) {
		r.Get("/products", v2ListProducts)
		r.Get("/products/{id}", v2GetProduct)
		r.Post("/products", v2CreateProduct)
	})

	// Variation 2: Version after /api (RECOMMENDED) — /api/v1/products
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/products", v1ListProducts)
		r.Get("/products/{id}", v1GetProduct)
		r.Post("/products", v1CreateProduct)
	})
	r.Route("/api/v2", func(r chi.Router) {
		r.Get("/products", v2ListProducts)
		r.Get("/products/{id}", v2GetProduct)
		r.Post("/products", v2CreateProduct)
	})

	// Variation 3: Version as resource suffix — /api/products/v1
	r.Route("/api/products/v1", func(r chi.Router) {
		r.Get("/", v1ListProducts)
		r.Get("/{id}", v1GetProduct)
		r.Post("/", v1CreateProduct)
	})
	r.Route("/api/products/v2", func(r chi.Router) {
		r.Get("/", v2ListProducts)
		r.Get("/{id}", v2GetProduct)
		r.Post("/", v2CreateProduct)
	})

	// Variation 4: Version baked into name (anti-pattern) — /api/products-v1
	r.Route("/api/products-v1", func(r chi.Router) {
		r.Get("/", v1ListProducts)
		r.Get("/{id}", v1GetProduct)
		r.Post("/", v1CreateProduct)
	})
	r.Route("/api/products-v2", func(r chi.Router) {
		r.Get("/", v2ListProducts)
		r.Get("/{id}", v2GetProduct)
		r.Post("/", v2CreateProduct)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	discovery := map[string]interface{}{
		"service": "URL Path Versioning Demo",
		"variations": []map[string]interface{}{
			{
				"name":        "Root Prefix",
				"pattern":     "/v{version}/api/products",
				"example_v1":  "/v1/api/products",
				"example_v2":  "/v2/api/products",
				"recommended": false,
				"notes":       "Version before /api — uncommon, breaks gateway routing rules",
			},
			{
				"name":        "After /api (Recommended)",
				"pattern":     "/api/v{version}/products",
				"example_v1":  "/api/v1/products",
				"example_v2":  "/api/v2/products",
				"recommended": true,
				"notes":       "Industry standard — used by GitHub, Stripe, Twilio",
			},
			{
				"name":        "Resource Suffix",
				"pattern":     "/api/products/v{version}",
				"example_v1":  "/api/products/v1",
				"example_v2":  "/api/products/v2",
				"recommended": false,
				"notes":       "Clashes with sub-resources — /api/products/v1 looks like a nested resource",
			},
			{
				"name":        "Baked into Name (Anti-pattern)",
				"pattern":     "/api/products-v{version}",
				"example_v1":  "/api/products-v1",
				"example_v2":  "/api/products-v2",
				"recommended": false,
				"notes":       "Not real versioning — just different resource names. No framework support.",
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
