# Lab 03 - Path Parameters

## Learning Objectives

- Use URL path parameters to identify specific resources
- Introduce the chi router for flexible routing in Go
- Return appropriate HTTP status codes (200, 400, 404)
- Parse and validate path parameters

## Prerequisites

- Go 1.24 or later installed
- Docker and Docker Compose installed
- Completion of Lab 02 or equivalent understanding of basic HTTP handlers and JSON responses
- A terminal and a tool for making HTTP requests (curl, Postman, or a browser)

## Getting Started

### Option A: Run with Go

```bash
go mod tidy
go run main.go
```

### Option B: Run with Docker Compose

```bash
docker compose up --build
```

The server will start on http://localhost:8080.

## Test Commands

### List all items

```bash
curl http://localhost:8080/items
```

Expected response (HTTP 200):

```json
[
  {"id": 1, "name": "Laptop", "price": 999.99},
  {"id": 2, "name": "Mouse", "price": 29.99},
  {"id": 3, "name": "Keyboard", "price": 79.99},
  {"id": 4, "name": "Monitor", "price": 549.99}
]
```

### Get a specific item by ID

```bash
curl http://localhost:8080/items/1
```

Expected response (HTTP 200):

```json
{"id": 1, "name": "Laptop", "price": 999.99}
```

### Request an item that does not exist (404)

```bash
curl http://localhost:8080/items/999
```

Expected response (HTTP 404):

```json
{"error": "Item not found"}
```

### Request with an invalid ID (400)

```bash
curl http://localhost:8080/items/abc
```

Expected response (HTTP 400):

```json
{"error": "Invalid ID format"}
```

## Code Walkthrough

### The chi Router

The standard library `net/http` package provides basic routing, but it does not support path parameters out of the box. The **chi** router is a lightweight, composable router for Go that adds support for URL path parameters, middleware, and more.

```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Get("/items/{id}", getItem)
```

The `{id}` syntax in the route pattern defines a named path parameter. Chi will match any value in that position and make it available to your handler.

### Extracting Path Parameters with URLParam

Inside your handler, use `chi.URLParam` to retrieve the value of a path parameter by name:

```go
idStr := chi.URLParam(r, "id")
```

This returns the raw string from the URL. Since path parameters are always strings, you need to parse them into the appropriate type. In this lab, we convert the ID to an integer using `strconv.Atoi`:

```go
id, err := strconv.Atoi(idStr)
if err != nil {
    // handle invalid input
}
```

### HTTP Status Codes

This lab demonstrates three important status codes:

- **200 OK** - The request succeeded. This is the default status code when you write a response without explicitly setting one. Used when returning the list of items or a single item.
- **400 Bad Request** - The client sent a request that the server cannot process. Used when the path parameter is not a valid integer (e.g., `/items/abc`).
- **404 Not Found** - The requested resource does not exist. Used when the ID is valid but no item matches (e.g., `/items/999`).

Set the status code before writing the response body:

```go
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
```

**Important:** You must call `w.WriteHeader()` before writing the body. If you call `w.WriteHeader()` after writing to the response, it will have no effect because the default 200 status is sent with the first write.

## Exercises

### Exercise 1: Add a Price Endpoint

Add a `GET /items/{id}/price` endpoint that returns only the price of an item.

Expected response for `curl http://localhost:8080/items/1/price`:

```json
{"price": 999.99}
```

Handle the same error cases (400 for invalid ID, 404 for missing item).

### Exercise 2: Filter by Category

Add a `Category` field to the `Item` struct and assign categories to each item (e.g., "electronics", "peripherals"). Then create a `GET /categories/{category}/items` endpoint that returns all items matching the given category.

Expected response for `curl http://localhost:8080/categories/peripherals/items`:

```json
[
  {"id": 2, "name": "Mouse", "price": 29.99, "category": "peripherals"},
  {"id": 3, "name": "Keyboard", "price": 79.99, "category": "peripherals"}
]
```

Return an empty array `[]` if no items match the category.

### Exercise 3: Multiple Path Parameters

Add a `StoreID` field to the `Item` struct and create a `GET /stores/{storeId}/items/{itemId}` endpoint that looks up an item by both store and item ID.

This exercise demonstrates how to extract and use multiple path parameters in a single handler:

```go
storeID := chi.URLParam(r, "storeId")
itemID := chi.URLParam(r, "itemId")
```

Return 404 if the store or item is not found, and 400 if either parameter is not a valid integer.

## Key Concepts

### Path Parameters

Path parameters are variable segments in a URL path, identified by a placeholder name inside curly braces (e.g., `{id}`). They allow a single route to handle requests for many different resources. Unlike query parameters (which appear after `?`), path parameters are part of the URL path itself and typically identify a specific resource.

| Feature         | Path Parameter          | Query Parameter            |
| --------------- | ----------------------- | -------------------------- |
| Syntax          | `/items/{id}`           | `/items?id=1`              |
| Purpose         | Identify a resource     | Filter or modify a request |
| Required?       | Yes (part of the route) | Usually optional           |
| Example         | `/items/1`              | `/items?sort=name`         |

### Router Libraries

Go's standard `net/http` package handles basic routing, but third-party routers like **chi** provide features that simplify API development:

- Named path parameters (`{id}`, `{category}`)
- Route grouping and sub-routers
- Middleware support
- Method-based routing (`r.Get`, `r.Post`, `r.Put`, `r.Delete`)

Chi is a popular choice because it is compatible with the standard `net/http` interfaces, meaning your handlers have the same signature (`func(http.ResponseWriter, *http.Request)`) whether you use chi or the standard library.

### HTTP Status Codes

Status codes communicate the result of a request to the client:

| Code | Name        | When to Use                                     |
| ---- | ----------- | ----------------------------------------------- |
| 200  | OK          | The request succeeded and a response is returned |
| 400  | Bad Request | The client sent invalid data (e.g., non-numeric ID) |
| 404  | Not Found   | The requested resource does not exist            |

Always return meaningful status codes so that API clients can handle errors programmatically, rather than relying on parsing error messages.

## Cleanup

Stop the running server with `Ctrl+C`, or if using Docker Compose:

```bash
docker compose down
```
