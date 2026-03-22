package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	db          *sql.DB
	minioClient *minio.Client
	bucketName  = "uploads"
)

type FileMetadata struct {
	ID          int       `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	ObjectKey   string    `json:"object_key"`
	UploadedAt  time.Time `json:"uploaded_at"`
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

	createTable()

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}
	minioClient, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	ctx := context.Background()
	exists, _ := minioClient.BucketExists(ctx, bucketName)
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatal("Failed to create bucket:", err)
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/upload", uploadFile)
	r.Get("/files", listFiles)
	r.Get("/files/{id}", getFileMetadata)
	r.Get("/files/{id}/download", downloadFile)
	r.Delete("/files/{id}", deleteFile)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size BIGINT NOT NULL,
		object_key TEXT NOT NULL,
		uploaded_at TIMESTAMP DEFAULT NOW()
	)`)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Max 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File is required. Use form field 'file'"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	objectKey := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		objectKey,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to upload file"})
		return
	}

	var meta FileMetadata
	err = db.QueryRow(
		"INSERT INTO files (filename, content_type, size, object_key) VALUES ($1, $2, $3, $4) RETURNING id, filename, content_type, size, object_key, uploaded_at",
		header.Filename, contentType, header.Size, objectKey,
	).Scan(&meta.ID, &meta.Filename, &meta.ContentType, &meta.Size, &meta.ObjectKey, &meta.UploadedAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save metadata"})
		return
	}

	writeJSON(w, http.StatusCreated, meta)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, filename, content_type, size, object_key, uploaded_at FROM files ORDER BY uploaded_at DESC")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	files := []FileMetadata{}
	for rows.Next() {
		var f FileMetadata
		rows.Scan(&f.ID, &f.Filename, &f.ContentType, &f.Size, &f.ObjectKey, &f.UploadedAt)
		files = append(files, f)
	}
	writeJSON(w, http.StatusOK, files)
}

func getFileMetadata(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var f FileMetadata
	err := db.QueryRow("SELECT id, filename, content_type, size, object_key, uploaded_at FROM files WHERE id = $1", id).
		Scan(&f.ID, &f.Filename, &f.ContentType, &f.Size, &f.ObjectKey, &f.UploadedAt)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "File not found"})
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var f FileMetadata
	err := db.QueryRow("SELECT id, filename, content_type, size, object_key, uploaded_at FROM files WHERE id = $1", id).
		Scan(&f.ID, &f.Filename, &f.ContentType, &f.Size, &f.ObjectKey, &f.UploadedAt)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "File not found"})
		return
	}

	object, err := minioClient.GetObject(context.Background(), bucketName, f.ObjectKey, minio.GetObjectOptions{})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve file"})
		return
	}
	defer object.Close()

	w.Header().Set("Content-Type", f.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", f.Filename))
	io.Copy(w, object)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var objectKey string
	err := db.QueryRow("SELECT object_key FROM files WHERE id = $1", id).Scan(&objectKey)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "File not found"})
		return
	}

	minioClient.RemoveObject(context.Background(), bucketName, objectKey, minio.RemoveObjectOptions{})
	db.Exec("DELETE FROM files WHERE id = $1", id)

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
