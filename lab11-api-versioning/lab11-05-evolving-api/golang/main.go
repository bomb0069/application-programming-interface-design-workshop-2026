package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lib/pq"
)

var db *sql.DB

// Product represents the evolved product response.
// Fields are only added, never removed.
// Deprecated fields are kept for backward compatibility.
type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`             // Deprecated: use Categories instead
	Description string   `json:"description,omitempty"` // Added in evolution 1
	Tags        []string `json:"tags,omitempty"`        // Added in evolution 2
	SKU         string   `json:"sku,omitempty"`         // Added in evolution 3
	Categories  []string `json:"categories,omitempty"`  // Added in evolution 4 (replaces category)
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
	r.Use(deprecationHeaderMiddleware)

	r.Get("/api/products", listProducts)
	r.Get("/api/products/{id}", getProduct)
	r.Post("/api/products", createProduct)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// deprecationHeaderMiddleware adds sunset headers for deprecated fields
func deprecationHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Deprecated-Fields", "category")
		w.Header().Set("X-Deprecated-Message", "The 'category' field is deprecated. Use 'categories' array instead.")
		w.Header().Set("X-API-Sunset", "2026-12-31")
		next.ServeHTTP(w, r)
	})
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, category, description, tags, sku FROM products ORDER BY id")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, pq.Array(&p.Tags), &p.SKU)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
		// Evolution: populate categories from category
		p.Categories = []string{p.Category}
		products = append(products, p)
	}
	writeJSON(w, http.StatusOK, products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var p Product
	err = db.QueryRow("SELECT id, name, price, category, description, tags, sku FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, pq.Array(&p.Tags), &p.SKU)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}
	p.Categories = []string{p.Category}
	writeJSON(w, http.StatusOK, p)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string   `json:"name"`
		Price       float64  `json:"price"`
		Category    string   `json:"category"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		SKU         string   `json:"sku"`
		Categories  []string `json:"categories"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	if input.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name is required"})
		return
	}
	if input.Tags == nil {
		input.Tags = []string{}
	}
	// Support both category and categories
	category := input.Category
	if len(input.Categories) > 0 {
		category = input.Categories[0]
	}
	if input.SKU == "" {
		input.SKU = "N/A"
	}

	var p Product
	err := db.QueryRow(
		"INSERT INTO products (name, price, category, description, tags, sku) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, price, category, description, tags, sku",
		input.Name, input.Price, category, input.Description, pq.Array(input.Tags), input.SKU,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, pq.Array(&p.Tags), &p.SKU)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	p.Categories = []string{p.Category}

	writeJSON(w, http.StatusCreated, p)
}

func createTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL DEFAULT '',
		description TEXT DEFAULT '',
		tags TEXT[] DEFAULT '{}',
		sku TEXT DEFAULT 'N/A'
	)`)

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO products (name, price, category, description, tags, sku) VALUES
			('Laptop', 999.99, 'electronics', 'A powerful laptop for developers', '{portable,computing}', 'ELEC-001'),
			('Go Book', 39.99, 'books', 'Learn Go programming', '{programming,education}', 'BOOK-001'),
			('T-Shirt', 19.99, 'clothing', 'Comfortable cotton t-shirt', '{casual,cotton}', 'CLTH-001')`)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
