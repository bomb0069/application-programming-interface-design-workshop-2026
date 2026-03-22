# Lab 04 - CRUD In-Memory

Build a complete CRUD (Create, Read, Update, Delete) REST API backed by a thread-safe in-memory data store. This lab covers every fundamental operation you need when working with resources in a REST API.

## Learning Objectives

- Implement full CRUD operations for a resource
- Use proper HTTP methods (GET, POST, PUT, DELETE)
- Return correct HTTP status codes (200, 201, 204, 400, 404)
- Build a thread-safe in-memory data store using `sync.RWMutex`

## Getting Started

### Run Locally

```bash
go run main.go
```

### Run with Docker Compose

```bash
docker compose up --build
```

The server starts on http://localhost:8080.

## Test with curl

### Create a Todo (POST)

```bash
curl -s -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go"}' | jq
```

Expected response (`201 Created`):

```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": false
}
```

### List All Todos (GET)

```bash
curl -s http://localhost:8080/todos | jq
```

Expected response (`200 OK`):

```json
[
  {
    "id": 1,
    "title": "Learn Go",
    "completed": false
  }
]
```

### Get a Single Todo (GET by ID)

```bash
curl -s http://localhost:8080/todos/1 | jq
```

Expected response (`200 OK`):

```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": false
}
```

### Update a Todo (PUT)

```bash
curl -s -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}' | jq
```

Expected response (`200 OK`):

```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": true
}
```

### Delete a Todo (DELETE)

```bash
curl -s -X DELETE http://localhost:8080/todos/1 -w "\nHTTP Status: %{http_code}\n"
```

Expected response (`204 No Content`):

```
HTTP Status: 204
```

### Verify Deletion

```bash
curl -s http://localhost:8080/todos/1 | jq
```

Expected response (`404 Not Found`):

```json
{
  "error": "Todo not found"
}
```

## Code Walkthrough

### The Store Struct

```go
type Store struct {
    mu     sync.RWMutex
    todos  map[int]Todo
    nextID int
}
```

The `Store` holds all todos in a map keyed by ID. The `nextID` field acts as an auto-incrementing primary key. The `sync.RWMutex` protects concurrent access to the map.

### Thread Safety with sync.RWMutex

- **`RLock()` / `RUnlock()`** -- Used for read operations (`listTodos`, `getTodo`). Multiple goroutines can hold a read lock simultaneously, so read-heavy workloads stay fast.
- **`Lock()` / `Unlock()`** -- Used for write operations (`createTodo`, `updateTodo`, `deleteTodo`). A write lock is exclusive; no other readers or writers can proceed until it is released.

This read-write lock pattern is ideal when reads far outnumber writes, which is typical for most APIs.

### Handler Breakdown

| Handler       | Method   | Path           | Description                        |
|---------------|----------|----------------|------------------------------------|
| `listTodos`   | `GET`    | `/todos`       | Return all todos as a JSON array   |
| `createTodo`  | `POST`   | `/todos`       | Create a new todo from JSON body   |
| `getTodo`     | `GET`    | `/todos/{id}`  | Return a single todo by ID         |
| `updateTodo`  | `PUT`    | `/todos/{id}`  | Update fields on an existing todo  |
| `deleteTodo`  | `DELETE` | `/todos/{id}`  | Remove a todo by ID                |

### Status Codes Used

| Code  | Meaning        | When Used                                      |
|-------|----------------|-------------------------------------------------|
| `200` | OK             | Successful GET or PUT                           |
| `201` | Created        | Successful POST that creates a new resource     |
| `204` | No Content     | Successful DELETE with no response body         |
| `400` | Bad Request    | Invalid JSON body or missing required fields    |
| `404` | Not Found      | Requested todo ID does not exist                |

### Pointer Fields in Update Input

```go
var input struct {
    Title     *string `json:"title"`
    Completed *bool   `json:"completed"`
}
```

Using pointer types (`*string`, `*bool`) lets us distinguish between "field not provided" (`nil`) and "field set to zero value" (empty string or `false`). This allows partial updates -- only the fields included in the request body are changed.

## Exercises

1. **Add a "priority" field** -- Extend the `Todo` struct with a `Priority` field that accepts `"low"`, `"medium"`, or `"high"`. Validate the value on create and update.

2. **Add PATCH support** -- Implement `PATCH /todos/{id}` for partial updates. Compare it with PUT: PUT traditionally replaces the entire resource, while PATCH applies only the provided changes.

3. **List completed todos** -- Add a `GET /todos/completed` endpoint that returns only todos where `completed` is `true`.

4. **Search by title** -- Add a `GET /todos/search?q=keyword` endpoint that returns todos whose title contains the given keyword (case-insensitive).

## HTTP Methods Reference

| Method   | Description                         | Typical Status Codes |
|----------|-------------------------------------|----------------------|
| `GET`    | Retrieve a resource or collection   | 200, 404             |
| `POST`   | Create a new resource               | 201, 400             |
| `PUT`    | Replace / update an existing resource | 200, 400, 404      |
| `PATCH`  | Partially update an existing resource | 200, 400, 404      |
| `DELETE` | Remove a resource                   | 204, 404             |

## Cleanup

```bash
docker compose down
```
