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
	_ "github.com/lib/pq"
)

// @title Products API
// @version 1.0
// @description A sample Products API for the API Design Workshop
// @host localhost:8080
// @BasePath /

var db *sql.DB

// Product represents a product in the store
type Product struct {
	ID       int     `json:"id" example:"1"`
	Name     string  `json:"name" example:"Laptop"`
	Price    float64 `json:"price" example:"999.99"`
	Category string  `json:"category" example:"electronics"`
}

// CreateProductInput represents the input for creating a product
type CreateProductInput struct {
	Name     string  `json:"name" example:"Laptop"`
	Price    float64 `json:"price" example:"999.99"`
	Category string  `json:"category" example:"electronics"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Code    string `json:"code" example:"NOT_FOUND"`
	Message string `json:"message" example:"Product not found"`
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

	createTable()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/products", listProducts)
	r.Post("/products", createProduct)
	r.Get("/products/{id}", getProduct)
	r.Put("/products/{id}", updateProduct)
	r.Delete("/products/{id}", deleteProduct)

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL
	)`
	db.Exec(query)
}

// listProducts godoc
// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func listProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, category FROM products ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category)
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// createProduct godoc
// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductInput true "Product to create"
// @Success 201 {object} Product
// @Failure 400 {object} ErrorResponse
// @Router /products [post]
func createProduct(w http.ResponseWriter, r *http.Request) {
	var input CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	if input.Name == "" || input.Price <= 0 || input.Category == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Name, price (>0), and category are required")
		return
	}

	var product Product
	err := db.QueryRow(
		"INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id, name, price, category",
		input.Name, input.Price, input.Category,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create product")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// getProduct godoc
// @Summary Get a product
// @Description Get a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ID")
		return
	}

	var product Product
	err = db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// updateProduct godoc
// @Summary Update a product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body CreateProductInput true "Updated product data"
// @Success 200 {object} Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [put]
func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ID")
		return
	}

	var input CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	var product Product
	err = db.QueryRow(
		"UPDATE products SET name=$1, price=$2, category=$3 WHERE id=$4 RETURNING id, name, price, category",
		input.Name, input.Price, input.Category, id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// deleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [delete]
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ID")
		return
	}

	result, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{Code: code, Message: message},
	})
}
