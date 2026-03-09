package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/example/lab17-grpc-basics/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type productServer struct {
	pb.UnimplementedProductServiceServer
	mu       sync.RWMutex
	products map[int32]*pb.Product
	nextID   int32
}

func newServer() *productServer {
	s := &productServer{
		products: make(map[int32]*pb.Product),
		nextID:   4,
	}
	s.products[1] = &pb.Product{Id: 1, Name: "Laptop", Price: 999.99, Category: "electronics"}
	s.products[2] = &pb.Product{Id: 2, Name: "Go Book", Price: 39.99, Category: "books"}
	s.products[3] = &pb.Product{Id: 3, Name: "T-Shirt", Price: 19.99, Category: "clothing"}
	return s
}

func (s *productServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []*pb.Product
	for _, p := range s.products {
		products = append(products, p)
	}
	return &pb.ListProductsResponse{Products: products}, nil
}

func (s *productServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	return p, nil
}

func (s *productServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	p := &pb.Product{
		Id:       s.nextID,
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
	}
	s.products[p.Id] = p
	s.nextID++
	return p, nil
}

func (s *productServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.Category != "" {
		p.Category = req.Category
	}
	return p, nil
}

func (s *productServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[req.Id]; !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	delete(s.products, req.Id)
	return &pb.DeleteProductResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, newServer())
	reflection.Register(s)

	fmt.Println("gRPC server starting on :50051")
	log.Fatal(s.Serve(lis))
}
