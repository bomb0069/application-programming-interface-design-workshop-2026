package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/example/lab18-grpc-advanced/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewProductServiceClient(conn)
	ctx := context.Background()

	// 1. Server-side streaming
	fmt.Println("=== Server-Side Streaming: ListProducts ===")
	stream, err := client.ListProducts(ctx, &pb.ListProductsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Received: [%d] %s - $%.2f\n", product.Id, product.Name, product.Price)
	}

	// 2. Client-side streaming
	fmt.Println("\n=== Client-Side Streaming: BatchCreateProducts ===")
	batchStream, err := client.BatchCreateProducts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	newProducts := []pb.CreateProductRequest{
		{Name: "Mouse", Price: 29.99, Category: "electronics"},
		{Name: "Keyboard", Price: 79.99, Category: "electronics"},
		{Name: "Monitor", Price: 449.99, Category: "electronics"},
	}
	for _, p := range newProducts {
		fmt.Printf("  Sending: %s\n", p.Name)
		batchStream.Send(&p)
		time.Sleep(300 * time.Millisecond)
	}
	batchResp, err := batchStream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Batch created %d products\n", batchResp.Count)

	// 3. Bidirectional streaming
	fmt.Println("\n=== Bidirectional Streaming: ProductChat ===")
	chatStream, err := client.ProductChat(ctx)
	if err != nil {
		log.Fatal(err)
	}

	queries := []string{"electronics", "book", "shirt"}
	for _, q := range queries {
		fmt.Printf("  Searching: %s\n", q)
		chatStream.Send(&pb.ProductQuery{Search: q})

		// Small delay to receive responses
		time.Sleep(500 * time.Millisecond)
	}
	chatStream.CloseSend()

	for {
		product, err := chatStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Found: [%d] %s - $%.2f (%s)\n", product.Id, product.Name, product.Price, product.Category)
	}

	fmt.Println("\nDone!")
}
