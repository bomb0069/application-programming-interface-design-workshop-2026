# Lab 05 - CRUD with Database

In this lab, we replace the in-memory store from Lab 04 with a real **PostgreSQL** database. The API endpoints remain identical, but data now persists across restarts.

## Learning Objectives

- Connect a Go application to PostgreSQL using `database/sql` and `lib/pq`
- Execute SQL queries: `SELECT`, `INSERT`, `UPDATE`, `DELETE`
- Use **parameterized queries** (`$1`, `$2`, ...) to prevent SQL injection
- Orchestrate multiple services with **Docker Compose**
- Configure **database health checks** so the API waits for PostgreSQL to be ready

## Getting Started

Start both the API and database with a single command:

```bash
docker-compose up --build
```

Docker Compose will:

1. Start a PostgreSQL 16 container
2. Wait until PostgreSQL passes its health check (`pg_isready`)
3. Build and start the Go API, which connects to PostgreSQL and auto-creates the `todos` table

The API is available at `http://localhost:8080`.

## Test with curl

These are the same endpoints as Lab 04 -- the only difference is that data is now stored in PostgreSQL.

**Create a todo:**

```bash
curl -s -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries"}' | jq
```

**List all todos:**

```bash
curl -s http://localhost:8080/todos | jq
```

**Get a single todo:**

```bash
curl -s http://localhost:8080/todos/1 | jq
```

**Update a todo:**

```bash
curl -s -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}' | jq
```

**Delete a todo:**

```bash
curl -s -X DELETE http://localhost:8080/todos/1 -w "\nHTTP Status: %{http_code}\n"
```

## Code Walkthrough

### database/sql

Go's standard library provides the `database/sql` package, which offers a generic interface around SQL databases. It handles connection pooling, prepared statements, and transactions.

```go
import "database/sql"

db, err = sql.Open("postgres", dsn)
```

### lib/pq Driver

The `lib/pq` package is a PostgreSQL driver for `database/sql`. It is imported with a blank identifier so its `init()` function registers the driver:

```go
import _ "github.com/lib/pq"
```

### Connection String

The connection string (DSN) follows the PostgreSQL URI format:

```
postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable
```

In Docker Compose, the hostname is the service name (`db`), so the API uses:

```
postgres://postgres:postgres@db:5432/workshop?sslmode=disable
```

### createTable

On startup, the application runs `CREATE TABLE IF NOT EXISTS` to ensure the `todos` table exists. This is a simple approach suitable for development. In production, you would use a migration tool.

### Parameterized Queries

All user input is passed through parameterized placeholders (`$1`, `$2`, `$3`) instead of string concatenation. This prevents SQL injection:

```go
db.QueryRow("SELECT id, title, completed FROM todos WHERE id = $1", id)
db.Exec("UPDATE todos SET title = $1, completed = $2 WHERE id = $3", todo.Title, todo.Completed, id)
```

### sql.ErrNoRows

When a `QueryRow` finds no matching record, it returns `sql.ErrNoRows`. We check for this to return a proper 404 response:

```go
if err == sql.ErrNoRows {
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"error": "Todo not found"})
    return
}
```

## Comparing with Lab 04

| Aspect | Lab 04 (In-Memory) | Lab 05 (Database) |
|---|---|---|
| Storage | Go slice + mutex | PostgreSQL table |
| Persistence | Lost on restart | Survives restarts |
| ID Generation | Manual counter | `SERIAL` (auto-increment) |
| Concurrency | `sync.Mutex` | Database handles it |
| Infrastructure | Single binary | Docker Compose (API + DB) |
| Not Found | Linear search | `sql.ErrNoRows` |

## Exercises

1. **Add a `created_at` column** -- Add a `TIMESTAMP DEFAULT NOW()` column to the `todos` table. Include it in the `Todo` struct and return it in API responses.

2. **Add a `description` column** -- Add an optional `TEXT` column for a longer description. Update the create and update handlers to accept and return it.

3. **Add a seed data script** -- Create a `seed.sql` file that inserts sample todos, and mount it into the PostgreSQL container at `/docker-entrypoint-initdb.d/` so it runs on first startup.

4. **Connect with psql** -- Use the PostgreSQL CLI to inspect your data directly:
   ```bash
   docker-compose exec db psql -U postgres workshop
   ```
   Try running `SELECT * FROM todos;` and `\d todos` to see the table schema.

## Key Concepts

### database/sql

Go's `database/sql` package provides a generic interface for SQL databases. It manages connection pooling automatically -- you call `sql.Open()` once and share the `*sql.DB` across your application.

### Parameterized Queries

Never build SQL strings by concatenating user input. Always use placeholders (`$1`, `$2` for PostgreSQL) and pass values as arguments. The driver handles escaping and type conversion.

### Docker Compose Services

Docker Compose lets you define multi-container applications in a single YAML file. Each service gets its own container, network alias, and configuration. Services can reference each other by name (e.g., the API connects to `db:5432`).

### Database Health Checks

The `healthcheck` configuration tells Docker Compose how to determine if a container is ready. The `depends_on` condition `service_healthy` ensures the API does not start until PostgreSQL is accepting connections.

## Cleanup

Stop and remove all containers and the database volume:

```bash
docker-compose down -v
```

The `-v` flag removes the `pgdata` volume, which deletes all stored data. Omit `-v` if you want to keep your data for next time.
