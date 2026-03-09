# Lab 07 - Error Handling

## Learning Objectives

- Define a **consistent error response format** across all API endpoints
- Create **centralized error types** with factory functions
- Use **error handling middleware** (Logger, Recoverer) from chi
- **Separate error concerns from business logic** for cleaner handlers

## Error Response Format

All errors from this API follow a standard JSON structure:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Product not found"
  }
}
```

| Field     | Description                                         |
|-----------|-----------------------------------------------------|
| `code`    | Machine-readable error code (e.g., `BAD_REQUEST`)   |
| `message` | Human-readable description of what went wrong       |

The HTTP status code is set on the response header but not repeated in the body.

## Getting Started

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

To stop the services:

```bash
docker compose down
```

To stop and remove the database volume:

```bash
docker compose down -v
```

## Test Examples

### 400 Bad Request - Missing required field

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"price": 9.99, "category": "books"}' | jq
```

Response:

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "Name is required"
  }
}
```

### 400 Bad Request - Invalid ID format

```bash
curl -s http://localhost:8080/products/abc | jq
```

Response:

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid ID format"
  }
}
```

### 400 Bad Request - Invalid price

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Widget", "price": -5, "category": "tools"}' | jq
```

Response:

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "Price must be greater than 0"
  }
}
```

### 404 Not Found - Product does not exist

```bash
curl -s http://localhost:8080/products/9999 | jq
```

Response:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Product not found"
  }
}
```

### 409 Conflict - Duplicate product name

First, create a product:

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Gadget", "price": 29.99, "category": "electronics"}' | jq
```

Then try to create another with the same name:

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Gadget", "price": 19.99, "category": "electronics"}' | jq
```

Response:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "A product with this name already exists"
  }
}
```

### 201 Created - Successful creation

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Book", "price": 12.99, "category": "books"}' | jq
```

Response:

```json
{
  "id": 1,
  "name": "Book",
  "price": 12.99,
  "category": "books"
}
```

## Code Walkthrough

### errors.go - Centralized Error Types

The `APIError` struct provides a consistent shape for all error responses:

```go
type APIError struct {
    StatusCode int    `json:"-"`       // HTTP status code (not included in JSON body)
    Code       string `json:"code"`    // Machine-readable error code
    Message    string `json:"message"` // Human-readable message
}
```

Factory functions create specific error types, keeping the details in one place:

- `NewBadRequestError(message)` -- 400 for validation failures and malformed input
- `NewNotFoundError(resource)` -- 404 when a resource does not exist
- `NewConflictError(message)` -- 409 for uniqueness constraint violations
- `NewInternalError()` -- 500 for unexpected server errors (message is generic on purpose)

The `Send(w)` method writes the status code and JSON body to the response:

```go
func (e *APIError) Send(w http.ResponseWriter) {
    w.WriteHeader(e.StatusCode)
    json.NewEncoder(w).Encode(ErrorResponse{Error: *e})
}
```

### Handlers Using Centralized Errors

Handlers call error factory functions and `Send()` instead of manually building responses. This pattern keeps handlers focused on business logic:

```go
func getProduct(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        NewBadRequestError("Invalid ID format").Send(w)
        return
    }

    // ... query database ...

    if err == sql.ErrNoRows {
        NewNotFoundError("Product").Send(w)
        return
    }
}
```

### Chi Middleware

Three middleware functions are applied to every request:

1. **`middleware.Logger`** -- Logs every request with method, path, status code, and duration.
2. **`middleware.Recoverer`** -- Catches panics in handlers and returns a 500 instead of crashing the server.
3. **`JSONContentType`** (custom) -- Sets `Content-Type: application/json` on all responses so individual handlers do not need to.

```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(JSONContentType)
```

## Exercises

### 1. Add a NewValidationError with Field Errors

Create a `NewValidationError` factory that accepts a list of field-level errors and returns them in a structured format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": [
      {"field": "name", "message": "Name is required"},
      {"field": "price", "message": "Price must be greater than 0"}
    ]
  }
}
```

Hints:
- Add a `Fields` slice to `APIError` (or create a new `ValidationError` struct)
- Validate all fields at once instead of returning on the first error
- Return 422 Unprocessable Entity

### 2. Add Request ID to Error Responses

Use `middleware.RequestID` to attach a unique request ID to every response, and include it in error responses:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Product not found",
    "request_id": "abc-123-def"
  }
}
```

Hints:
- Add `r.Use(middleware.RequestID)` to the middleware chain
- Use `middleware.GetReqID(r.Context())` in the `Send` method (you will need to pass the request or context)
- Set the `X-Request-Id` response header as well

### 3. Custom JSON Recovery Middleware

The default `middleware.Recoverer` returns plain text on panic. Write a custom recovery middleware that returns a JSON error response instead:

```go
func JSONRecoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rvr := recover(); rvr != nil {
                // Log the panic, return JSON 500 error
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

Test it by adding a temporary handler that panics:

```go
r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
    panic("something went wrong")
})
```

### 4. Error Logging Middleware

Create middleware that wraps `http.ResponseWriter` to capture the status code, then logs details for any response with status >= 400:

```
ERROR [POST /products] 409 - 2.3ms - Body: {"name": "Gadget", ...}
```

Hints:
- Create a `statusResponseWriter` that wraps `http.ResponseWriter` and records the status code
- Read and restore `r.Body` using `io.ReadAll` and `io.NopCloser` so the body is still available to the handler
- Log method, path, status, duration, and request body for error responses

## Key Concepts

### Consistent Error Format

Every error response uses the same JSON structure. Clients can rely on parsing `error.code` for programmatic handling and `error.message` for display. This eliminates guesswork about what shape an error response will take.

### Error Factory Functions

Factory functions like `NewNotFoundError("Product")` centralize error creation. If you need to change the format, status code, or add fields, you change it in one place. They also make handlers more readable -- the intent is immediately clear.

### Middleware Chain

Middleware runs in order for every request. The chi middleware chain in this lab:

```
Request -> Logger -> Recoverer -> JSONContentType -> Handler -> Response
```

- **Logger** wraps the request to log timing and status after the handler runs
- **Recoverer** catches any panics so the server stays up
- **JSONContentType** sets the content type header before the handler runs

### Separation of Concerns

Error formatting is in `errors.go`. Business logic and routing are in `handlers.go` and `main.go`. Handlers do not need to know how errors are serialized -- they just call `Send(w)`. This makes it straightforward to change error formatting, add fields, or switch serialization formats without touching handler code.

## Cleanup

```bash
docker compose down -v
```
