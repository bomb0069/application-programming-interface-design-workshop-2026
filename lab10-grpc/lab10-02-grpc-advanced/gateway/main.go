package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/example/lab18-grpc-advanced/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.ProductServiceClient

func main() {
	conn, err := grpc.NewClient("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client = pb.NewProductServiceClient(conn)

	http.HandleFunc("/api/products", productsHandler)
	http.HandleFunc("/api/products/", productHandler)

	log.Println("REST-to-gRPC Gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		stream, err := client.ListProducts(ctx, &pb.ListProductsRequest{
			Category: r.URL.Query().Get("category"),
		})
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		var products []*pb.Product
		for {
			p, err := stream.Recv()
			if err != nil {
				break
			}
			products = append(products, p)
		}
		if products == nil {
			products = []*pb.Product{}
		}
		json.NewEncoder(w).Encode(products)

	case http.MethodPost:
		var input struct {
			Name     string  `json:"name"`
			Price    float64 `json:"price"`
			Category string  `json:"category"`
		}
		json.NewDecoder(r.Body).Decode(&input)
		product, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
			Name: input.Name, Price: input.Price, Category: input.Category,
		})
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)

	default:
		writeError(w, 405, "Method not allowed")
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/products/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		writeError(w, 400, "Invalid ID")
		return
	}

	product, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: int32(id)})
	if err != nil {
		writeError(w, 404, "Product not found")
		return
	}
	json.NewEncoder(w).Encode(product)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
