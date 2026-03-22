package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	pb "github.com/example/lab18-grpc-advanced/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedProductServiceServer
	mu       sync.RWMutex
	products map[int32]*pb.Product
	nextID   int32
}

func newServer() *server {
	s := &server{products: make(map[int32]*pb.Product), nextID: 4}
	s.products[1] = &pb.Product{Id: 1, Name: "Laptop", Price: 999.99, Category: "electronics"}
	s.products[2] = &pb.Product{Id: 2, Name: "Go Book", Price: 39.99, Category: "books"}
	s.products[3] = &pb.Product{Id: 3, Name: "T-Shirt", Price: 19.99, Category: "clothing"}
	return s
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	return p, nil
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := &pb.Product{Id: s.nextID, Name: req.Name, Price: req.Price, Category: req.Category}
	s.products[p.Id] = p
	s.nextID++
	return p, nil
}

// Server-side streaming: sends products one by one
func (s *server) ListProducts(req *pb.ListProductsRequest, stream pb.ProductService_ListProductsServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.products {
		if req.Category != "" && p.Category != req.Category {
			continue
		}
		time.Sleep(500 * time.Millisecond) // Simulate delay
		if err := stream.Send(p); err != nil {
			return err
		}
		log.Printf("Streamed product: %s", p.Name)
	}
	return nil
}

// Client-side streaming: receives multiple products from client
func (s *server) BatchCreateProducts(stream pb.ProductService_BatchCreateProductsServer) error {
	var created []*pb.Product
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.BatchCreateResponse{
				Count:    int32(len(created)),
				Products: created,
			})
		}
		if err != nil {
			return err
		}

		s.mu.Lock()
		p := &pb.Product{Id: s.nextID, Name: req.Name, Price: req.Price, Category: req.Category}
		s.products[p.Id] = p
		s.nextID++
		s.mu.Unlock()

		created = append(created, p)
		log.Printf("Batch created: %s", p.Name)
	}
}

// Bidirectional streaming: receive search queries, send matching products
func (s *server) ProductChat(stream pb.ProductService_ProductChatServer) error {
	for {
		query, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.mu.RLock()
		for _, p := range s.products {
			if strings.Contains(strings.ToLower(p.Name), strings.ToLower(query.Search)) ||
				strings.Contains(strings.ToLower(p.Category), strings.ToLower(query.Search)) {
				if err := stream.Send(p); err != nil {
					s.mu.RUnlock()
					return err
				}
			}
		}
		s.mu.RUnlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, newServer())
	reflection.Register(s)
	fmt.Println("gRPC server with streaming on :50051")
	log.Fatal(s.Serve(lis))
}
