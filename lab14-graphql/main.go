package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type Category struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
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

	createTableAndSeed()

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    buildQueryType(),
		Mutation: buildMutationType(),
	})
	if err != nil {
		log.Fatal("Failed to create GraphQL schema:", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// GraphQL endpoint
	r.Handle("/graphql", handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	}))

	// REST endpoint for comparison
	r.Get("/api/products", restListProducts)

	log.Println("Server starting on :8080")
	log.Println("GraphQL Playground: http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTableAndSeed() {
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL
	)`)

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 {
		return
	}

	products := []struct {
		name, category string
		price          float64
	}{
		{"Laptop Pro", "electronics", 1299.99},
		{"Wireless Mouse", "electronics", 29.99},
		{"Mechanical Keyboard", "electronics", 89.99},
		{"Go Programming Language", "books", 39.99},
		{"Clean Code", "books", 34.99},
		{"Design Patterns", "books", 44.99},
		{"Cotton T-Shirt", "clothing", 19.99},
		{"Denim Jeans", "clothing", 59.99},
		{"Running Shoes", "clothing", 89.99},
	}
	for _, p := range products {
		db.Exec("INSERT INTO products (name, price, category) VALUES ($1, $2, $3)", p.name, p.price, p.category)
	}
}

func restListProducts(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name, price, category FROM products ORDER BY id")
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
