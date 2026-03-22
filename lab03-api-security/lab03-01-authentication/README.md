# Lab 10 - Authentication

## Learning Objectives

- Implement JWT (JSON Web Token) authentication
- Password hashing with bcrypt
- Auth middleware for protecting routes
- Protected vs public routes
- Token-based API security

## Getting Started

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

## Test Workflow

### 1. Health Check (Public)

```bash
curl http://localhost:8080/health
```

### 2. Register a New User

```bash
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123"}'
```

Response:
```json
{"id":1,"username":"john","email":"john@example.com"}
```

### 3. Login to Get a Token

```bash
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}'
```

Response:
```json
{"token":"eyJhbGciOiJIUzI1NiIs...","expires_in":86400}
```

Save the token value for subsequent requests.

### 4. Access Protected Route with Token

```bash
curl -s http://localhost:8080/products \
  -H "Authorization: Bearer <token>"
```

Replace `<token>` with the token from the login response.

### 5. Access Protected Route Without Token (401)

```bash
curl -s http://localhost:8080/products
```

Response:
```json
{"error":"Authorization header required"}
```

### 6. Get User Profile

```bash
curl -s http://localhost:8080/me \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{"id":1,"username":"john","email":"john@example.com"}
```

### 7. Create a Product (Authenticated)

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Keyboard","price":79.99,"category":"electronics"}'
```

## Code Walkthrough

### JWT Claims

When a user logs in, the server creates a JWT containing claims:

```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id":  user.ID,
    "username": user.Username,
    "exp":      time.Now().Add(24 * time.Hour).Unix(),
})
```

The token is signed with a secret key (`JWT_SECRET`) and returned to the client. The token has three parts separated by dots: `header.payload.signature`.

### Bcrypt Password Hashing

Passwords are never stored in plain text. We use bcrypt to hash them:

```go
// Hashing during registration
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

// Comparing during login
err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
```

### Middleware Pattern

The `AuthMiddleware` intercepts requests to protected routes, validates the JWT, and injects the user into the request context:

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Extract token from Authorization header
        // 2. Parse and validate the JWT
        // 3. Extract user info from claims
        // 4. Add user to request context
        // 5. Call next handler
    })
}
```

Routes are grouped so the middleware only applies to protected endpoints:

```go
r.Group(func(r chi.Router) {
    r.Use(AuthMiddleware)
    r.Get("/products", listProducts)
    r.Get("/me", meHandler)
})
```

### Context Values

The authenticated user is passed through the request context:

```go
// Setting in middleware
ctx := context.WithValue(r.Context(), userContextKey, user)

// Retrieving in handler
user := r.Context().Value(userContextKey).(User)
```

## Exercises

1. **Role-Based Access Control** - Add a `role` field to the users table (e.g., `admin`, `user`). Include the role in the JWT claims. Create a middleware that checks if the user has the required role before allowing access (e.g., only admins can create products).

2. **Refresh Token Endpoint** - Implement a `POST /refresh` endpoint that accepts a valid (non-expired) token and returns a new token with a refreshed expiration time. This allows clients to stay logged in without re-entering credentials.

3. **API Key Authentication** - Add an alternative authentication method using API keys. Create an `api_keys` table and a `POST /api-keys` endpoint (authenticated) to generate keys. Modify the `AuthMiddleware` to accept either a Bearer token or an `X-API-Key` header.

4. **Token Revocation (Logout)** - Implement a `POST /logout` endpoint that adds the current token to a blacklist (in-memory map or database table). Modify the `AuthMiddleware` to check the blacklist before allowing access.

## Key Concepts

| Concept | Description |
|---------|-------------|
| **JWT** | JSON Web Token with three parts: header (algorithm), payload (claims), and signature. Stateless authentication -- the server does not need to store session data. |
| **Password Hashing** | Bcrypt is a one-way hashing algorithm designed for passwords. It includes a salt and a cost factor to resist brute-force attacks. |
| **Middleware** | A function that wraps HTTP handlers to add cross-cutting behavior (authentication, logging, etc.) without modifying individual handlers. |
| **Bearer Token** | An authentication scheme where the client sends the token in the `Authorization: Bearer <token>` header with each request. |

## Cleanup

```bash
docker compose down -v
```
