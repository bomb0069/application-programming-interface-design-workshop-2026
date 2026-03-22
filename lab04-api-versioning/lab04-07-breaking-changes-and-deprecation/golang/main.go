package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Deprecation schedule
var deprecationSchedule = map[string]struct {
	Deprecation string
	Sunset      time.Time
}{
	"v1": {
		Deprecation: "Sat, 01 Mar 2026 00:00:00 GMT",
		Sunset:      time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
	},
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
	r.Use(middleware.Logger)

	// V1 — deprecated, with sunset filter
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(sunsetFilter("v1"))
		r.Use(deprecationMiddleware("v1"))
		r.Get("/products", v1ListProducts)
		r.Get("/products/{id}", v1GetProduct)
		r.Post("/products", v1CreateProduct)
	})

	// V2 — current version
	r.Route("/api/v2", func(r chi.Router) {
		r.Get("/products", v2ListProducts)
		r.Get("/products/{id}", v2GetProduct)
		r.Post("/products", v2CreateProduct)
	})

	// Breaking changes classification endpoint
	r.Get("/api/changes", breakingChangesHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// deprecationMiddleware adds Deprecation, Sunset, and Link headers
func deprecationMiddleware(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if schedule, ok := deprecationSchedule[version]; ok {
				w.Header().Set("Deprecation", schedule.Deprecation)
				w.Header().Set("Sunset", schedule.Sunset.Format(http.TimeFormat))
				w.Header().Set("Link", `</docs/migrate-v1-v2>; rel="successor-version"`)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// sunsetFilter returns 410 Gone after the sunset date
func sunsetFilter(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if schedule, ok := deprecationSchedule[version]; ok {
				if time.Now().UTC().After(schedule.Sunset) {
					writeJSON(w, http.StatusGone, map[string]string{
						"error":      "VERSION_SUNSET",
						"message":    "API " + version + " was sunset on " + schedule.Sunset.Format("2006-01-02"),
						"migrateUrl": "/docs/migrate-v1-v2",
					})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ChangeClassification represents a proposed API change and its classification
type ChangeClassification struct {
	Change         string `json:"change"`
	Classification string `json:"classification"`
	Explanation    string `json:"explanation"`
}

func breakingChangesHandler(w http.ResponseWriter, r *http.Request) {
	changes := []ChangeClassification{
		{Change: "Add optional query parameter 'sort'", Classification: "safe", Explanation: "New optional params don't affect existing clients"},
		{Change: "Add 'created_at' field to response", Classification: "safe", Explanation: "Additive response fields are backward compatible"},
		{Change: "Remove 'legacy_id' field from response", Classification: "breaking", Explanation: "Clients may depend on this field"},
		{Change: "Rename 'name' to 'title' in response", Classification: "breaking", Explanation: "Field rename breaks all existing clients"},
		{Change: "Change 'price' from number to string", Classification: "breaking", Explanation: "Type change breaks deserialization"},
		{Change: "Add new endpoint POST /api/v1/orders", Classification: "safe", Explanation: "New endpoints don't affect existing ones"},
		{Change: "Add new enum value 'premium' to 'tier'", Classification: "context-dependent", Explanation: "Safe if clients handle unknown values; breaking if they use exhaustive switch"},
		{Change: "Change error response from string to object", Classification: "context-dependent", Explanation: "Safe if clients only check HTTP status; breaking if they parse error body"},
		{Change: "Make 'email' field required (was optional)", Classification: "breaking", Explanation: "Existing requests without email will now fail"},
		{Change: "Change pagination from offset to cursor", Classification: "breaking", Explanation: "Completely changes how clients navigate results"},
	}
	writeJSON(w, http.StatusOK, changes)
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
