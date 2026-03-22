package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
)

// V2Product is the enriched product representation returned by API v2,
// including description and tags fields.
type V2Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// V2Response wraps v2 payloads with a version indicator.
type V2Response struct {
	Data    interface{} `json:"data"`
	Version string      `json:"version"`
}

func v2ListProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, category, description, tags FROM products ORDER BY id")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	products := []V2Product{}
	for rows.Next() {
		var p V2Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, pq.Array(&p.Tags))
		products = append(products, p)
	}

	writeJSON(w, http.StatusOK, V2Response{Data: products, Version: "2.0"})
}

func v2GetProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var p V2Product
	err := db.QueryRow("SELECT id, name, price, category, description, tags FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, pq.Array(&p.Tags))
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}

	writeJSON(w, http.StatusOK, V2Response{Data: p, Version: "2.0"})
}

func v2CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string   `json:"name"`
		Price       float64  `json:"price"`
		Category    string   `json:"category"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
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

	var p V2Product
	db.QueryRow(
		"INSERT INTO products (name, price, category, description, tags) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, price, category, description, tags",
		input.Name, input.Price, input.Category, input.Description, pq.Array(input.Tags),
	).Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, pq.Array(&p.Tags))

	writeJSON(w, http.StatusCreated, V2Response{Data: p, Version: "2.0"})
}
