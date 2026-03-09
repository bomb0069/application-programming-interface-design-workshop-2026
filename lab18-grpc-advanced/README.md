# Lab 18 - gRPC Advanced: Streaming and REST Gateway

## Learning Objectives

- Implement all four gRPC communication patterns (unary, server streaming, client streaming, bidirectional streaming)
- Build a REST-to-gRPC gateway that exposes gRPC services as REST endpoints
- Understand streaming use cases and lifecycle management

## Architecture

```
                    +-----------+
                    |  gRPCUI   |
                    |  :8081    |
                    +-----+-----+
                          |
+--------+          +-----v-----+          +---------+
| Client  +--------->  gRPC     <----------+ Gateway |
| (CLI)   |  gRPC   |  Server   |   gRPC   | :8080   |
+--------+          |  :50051   |          +----^----+
                    +-----------+               |
                                           REST | HTTP
                                           curl/browser
```

| Service   | Port  | Protocol | Description                      |
|-----------|-------|----------|----------------------------------|
| server    | 50051 | gRPC     | Product service with streaming   |
| client    | -     | gRPC     | CLI client demonstrating streams |
| gateway   | 8080  | HTTP     | REST-to-gRPC proxy               |
| grpcui    | 8081  | HTTP     | Web UI for gRPC testing          |

## Getting Started

```bash
docker-compose up --build
```

## Streaming Patterns Explained

### 1. Unary RPC (single request, single response)

Standard request-response, like a regular function call.

```
Client --[Request]--> Server
Client <--[Response]-- Server
```

**RPCs:** `GetProduct`, `CreateProduct`

### 2. Server-Side Streaming (single request, multiple responses)

Client sends one request, server streams back multiple responses. Useful for listing data, real-time feeds, or large result sets.

```
Client --[Request]-----> Server
Client <--[Response 1]-- Server
Client <--[Response 2]-- Server
Client <--[Response N]-- Server
Client <--[EOF]--------- Server
```

**RPC:** `ListProducts` - client requests products (optionally filtered by category), server streams them one by one.

### 3. Client-Side Streaming (multiple requests, single response)

Client sends multiple messages, server processes them and returns a single response. Useful for batch uploads, aggregation, or file uploads.

```
Client --[Request 1]--> Server
Client --[Request 2]--> Server
Client --[Request N]--> Server
Client --[EOF]--------> Server
Client <--[Response]--- Server
```

**RPC:** `BatchCreateProducts` - client streams multiple product creation requests, server responds with the total count and created products.

### 4. Bidirectional Streaming (multiple requests AND responses)

Both client and server send streams of messages independently. Useful for chat, real-time collaboration, or interactive queries.

```
Client --[Request 1]---> Server
Client <--[Response 1]-- Server
Client --[Request 2]---> Server
Client <--[Response 2]-- Server
Client <--[Response 3]-- Server
Client --[EOF]---------> Server
Client <--[EOF]--------- Server
```

**RPC:** `ProductChat` - client sends search queries, server responds with matching products for each query.

## Test the Gateway

The REST gateway translates HTTP requests into gRPC calls.

### List all products

```bash
curl http://localhost:8080/api/products
```

### Filter by category

```bash
curl http://localhost:8080/api/products?category=electronics
```

### Get a single product

```bash
curl http://localhost:8080/api/products/1
```

### Create a product

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Headphones", "price": 149.99, "category": "electronics"}'
```

## View Client Output

The client container runs all streaming demos automatically:

```bash
docker-compose logs client
```

Expected output:
```
=== Server-Side Streaming: ListProducts ===
  Received: [1] Laptop - $999.99
  Received: [2] Go Book - $39.99
  Received: [3] T-Shirt - $19.99

=== Client-Side Streaming: BatchCreateProducts ===
  Sending: Mouse
  Sending: Keyboard
  Sending: Monitor
  Batch created 3 products

=== Bidirectional Streaming: ProductChat ===
  Searching: electronics
  Searching: book
  Searching: shirt
  Found: [1] Laptop - $999.99 (electronics)
  Found: [2] Go Book - $39.99 (books)
  Found: [3] T-Shirt - $19.99 (clothing)

Done!
```

## Test with gRPCUI

Open [http://localhost:8081](http://localhost:8081) in your browser. gRPCUI provides a web interface to call any gRPC method, including streaming RPCs.

## Code Walkthrough

### Server-Side Streaming (server)

```go
func (s *server) ListProducts(req *pb.ListProductsRequest, stream pb.ProductService_ListProductsServer) error {
    for _, p := range s.products {
        stream.Send(p)  // Send each product individually
    }
    return nil  // Returning nil closes the stream
}
```

### Server-Side Streaming (client)

```go
stream, _ := client.ListProducts(ctx, &pb.ListProductsRequest{})
for {
    product, err := stream.Recv()
    if err == io.EOF { break }  // Stream closed
    // Process product...
}
```

### Client-Side Streaming (server)

```go
func (s *server) BatchCreateProducts(stream pb.ProductService_BatchCreateProductsServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&response)  // Send final response
        }
        // Process each request...
    }
}
```

### Client-Side Streaming (client)

```go
stream, _ := client.BatchCreateProducts(ctx)
for _, p := range products {
    stream.Send(&p)           // Send each product
}
resp, _ := stream.CloseAndRecv()  // Close stream and get response
```

### Bidirectional Streaming (server)

```go
func (s *server) ProductChat(stream pb.ProductService_ProductChatServer) error {
    for {
        query, err := stream.Recv()  // Receive query
        if err == io.EOF { return nil }
        // Send matching products back
        stream.Send(matchingProduct)
    }
}
```

### Bidirectional Streaming (client)

```go
stream, _ := client.ProductChat(ctx)
// Send queries
stream.Send(&pb.ProductQuery{Search: "electronics"})
stream.CloseSend()  // Done sending
// Receive results
for {
    product, err := stream.Recv()
    if err == io.EOF { break }
}
```

## Exercises

1. **Real-Time Price Updates** - Add a server-side streaming RPC `WatchPriceUpdates` that streams price changes in real-time. Simulate random price fluctuations every second.

2. **gRPC-to-HTTP Status Mapping** - Improve the gateway error handling to map gRPC status codes to appropriate HTTP status codes (e.g., `codes.NotFound` to `404`, `codes.InvalidArgument` to `400`).

3. **Logging Interceptor** - Add a gRPC unary and stream interceptor that logs the method name, duration, and status code for every RPC call.

4. **Deadline Propagation** - Add timeout/deadline support in the gateway so that REST requests with a `?timeout=5s` query parameter propagate deadlines to the gRPC server.

## Key Concepts

| Concept                  | Description                                                    |
|--------------------------|----------------------------------------------------------------|
| **Server Streaming**     | Server sends multiple messages; client reads until EOF         |
| **Client Streaming**     | Client sends multiple messages; server reads until EOF         |
| **Bidi Streaming**       | Both sides send/receive independently                          |
| **stream.Send(msg)**     | Send a message on the stream                                   |
| **stream.Recv()**        | Receive a message; returns `io.EOF` when stream ends           |
| **stream.CloseSend()**   | Client signals it is done sending                              |
| **SendAndClose(resp)**   | Server sends final response and closes client stream           |
| **CloseAndRecv()**       | Client closes send side and receives server's final response   |
| **REST-to-gRPC Gateway** | HTTP server that translates REST calls into gRPC calls         |
| **gRPC Reflection**      | Allows tools like gRPCUI to discover services at runtime       |

## Cleanup

```bash
docker-compose down
```
