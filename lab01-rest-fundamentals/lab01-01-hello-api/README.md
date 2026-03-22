# Lab 01 - Hello API

## Learning Objectives

- Understand how to create a basic HTTP server in Go
- Return JSON responses
- Use Docker to containerize a Go application

## Prerequisites

- Go 1.24+
- Docker and Docker Compose

## Getting Started

1. Build and run the application using Docker Compose:

```bash
docker-compose up --build
```

2. Test the API by sending a request:

```bash
curl http://localhost:8080/
```

You should see the following response:

```json
{"message": "Hello, World!"}
```

## Explain the Code

Let's walk through `main.go` step by step.

### Importing Packages

```go
import (
    "encoding/json"
    "log"
    "net/http"
)
```

- `net/http` is Go's built-in HTTP package. It provides everything you need to create an HTTP server and handle requests without any external dependencies.
- `encoding/json` is used to encode Go data structures into JSON format.
- `log` is used to print log messages to the console.

### Registering a Route with HandleFunc

```go
http.HandleFunc("/", helloHandler)
```

`http.HandleFunc` registers a handler function for a given URL pattern. In this case, any request to the root path `/` will be handled by the `helloHandler` function.

### Starting the Server

```go
log.Println("Server starting on :8080")
log.Fatal(http.ListenAndServe(":8080", nil))
```

`http.ListenAndServe` starts an HTTP server on port 8080. The second argument `nil` tells Go to use the default request multiplexer (the one where we registered our handler with `HandleFunc`). If the server fails to start, `log.Fatal` will print the error and exit the program.

### The Handler Function

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Hello, World!",
    })
}
```

Every HTTP handler in Go receives two arguments:

- `w http.ResponseWriter` - used to write the response back to the client.
- `r *http.Request` - contains all the information about the incoming request (method, headers, body, URL, etc.).

Inside the handler:

1. We set the `Content-Type` header to `application/json` so the client knows the response is JSON.
2. We use `json.NewEncoder(w).Encode(...)` to serialize a Go map into JSON and write it directly to the response. This is equivalent to manually calling `json.Marshal` and then `w.Write`, but more concise.

## Exercises

1. **Add a `/health` endpoint** - Create a new handler that responds with `{"status": "ok"}` when a client sends a request to `/health`. This is a common pattern used by load balancers and orchestrators to check if a service is running.

2. **Add a `/time` endpoint** - Create a new handler that returns the current server time. You will need to import the `time` package and use `time.Now().Format(time.RFC3339)` to get a formatted timestamp. Return it as `{"current_time": "2026-03-09T10:30:00Z"}`.

3. **Add your name to the response** - Modify the `helloHandler` to include an `author` field in the JSON response: `{"message": "Hello, World!", "author": "Your Name"}`.

## Key Concepts

### HTTP Methods

HTTP defines several request methods (also called verbs). The most common ones are:

- **GET** - Retrieve data from the server
- **POST** - Send data to the server to create a resource
- **PUT** - Update an existing resource
- **DELETE** - Remove a resource

In this lab, our handler responds to all HTTP methods. In later labs, you will learn how to handle specific methods.

### JSON Response

JSON (JavaScript Object Notation) is the most common format for API responses. In Go, you can use the `encoding/json` package to convert Go data structures (maps, structs) into JSON strings.

### Content-Type Header

The `Content-Type` header tells the client what format the response body is in. For JSON APIs, you should always set this to `application/json`. Without this header, clients may not correctly interpret the response.

## Cleanup

When you are done, stop and remove the containers:

```bash
docker-compose down
```
