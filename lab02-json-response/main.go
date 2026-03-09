package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var books = []Book{
	{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015},
	{ID: 2, Title: "Go in Action", Author: "William Kennedy", Year: 2015},
	{ID: 3, Title: "Learning Go", Author: "Jon Bodner", Year: 2021},
}

func main() {
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/count", booksCountHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func booksCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": len(books)})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
