# Lab 06 - Request Validation

Learn how to validate incoming request bodies in a Go REST API using the [go-playground/validator](https://github.com/go-playground/validator) library.

## Learning Objectives

- Validate request bodies before processing them
- Return structured validation errors with field-level messages
- Use the `go-playground/validator` library for declarative validation
- Understand common validation tags and how to combine them

## Prerequisites

- Go 1.24 or later
- Docker and Docker Compose
- curl or a similar HTTP client

## Getting Started

Start the application with Docker Compose:

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

## API Endpoints

| Method | Path             | Description          |
|--------|------------------|----------------------|
| GET    | /products        | List all products    |
| POST   | /products        | Create a product     |
| GET    | /products/{id}   | Get a product by ID  |
| PUT    | /products/{id}   | Update a product     |
| DELETE | /products/{id}   | Delete a product     |

## Test Examples

### Create a valid product

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Mouse",
    "price": 29.99,
    "category": "electronics",
    "sku": "WMSE1234"
  }' | jq
```

Expected response (201 Created):

```json
{
  "id": 1,
  "name": "Wireless Mouse",
  "price": 29.99,
  "category": "electronics",
  "sku": "WMSE1234"
}
```

### Missing required fields

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{}' | jq
```

Expected response (400 Bad Request):

```json
{
  "error": "Validation failed",
  "details": [
    { "field": "name", "message": "name is required" },
    { "field": "price", "message": "price must be greater than 0" },
    { "field": "category", "message": "category is required" },
    { "field": "sku", "message": "sku is required" }
  ]
}
```

### Price must be greater than zero

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Free Item",
    "price": 0,
    "category": "books",
    "sku": "FREE0001"
  }' | jq
```

Expected response (400 Bad Request):

```json
{
  "error": "Validation failed",
  "details": [
    { "field": "price", "message": "price must be greater than 0" }
  ]
}
```

### Invalid category

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mystery Item",
    "price": 9.99,
    "category": "toys",
    "sku": "TOYS0001"
  }' | jq
```

Expected response (400 Bad Request):

```json
{
  "error": "Validation failed",
  "details": [
    { "field": "category", "message": "category must be one of: electronics books clothing food" }
  ]
}
```

### Wrong SKU length

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Short SKU Item",
    "price": 5.00,
    "category": "food",
    "sku": "AB12"
  }' | jq
```

Expected response (400 Bad Request):

```json
{
  "error": "Validation failed",
  "details": [
    { "field": "sku", "message": "sku must be exactly 8 characters" }
  ]
}
```

### Multiple validation errors at once

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "X",
    "price": -5,
    "category": "toys",
    "sku": "AB!@"
  }' | jq
```

Expected response (400 Bad Request):

```json
{
  "error": "Validation failed",
  "details": [
    { "field": "name", "message": "name must be at least 2 characters" },
    { "field": "price", "message": "price must be greater than 0" },
    { "field": "category", "message": "category must be one of: electronics books clothing food" },
    { "field": "sku", "message": "sku must be exactly 8 characters" }
  ]
}
```

## Code Walkthrough

### Validation Tags on the Struct

The `CreateProductInput` struct uses `validate` tags to declare rules:

```go
type CreateProductInput struct {
    Name     string  `json:"name" validate:"required,min=2,max=100"`
    Price    float64 `json:"price" validate:"required,gt=0"`
    Category string  `json:"category" validate:"required,oneof=electronics books clothing food"`
    SKU      string  `json:"sku" validate:"required,len=8,alphanum"`
}
```

Each field has a `validate` tag with comma-separated rules. Multiple rules are combined -- all must pass for the field to be valid.

### The `formatValidationErrors` Function

When `validate.Struct()` returns an error, it contains a slice of `validator.ValidationErrors`. The `formatValidationErrors` function iterates over each error and produces a human-readable message based on the tag that failed:

