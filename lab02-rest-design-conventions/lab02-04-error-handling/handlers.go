package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CreateProductInput struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, category FROM products ORDER BY id")
	if err != nil {
		log.Printf("Database error: %v", err)
		NewInternalError().Send(w)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			log.Printf("Scan error: %v", err)
			NewInternalError().Send(w)
			return
		}
		products = append(products, p)
	}
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var input CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewBadRequestError("Invalid request body").Send(w)
		return
	}

	if input.Name == "" {
		NewBadRequestError("Name is required").Send(w)
		return
	}
	if input.Price <= 0 {
		NewBadRequestError("Price must be greater than 0").Send(w)
		return
	}
	if input.Category == "" {
		NewBadRequestError("Category is required").Send(w)
		return
	}

	var product Product
	err := db.QueryRow(
		"INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id, name, price, category",
		input.Name, input.Price, input.Category,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			NewConflictError("A product with this name already exists").Send(w)
			return
		}
		log.Printf("Database error: %v", err)
		NewInternalError().Send(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		NewBadRequestError("Invalid ID format").Send(w)
		return
	}

	var product Product
	err = db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err == sql.ErrNoRows {
		NewNotFoundError("Product").Send(w)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		NewInternalError().Send(w)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		NewBadRequestError("Invalid ID format").Send(w)
		return
	}

	var input CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewBadRequestError("Invalid request body").Send(w)
		return
	}

	if input.Name == "" {
		NewBadRequestError("Name is required").Send(w)
		return
	}

	var product Product
	err = db.QueryRow(
		"UPDATE products SET name=$1, price=$2, category=$3 WHERE id=$4 RETURNING id, name, price, category",
		input.Name, input.Price, input.Category, id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err == sql.ErrNoRows {
		NewNotFoundError("Product").Send(w)
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			NewConflictError("A product with this name already exists").Send(w)
			return
		}
		log.Printf("Database error: %v", err)
		NewInternalError().Send(w)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		NewBadRequestError("Invalid ID format").Send(w)
		return
	}

	result, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Printf("Database error: %v", err)
		NewInternalError().Send(w)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		NewNotFoundError("Product").Send(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
