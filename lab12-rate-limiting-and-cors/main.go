package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	createTable()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting middleware
	r.Use(RateLimitMiddleware(10, time.Minute))

	r.Get("/products", listProducts)
	r.Post("/products", createProduct)
	r.Get("/products/{id}", getProduct)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// TokenBucket implements a simple token bucket rate limiter
type TokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
}

func NewTokenBucket(maxTokens float64, refillPeriod time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: maxTokens / refillPeriod.Seconds(),
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() (bool, float64, time.Time) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	if tb.tokens < 1 {
		resetTime := now.Add(time.Duration((1 - tb.tokens) / tb.refillRate * float64(time.Second)))
		return false, 0, resetTime
	}

	tb.tokens--
	resetTime := now.Add(time.Duration(tb.maxTokens/tb.refillRate) * time.Second)
	return true, tb.tokens, resetTime
}

// RateLimitMiddleware creates a per-IP rate limiter
func RateLimitMiddleware(maxRequests float64, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	buckets := make(map[string]*TokenBucket)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			bucket, exists := buckets[ip]
			if !exists {
				bucket = NewTokenBucket(maxRequests, window)
				buckets[ip] = bucket
			}
			mu.Unlock()

			allowed, remaining, resetTime := bucket.Allow()

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(maxRequests)))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", strconv.FormatInt(resetTime.Unix()-time.Now().Unix(), 10))
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Rate limit exceeded. Try again later.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func createTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL
	)`)
	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO products (name, price, category) VALUES ('Laptop', 999.99, 'electronics'), ('Go Book', 39.99, 'books'), ('T-Shirt', 19.99, 'clothing')")
	}
}

func listProducts(w http.ResponseWriter, r *http.Request) {
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

func createProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Category string  `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	var p Product
	db.QueryRow("INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id, name, price, category",
		input.Name, input.Price, input.Category).Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var p Product
	err := db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
