package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

var (
	db       *sql.DB
	validate *validator.Validate
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Category string  `json:"category" validate:"required,oneof=electronics books clothing food"`
	SKU      string  `json:"sku" validate:"required,len=8,alphanum"`
}

type CreateProductInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Category string  `json:"category" validate:"required,oneof=electronics books clothing food"`
	SKU      string  `json:"sku" validate:"required,len=8,alphanum"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func main() {
	validate = validator.New()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTable()

	r := chi.NewRouter()
	r.Get("/products", listProducts)
	r.Post("/products", createProduct)
	r.Get("/products/{id}", getProduct)
	r.Put("/products/{id}", updateProduct)
	r.Delete("/products/{id}", deleteProduct)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL,
		sku VARCHAR(8) UNIQUE NOT NULL
	)`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func formatValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		var msg string
		switch e.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", strings.ToLower(e.Field()))
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", strings.ToLower(e.Field()), e.Param())
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", strings.ToLower(e.Field()), e.Param())
		case "gt":
			msg = fmt.Sprintf("%s must be greater than %s", strings.ToLower(e.Field()), e.Param())
		case "oneof":
			msg = fmt.Sprintf("%s must be one of: %s", strings.ToLower(e.Field()), e.Param())
		case "len":
			msg = fmt.Sprintf("%s must be exactly %s characters", strings.ToLower(e.Field()), e.Param())
		case "alphanum":
			msg = fmt.Sprintf("%s must contain only alphanumeric characters", strings.ToLower(e.Field()))
		default:
			msg = fmt.Sprintf("%s is invalid", strings.ToLower(e.Field()))
		}
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(e.Field()),
			Message: msg,
		})
	}
	return errors
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, category, sku FROM products ORDER BY id")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.SKU); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		products = append(products, p)
	}
	respondJSON(w, http.StatusOK, products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var input CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := validate.Struct(input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": formatValidationErrors(err),
		})
		return
	}

	var product Product
	err := db.QueryRow(
		"INSERT INTO products (name, price, category, sku) VALUES ($1, $2, $3, $4) RETURNING id, name, price, category, sku",
		input.Name, input.Price, input.Category, input.SKU,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.SKU)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondJSON(w, http.StatusConflict, map[string]string{"error": "Product with this SKU already exists"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var product Product
	err = db.QueryRow("SELECT id, name, price, category, sku FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.SKU)
	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var input CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := validate.Struct(input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": formatValidationErrors(err),
		})
		return
	}

	var product Product
	err = db.QueryRow(
		"UPDATE products SET name=$1, price=$2, category=$3, sku=$4 WHERE id=$5 RETURNING id, name, price, category, sku",
		input.Name, input.Price, input.Category, input.SKU, id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.SKU)
	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	result, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
