# Lab 17 - gRPC Basics

## Learning Objectives

- Define services and messages with Protocol Buffers
- Implement a gRPC server and client in Go
- Understand protobuf serialization and code generation
- Use gRPCUI for interactive testing
- Compare gRPC with REST

## Architecture

```
┌──────────────────┐     gRPC (HTTP/2)     ┌──────────────────┐
│   gRPC Client    │ ───────────────────>   │   gRPC Server    │
│   (Go CLI)       │     protobuf binary    │   :50051         │
└──────────────────┘                        └──────────────────┘
                                                    ^
┌──────────────────┐     gRPC reflection            │
│   gRPCUI         │ ──────────────────────────────-┘
│   :8080          │
└──────────────────┘
```

- **Server** (:50051) - gRPC server implementing the ProductService
- **Client** (CLI) - Go client that calls all five RPCs and prints results
- **gRPCUI** (:8080) - Web UI for interactively calling gRPC methods

## Prerequisites

- Docker and Docker Compose

## Getting Started

Start all services:

```bash
docker-compose up --build
```

The client runs once, calls all RPCs, and exits. View its output:

```bash
docker-compose logs client
```

Open gRPCUI at [http://localhost:8080](http://localhost:8080) to interactively call methods.

## Understanding the Proto File

The file `proto/product.proto` defines the entire API contract:

```protobuf
syntax = "proto3";           // Use proto3 syntax
package product;             // Protobuf package namespace
option go_package = "...";   // Go import path for generated code
```

### Service Definition

The `service` block defines the RPC methods, similar to an interface:

```protobuf
service ProductService {
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc GetProduct(GetProductRequest) returns (Product);
  rpc CreateProduct(CreateProductRequest) returns (Product);
  rpc UpdateProduct(UpdateProductRequest) returns (Product);
  rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse);
}
```

Each RPC takes exactly one request message and returns exactly one response message (unary RPCs).

### Message Types

Messages define the data structures. Each field has a type, name, and unique field number:

```protobuf
message Product {
  int32 id = 1;        // Field number 1
  string name = 2;     // Field number 2
  double price = 3;    // Field number 3
  string category = 4; // Field number 4
}
```

Field numbers are used in the binary encoding -- they must never be changed once the schema is in use.

### Common Proto Types

| Proto Type | Go Type   | Description           |
|------------|-----------|-----------------------|
| `int32`    | `int32`   | Variable-length int   |
| `int64`    | `int64`   | Variable-length int   |
| `double`   | `float64` | 64-bit floating point |
| `string`   | `string`  | UTF-8 string          |
| `bool`     | `bool`    | Boolean               |
| `repeated` | `[]T`     | List/slice of values  |

## Code Walkthrough

### Proto Code Generation

The Dockerfiles run `protoc` to generate Go code from the `.proto` file:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/product.proto
```

This produces two files:
- `proto/product.pb.go` - Message types (Product, ListProductsRequest, etc.)
- `proto/product_grpc.pb.go` - Service interface and client/server stubs

### Server Implementation

The server embeds `UnimplementedProductServiceServer` to satisfy the interface. This allows forward compatibility -- new RPCs added to the proto return "unimplemented" by default:

```go
type productServer struct {
    pb.UnimplementedProductServiceServer
    mu       sync.RWMutex
    products map[int32]*pb.Product
    nextID   int32
}
```

Each method implements the corresponding RPC:

```go
func (s *productServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    p, ok := s.products[req.Id]
    if !ok {
        return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
    }
    return p, nil
}
```

Register the server and enable reflection (required for gRPCUI):

```go
s := grpc.NewServer()
pb.RegisterProductServiceServer(s, newServer())
reflection.Register(s)
```

### Client Implementation

The client creates a connection and a typed client stub:

```go
conn, err := grpc.NewClient("server:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
client := pb.NewProductServiceClient(conn)
```

Then calls methods with full type safety:

```go
product, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
```

### gRPC Status Codes

gRPC uses its own status codes instead of HTTP status codes:

| gRPC Code          | HTTP Equivalent | Meaning                    |
|--------------------|-----------------|----------------------------|
| `OK`               | 200             | Success                    |
| `NotFound`         | 404             | Resource not found         |
| `InvalidArgument`  | 400             | Bad request                |
| `AlreadyExists`    | 409             | Resource already exists    |
| `Unauthenticated`  | 401             | Missing/invalid auth       |
| `PermissionDenied` | 403             | Insufficient permissions   |
| `Internal`         | 500             | Server error               |
| `Unimplemented`    | 501             | RPC not implemented        |
| `Unavailable`      | 503             | Service unavailable        |

Return errors with:

```go
return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
```

## gRPC vs REST Comparison

| Aspect            | gRPC                          | REST                         |
|-------------------|-------------------------------|------------------------------|
| Protocol          | HTTP/2                        | HTTP/1.1 or HTTP/2           |
| Data Format       | Protocol Buffers (binary)     | JSON (text)                  |
| API Contract      | `.proto` file (strict)        | OpenAPI/informal (flexible)  |
| Code Generation   | Built-in (protoc)             | Optional (openapi-generator) |
| Streaming         | Native (4 patterns)           | SSE, WebSocket (separate)    |
| Browser Support   | Requires gRPC-Web proxy       | Native                       |
| Type Safety       | Strong (generated types)      | Weak (JSON parsing)          |
| Performance       | Faster (binary, HTTP/2)       | Slower (text, HTTP/1.1)      |
| Tooling           | grpcurl, gRPCUI, Buf          | curl, Postman, Swagger       |
| Human Readable    | No (binary wire format)       | Yes (JSON)                   |
| Load Balancing    | Requires L7 / client-side     | Standard L4/L7               |

**When to use gRPC:** microservice-to-microservice communication, high-performance requirements, streaming, polyglot environments.

**When to use REST:** public APIs, browser clients, simple CRUD, broad ecosystem compatibility.

## Exercises

### Exercise 1: Add SearchProducts RPC

Add a new RPC that searches products by a query string:

```protobuf
message SearchProductsRequest {
  string query = 1;
}

rpc SearchProducts(SearchProductsRequest) returns (ListProductsResponse);
```

Implement it in the server to search by name or category (case-insensitive substring match). Update the client to test it.

### Exercise 2: Add Field Validation

Improve input validation in `CreateProduct`:

- Name must be non-empty and under 100 characters
- Price must be positive
- Category must be one of: "electronics", "books", "clothing", "food"

Return `codes.InvalidArgument` with descriptive error messages for each violation.

### Exercise 3: Add Metadata (Headers)

gRPC supports metadata, similar to HTTP headers. Add a request ID:

**Server side** - read metadata:
```go
import "google.golang.org/grpc/metadata"

md, ok := metadata.FromIncomingContext(ctx)
if ok {
    if vals := md.Get("x-request-id"); len(vals) > 0 {
        log.Printf("Request ID: %s", vals[0])
    }
}
```

**Client side** - send metadata:
```go
ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", "req-123")
```

### Exercise 4: Add Interceptors

Interceptors are gRPC's middleware. Add a logging interceptor:

```go
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("Method: %s | Duration: %s | Error: %v",
        info.FullMethod, time.Since(start), err)
    return resp, err
}

s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
```

Add this to the server and observe the logs when the client makes calls.

## Key Concepts

| Concept                  | Description                                                       |
|--------------------------|-------------------------------------------------------------------|
| Protocol Buffers         | Language-neutral serialization format; defines messages and types  |
| gRPC Service Definition  | `.proto` service block defining RPCs with request/response types   |
| Unary RPC                | Single request, single response (like a function call)            |
| Code Generation          | `protoc` generates typed client stubs and server interfaces       |
| gRPC Status Codes        | Structured error codes (NotFound, InvalidArgument, etc.)          |
| Reflection               | Server exposes its schema at runtime for tools like gRPCUI        |
| Interceptors             | Middleware pattern for cross-cutting concerns (logging, auth)     |
| Metadata                 | Key-value pairs sent alongside RPCs (like HTTP headers)           |

## Cleanup

```bash
docker-compose down
```
