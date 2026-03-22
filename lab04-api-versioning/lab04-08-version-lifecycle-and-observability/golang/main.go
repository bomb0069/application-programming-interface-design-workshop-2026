package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type contextKey string

const versionKey contextKey = "api-version"

var db *sql.DB

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"version", "endpoint", "method", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"version", "endpoint", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// Version lifecycle definitions
type VersionInfo struct {
	Version     string `json:"version"`
	Stage       string `json:"stage"`
	ReleasedAt  string `json:"released_at"`
	Description string `json:"description"`
}

var versionLifecycle = []VersionInfo{
	{Version: "v1", Stage: "deprecated", ReleasedAt: "2025-01-01", Description: "Original API. Deprecated since 2026-03-01. Sunset: 2026-09-01."},
	{Version: "v2", Stage: "current", ReleasedAt: "2026-01-01", Description: "Current version. Enhanced responses with envelope, description, tags."},
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

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// V1 — deprecated
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(setVersion("v1"))
		r.Use(prometheusMiddleware)
		r.Use(structuredLogMiddleware)
		r.Get("/products", v1ListProducts)
		r.Get("/products/{id}", v1GetProduct)
		r.Post("/products", v1CreateProduct)
	})

	// V2 — current
	r.Route("/api/v2", func(r chi.Router) {
		r.Use(setVersion("v2"))
		r.Use(prometheusMiddleware)
		r.Use(structuredLogMiddleware)
		r.Get("/products", v2ListProducts)
		r.Get("/products/{id}", v2GetProduct)
		r.Post("/products", v2CreateProduct)
	})

	// Version lifecycle endpoint
	r.Get("/api/lifecycle", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"versions": versionLifecycle,
			"policy": map[string]string{
				"minimum_support_window": "12 months from release",
				"deprecation_notice":     "6 months before sunset",
				"sunset_trigger":         "Traffic below 1% OR support window expires",
			},
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func setVersion(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), versionKey, version)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getVersion(r *http.Request) string {
	if v, ok := r.Context().Value(versionKey).(string); ok {
		return v
	}
	return "v1"
}

// prometheusMiddleware records request metrics with version label
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &statusRecorder{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(wrapped, r)

		version := getVersion(r)
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(version, r.URL.Path, r.Method, strconv.Itoa(wrapped.statusCode)).Inc()
		httpRequestDuration.WithLabelValues(version, r.URL.Path, r.Method).Observe(duration)
	})
}

// structuredLogMiddleware emits structured JSON log entries
func structuredLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &statusRecorder{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(wrapped, r)

		version := getVersion(r)
		logEntry := map[string]interface{}{
			"ts":          time.Now().UTC().Format(time.RFC3339),
			"api_version": version,
			"method":      r.Method,
			"endpoint":    r.URL.Path,
			"status":      wrapped.statusCode,
			"latency_ms":  time.Since(start).Milliseconds(),
			"user_agent":  r.UserAgent(),
		}
		logJSON, _ := json.Marshal(logEntry)
		log.Println(string(logJSON))
	})
}

// statusRecorder wraps ResponseWriter to capture status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func createTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category TEXT NOT NULL,
		description TEXT DEFAULT '',
		tags TEXT[] DEFAULT '{}'
	)`)

	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO products (name, price, category, description, tags) VALUES
			('Laptop', 999.99, 'electronics', 'A powerful laptop for developers', '{portable,computing}'),
			('Go Book', 39.99, 'books', 'Learn Go programming', '{programming,education}'),
			('T-Shirt', 19.99, 'clothing', 'Comfortable cotton t-shirt', '{casual,cotton}')`)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