```go
func formatValidationErrors(err error) []ValidationError {
    var errors []ValidationError
    for _, e := range err.(validator.ValidationErrors) {
        var msg string
        switch e.Tag() {
        case "required":
            msg = fmt.Sprintf("%s is required", strings.ToLower(e.Field()))
        case "min":
            msg = fmt.Sprintf("%s must be at least %s characters", strings.ToLower(e.Field()), e.Param())
        // ... other cases
        }
        errors = append(errors, ValidationError{
            Field:   strings.ToLower(e.Field()),
            Message: msg,
        })
    }
    return errors
}
```

This produces structured JSON errors that API consumers can programmatically handle.

### Using Validation in Handlers

Validation is applied in the handler after decoding the JSON body:

```go
func createProduct(w http.ResponseWriter, r *http.Request) {
    var input CreateProductInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
        return
    }

    if err := validate.Struct(input); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]interface{}{
            "error":   "Validation failed",
            "details": formatValidationErrors(err),
        })
        return
    }

    // ... proceed with database insert
}
```

## Validation Tags Reference

| Tag        | Description                              | Example                      |
|------------|------------------------------------------|------------------------------|
| `required` | Field must not be zero value             | `validate:"required"`        |
| `min`      | Minimum length (string) or value (number)| `validate:"min=2"`           |
| `max`      | Maximum length (string) or value (number)| `validate:"max=100"`         |
| `len`      | Exact length                             | `validate:"len=8"`           |
| `gt`       | Greater than                             | `validate:"gt=0"`            |
| `gte`      | Greater than or equal                    | `validate:"gte=1"`           |
| `lt`       | Less than                                | `validate:"lt=1000"`         |
| `lte`      | Less than or equal                       | `validate:"lte=999"`         |
| `oneof`    | Must be one of the listed values         | `validate:"oneof=a b c"`     |
| `alphanum` | Alphanumeric characters only             | `validate:"alphanum"`        |
| `email`    | Valid email address                      | `validate:"email"`           |
| `url`      | Valid URL                                | `validate:"url"`             |
| `uuid`     | Valid UUID                               | `validate:"uuid"`            |

For the full list, see the [validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10).

## Exercises

### Exercise 1: Add Email Validation

Add a `contact_email` field to the product with the `email` validation tag:

- Add a `ContactEmail` field to both `Product` and `CreateProductInput` structs
- Use the tag `validate:"required,email"`
- Update the database table and queries to include the new column
- Test with valid and invalid email addresses

### Exercise 2: Custom Price Precision Validation

Add a custom validation to ensure the price has no fractions of cents (must be a multiple of 0.01):

- Register a custom validation function using `validate.RegisterValidation("price_precision", ...)`
- The validator should check that `price * 100` has no fractional remainder
- Add the `price_precision` tag to the `Price` field
- Test with values like `29.99` (valid) and `29.999` (invalid)

### Exercise 3: Optional Update Fields

Create an `UpdateProductInput` struct where all fields are optional but validated when present:

- Use pointer types (`*string`, `*float64`) for all fields
- Use the `omitempty` tag so validation only runs when a value is provided: `validate:"omitempty,min=2,max=100"`
- Modify the `updateProduct` handler to only update fields that are provided
- Build a dynamic SQL UPDATE query based on which fields are non-nil

### Exercise 4: Cross-Field Validation

Add a rule: if `category` is "electronics", then `price` must be greater than 10:

- Register a struct-level validation using `validate.RegisterStructValidation(...)`
- Access multiple fields within the validation function
- Return a meaningful error message when the rule fails
- Test with `category: "electronics"` and `price: 5` to verify the error

## Key Concepts

### Input Validation

Never trust data coming from API clients. Always validate request bodies before processing them. Validation serves as the first line of defense against invalid or malicious data reaching your database.

### Validation Tags

The `go-playground/validator` library uses struct tags to declare validation rules. Tags are comma-separated and attached to struct fields. This declarative approach keeps validation logic close to the data definition, making it easy to see what rules apply to each field.

### Structured Error Responses

Rather than returning a single error string, return an array of field-level errors. This allows API consumers to map errors to specific form fields and display targeted feedback to users. Each error includes the field name and a human-readable message.

## Cleanup

Stop and remove the containers:

```bash
docker compose down
```

To also remove the database volume:

```bash
docker compose down -v
```
