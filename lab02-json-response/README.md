# Lab 02 - JSON Response

## Learning Objectives

By the end of this lab, you will be able to:

- Define Go structs with JSON tags
- Marshal Go data structures to JSON
- Handle multiple endpoints that return different JSON shapes
- Understand how JSON struct tags control serialization

## Prerequisites

- Go 1.24 or later installed
- Docker and Docker Compose installed
- Basic understanding of HTTP (covered in Lab 01)
- A terminal and a text editor

## Getting Started

1. Navigate to the lab directory:

```bash
cd lab02-json-response
```

2. Start the server using Docker Compose:

```bash
docker-compose up --build
```

The server will start on port **8080**.

## Test the Endpoints

Open a new terminal and run the following commands:

### Get all books

```bash
curl http://localhost:8080/books
```

Expected response:

```json
[
  {"id":1,"title":"The Go Programming Language","author":"Alan Donovan","year":2015},
  {"id":2,"title":"Go in Action","author":"William Kennedy","year":2015},
  {"id":3,"title":"Learning Go","author":"Jon Bodner","year":2021}
]
```

### Get the book count

```bash
curl http://localhost:8080/books/count
```

Expected response:

```json
{"count":3}
```

### Health check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

> **Tip:** Pipe the output through `jq` for pretty-printed JSON:
> ```bash
> curl -s http://localhost:8080/books | jq .
> ```

## Code Walkthrough

Open `main.go` and follow along with the explanation below.

### 1. Struct Definition with JSON Tags

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}
```

A **struct** in Go is a typed collection of fields. Each field here has a **JSON struct tag** (the backtick-enclosed annotation). Struct tags are metadata that tell the `encoding/json` package how to map struct fields to JSON keys.

Without struct tags, `json.Marshal` would use the Go field names directly (e.g., `"ID"`, `"Title"`). The tags let you control the exact JSON key names so that they follow the conventional lowercase style used in JSON APIs.

### 2. The Books Slice

```go
var books = []Book{
    {ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015},
    {ID: 2, Title: "Go in Action", Author: "William Kennedy", Year: 2015},
    {ID: 3, Title: "Learning Go", Author: "Jon Bodner", Year: 2021},
}
```

This is a **package-level variable** holding a slice of `Book` structs. In a real application, this data would come from a database. For this lab, we use an in-memory slice to keep things simple and focus on JSON serialization.

### 3. Encoding JSON with `json.NewEncoder`

```go
func booksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}
```

`json.NewEncoder(w)` creates a new JSON encoder that writes directly to the `http.ResponseWriter`. Calling `.Encode(books)` converts the Go slice into JSON and writes it to the response body in a single step.

An alternative approach is to use `json.Marshal(books)` which returns a `[]byte`, and then write that to `w` manually. The encoder approach is preferred for HTTP handlers because it streams the output directly without buffering the entire response in memory.

### 4. How Struct Tags Control JSON Field Names

| Go Field | Struct Tag       | JSON Key   |
|----------|------------------|------------|
| `ID`     | `json:"id"`      | `"id"`     |
| `Title`  | `json:"title"`   | `"title"`  |
| `Author` | `json:"author"`  | `"author"` |
| `Year`   | `json:"year"`    | `"year"`   |

The struct tag format is `` `json:"keyname"` ``. Whatever string you put between the quotes becomes the key in the resulting JSON output. This is the standard way to ensure your API returns consistently formatted JSON regardless of Go naming conventions.

## Exercises

Try these exercises to deepen your understanding. Each one builds on the starter code.

### Exercise 1: Add a Genre Field

Add a new field `Genre` to the `Book` struct with the JSON tag `"genre"`. Then update the existing book data to include genres.

**Steps:**

1. Add `Genre string \`json:"genre"\`` to the `Book` struct.
2. Add a genre to each book in the `books` slice (e.g., `"Programming"`, `"Programming"`, `"Programming"`).
3. Rebuild and test:

```bash
docker-compose up --build
curl -s http://localhost:8080/books | jq .
```

4. Verify that each book now includes a `"genre"` field in the JSON output.

### Exercise 2: Create a `/books/summary` Endpoint

Create a new endpoint at `/books/summary` that returns an array of strings. Each string should describe a book in the format:

```
"{Title} by {Author} ({Year})"
```

**Expected response:**

```json
[
  "The Go Programming Language by Alan Donovan (2015)",
  "Go in Action by William Kennedy (2015)",
  "Learning Go by Jon Bodner (2021)"
]
```

**Hints:**

- Register a new handler with `http.HandleFunc("/books/summary", booksSummaryHandler)`.
- Use `fmt.Sprintf` to format each string.
- Build a `[]string` slice and encode it as JSON.

### Exercise 3: Create a Nested Response Structure

Instead of returning a plain array from `/books`, create a wrapper response that includes both the data and metadata.

**Target response structure:**

```json
{
  "data": [
    {"id":1,"title":"The Go Programming Language","author":"Alan Donovan","year":2015},
    {"id":2,"title":"Go in Action","author":"William Kennedy","year":2015},
    {"id":3,"title":"Learning Go","author":"Jon Bodner","year":2021}
  ],
  "count": 3
}
```

**Hints:**

- Define a new struct (e.g., `BooksResponse`) with fields `Data` and `Count`.
- Use appropriate JSON tags: `json:"data"` and `json:"count"`.
- Populate the struct and encode it in the handler.

### Exercise 4: Experiment with JSON Tags

Try modifying struct tags to observe different behaviors:

1. **Hide a field with `json:"-"`**: Change the `Year` field tag to `json:"-"` and observe that it no longer appears in the JSON output.

2. **Omit empty values with `omitempty`**: Change the `Title` field tag to `json:"title,omitempty"`. Then add a book with an empty title to the slice and see that the `"title"` key is omitted for that entry.

3. **Combine both**: Try `json:"author,omitempty"` and add a book with no author to see the effect.

After experimenting, remember to restore the original tags before moving on.

## Key Concepts

### JSON Marshaling

**Marshaling** is the process of converting a Go data structure (struct, slice, map, etc.) into JSON format. The `encoding/json` package provides two main ways to do this:

- `json.Marshal(v)` -- returns `[]byte` and an `error`
- `json.NewEncoder(w).Encode(v)` -- writes JSON directly to an `io.Writer`

For HTTP handlers, `json.NewEncoder` is generally preferred because it writes directly to the response without an intermediate buffer.

### Struct Tags

Struct tags are string annotations on struct fields that provide metadata for packages like `encoding/json`. The general syntax is:

```go
FieldName Type `tagname:"value"`
```

Common JSON struct tag options:

| Tag                        | Effect                                         |
|----------------------------|-------------------------------------------------|
| `json:"name"`             | Sets the JSON key to `"name"`                   |
| `json:"-"`                | Excludes the field from JSON output              |
| `json:"name,omitempty"`   | Omits the field if it has a zero value            |
| `json:",omitempty"`       | Uses the Go field name but omits if zero value    |

### Content-Type Header

Setting `Content-Type: application/json` in the response header tells the client that the response body is JSON. This is important because:

- Browsers and HTTP clients use it to parse the response correctly.
- API tools like Postman and curl use it for formatting and syntax highlighting.
- It is part of the HTTP specification for proper content negotiation.

Always set the `Content-Type` header **before** writing the response body.

## Cleanup

When you are finished with the lab, stop the running containers:

```bash
docker-compose down
```
