package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/example/lab17-grpc-basics/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("server:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List all products
	fmt.Println("=== List Products ===")
	listResp, err := client.ListProducts(ctx, &pb.ListProductsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range listResp.Products {
		fmt.Printf("  [%d] %s - $%.2f (%s)\n", p.Id, p.Name, p.Price, p.Category)
	}

	// Create a product
	fmt.Println("\n=== Create Product ===")
	created, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:     "Headphones",
		Price:    79.99,
		Category: "electronics",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Created: [%d] %s - $%.2f\n", created.Id, created.Name, created.Price)

	// Get single product
	fmt.Println("\n=== Get Product ===")
	product, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Got: [%d] %s - $%.2f (%s)\n", product.Id, product.Name, product.Price, product.Category)

	// Update product
	fmt.Println("\n=== Update Product ===")
	updated, err := client.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:    1,
		Price: 1099.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Updated: [%d] %s - $%.2f\n", updated.Id, updated.Name, updated.Price)

	// Delete product
	fmt.Println("\n=== Delete Product ===")
	delResp, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: 3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Deleted: %v\n", delResp.Success)

	// List again
	fmt.Println("\n=== List Products (after changes) ===")
	listResp, err = client.ListProducts(ctx, &pb.ListProductsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range listResp.Products {
		fmt.Printf("  [%d] %s - $%.2f (%s)\n", p.Id, p.Name, p.Price, p.Category)
	}
}
