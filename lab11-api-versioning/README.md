# Lab 11 - API Versioning

## Learning Objectives

- Implement URL-based API versioning
- Handle breaking vs non-breaking changes
- Run multiple API versions simultaneously
- Version migration strategies

## Getting Started

```bash
docker compose up --build
```

The API server starts on http://localhost:8080 with two versioned endpoints: `/v1` and `/v2`.

## Test Both Versions

### V1 - Simple Flat Responses

List all products (flat array):

```bash
curl http://localhost:8080/v1/products
```

```json
[
  {"id": 1, "name": "Laptop", "price": 999.99, "category": "electronics"},
  {"id": 2, "name": "Go Book", "price": 39.99, "category": "books"},
  {"id": 3, "name": "T-Shirt", "price": 19.99, "category": "clothing"}
]
```

Get a single product:

```bash
curl http://localhost:8080/v1/products/1
```

```json
{"id": 1, "name": "Laptop", "price": 999.99, "category": "electronics"}
```

### V2 - Enhanced Wrapped Responses

List all products (wrapped with `data` and `version`):

```bash
curl http://localhost:8080/v2/products
```

```json
{
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "price": 999.99,
      "category": "electronics",
      "description": "A powerful laptop for developers",
      "tags": ["portable", "computing"]
    }
  ],
  "version": "2.0"
}
```

Get a single product:

```bash
curl http://localhost:8080/v2/products/1
```

```json
{
  "data": {
    "id": 1,
    "name": "Laptop",
    "price": 999.99,
    "category": "electronics",
    "description": "A powerful laptop for developers",
    "tags": ["portable", "computing"]
  },
  "version": "2.0"
}
```

Create a product with description and tags (V2 only):

```bash
curl -X POST http://localhost:8080/v2/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Keyboard", "price": 79.99, "category": "electronics", "description": "Mechanical keyboard", "tags": ["peripheral", "mechanical"]}'
```

## V1 vs V2 Response Comparison

| Aspect          | V1                          | V2                                      |
|-----------------|-----------------------------|-----------------------------------------|
| Response shape  | Flat object / array         | Wrapped in `{"data": ..., "version": "2.0"}` |
| Fields          | id, name, price, category   | id, name, price, category, description, tags |
| Create input    | name, price, category       | name, price, category, description, tags |

## Breaking vs Non-Breaking Changes

### Non-Breaking Changes (safe to add without a new version)

- Adding new optional fields to the response
- Adding new endpoints
- Adding optional query parameters

### Breaking Changes (require a new version)

- Changing the response structure (e.g., wrapping in `{data, version}`)
- Removing or renaming existing fields
- Changing field types (e.g., `price` from number to string)
- Altering error response format

V2 in this lab introduces **breaking changes**: the response is wrapped in an envelope with `data` and `version` fields, and new fields (`description`, `tags`) are added to the product model.

## Exercises

1. **Header-based versioning** -- Add support for selecting the API version via the `Accept` header:
   ```
   Accept: application/vnd.api.v1+json
   Accept: application/vnd.api.v2+json
   ```

2. **Deprecation warning** -- Add middleware to v1 routes that sets a deprecation header:
   ```
   X-API-Deprecation: v1 will be removed on 2027-01-01
   ```

3. **V3 with cursor-based pagination** -- Create a `/v3/products` endpoint that uses cursor-based pagination instead of returning all results:
   ```
   GET /v3/products?cursor=abc123&limit=10
   ```

4. **Version redirect middleware** -- Write middleware that redirects unversioned requests (`/products`) to the latest version (`/v2/products`).

## Key Concepts

- **URL Versioning** -- The version is part of the URL path (`/v1/`, `/v2/`). This is the most common and visible approach. Clients explicitly choose which version to call.

- **Breaking vs Non-Breaking Changes** -- Non-breaking (additive) changes can be made within the same version. Breaking changes (structural, removals, renames) require a new version to avoid disrupting existing clients.

- **Version Coexistence** -- Multiple API versions run side by side in the same service, sharing the same database. This lets clients migrate at their own pace.

- **Deprecation Strategy** -- Older versions should be marked as deprecated with a timeline for removal. Use response headers, documentation, and client communication to drive migration.

## Cleanup

```bash
docker compose down -v
```
