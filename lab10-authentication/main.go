package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var jwtSecret []byte

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)

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

	createTables()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public routes
	r.Post("/register", registerHandler)
	r.Post("/login", loginHandler)
	r.Get("/health", healthHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/products", listProducts)
		r.Post("/products", createProduct)
		r.Get("/products/{id}", getProduct)
		r.Get("/me", meHandler)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
