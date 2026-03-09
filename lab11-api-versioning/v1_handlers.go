package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// V1 returns simple flat product objects
type V1Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func v1ListProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, category FROM products ORDER BY id")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	products := []V1Product{}
	for rows.Next() {
		var p V1Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category)
		products = append(products, p)
	}
	writeJSON(w, http.StatusOK, products)
}

func v1GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var p V1Product
	err = db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func v1CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Category string  `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	if input.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name is required"})
		return
	}

	var p V1Product
	db.QueryRow(
		"INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id, name, price, category",
		input.Name, input.Price, input.Category,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Category)

	writeJSON(w, http.StatusCreated, p)
}
