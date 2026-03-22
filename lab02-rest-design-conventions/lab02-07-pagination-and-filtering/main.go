package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	InStock  bool    `json:"in_stock"`
}

type PaginatedResponse struct {
	Data     []Product    `json:"data"`
	Metadata PageMetadata `json:"metadata"`
}

type PageMetadata struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
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

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTableAndSeed()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/products", listProducts)
	r.Get("/products/{id}", getProduct)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTableAndSeed() {
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL,
		in_stock BOOLEAN DEFAULT TRUE
	)`)

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 {
		return
	}

	products := []struct {
		name     string
		price    float64
		category string
		inStock  bool
	}{
		{"Laptop Pro 15", 1299.99, "electronics", true},
		{"Wireless Mouse", 29.99, "electronics", true},
		{"Mechanical Keyboard", 89.99, "electronics", true},
		{"USB-C Hub", 49.99, "electronics", false},
		{"Monitor 27 inch", 449.99, "electronics", true},
		{"Go Programming Language", 39.99, "books", true},
		{"Clean Code", 34.99, "books", true},
		{"Design Patterns", 44.99, "books", false},
		{"API Design Patterns", 49.99, "books", true},
		{"The Pragmatic Programmer", 42.99, "books", true},
		{"Cotton T-Shirt", 19.99, "clothing", true},
		{"Denim Jeans", 59.99, "clothing", true},
		{"Winter Jacket", 129.99, "clothing", false},
		{"Running Shoes", 89.99, "clothing", true},
		{"Baseball Cap", 14.99, "clothing", true},
		{"Organic Coffee", 12.99, "food", true},
		{"Green Tea Pack", 8.99, "food", true},
		{"Dark Chocolate", 5.99, "food", true},
		{"Protein Bars", 24.99, "food", false},
		{"Olive Oil", 15.99, "food", true},
	}

	for _, p := range products {
		db.Exec("INSERT INTO products (name, price, category, in_stock) VALUES ($1, $2, $3, $4)",
			p.name, p.price, p.category, p.inStock)
	}
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse filter params
	category := r.URL.Query().Get("category")
	inStockStr := r.URL.Query().Get("in_stock")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	// Parse sort params
	sortField := r.URL.Query().Get("sort")
	if sortField == "" {
		sortField = "id"
	}
	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Validate sort field
	validSortFields := map[string]bool{"id": true, "name": true, "price": true, "category": true}
	if !validSortFields[sortField] {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid sort field. Use: id, name, price, category"})
		return
	}

	// Build WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		where += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	if inStockStr != "" {
		inStock, err := strconv.ParseBool(inStockStr)
		if err == nil {
			where += fmt.Sprintf(" AND in_stock = $%d", argIdx)
			args = append(args, inStock)
			argIdx++
		}
	}
	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			where += fmt.Sprintf(" AND price >= $%d", argIdx)
			args = append(args, minPrice)
			argIdx++
		}
	}
	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			where += fmt.Sprintf(" AND price <= $%d", argIdx)
			args = append(args, maxPrice)
			argIdx++
		}
	}

	// Count total
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM products " + where
	db.QueryRow(countQuery, args...).Scan(&totalItems)

	// Fetch page
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT id, name, price, category, in_stock FROM products %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		where, sortField, sortOrder, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.InStock)
		products = append(products, p)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PaginatedResponse{
		Data: products,
		Metadata: PageMetadata{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	})
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var p Product
	err = db.QueryRow("SELECT id, name, price, category, in_stock FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.InStock)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
