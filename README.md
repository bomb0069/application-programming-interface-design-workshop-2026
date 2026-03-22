# API Design Workshop 2026

Welcome to the **API Design Workshop**! This is a hands-on, progressive workshop that takes you from building your first HTTP endpoint to mastering advanced API technologies — all in **Go**.

## What is this workshop?

This workshop contains **20 self-contained labs** organized into three parts. You will build real, working APIs from scratch — starting with simple JSON responses and progressing through database-backed CRUD, authentication, GraphQL, gRPC, WebSockets, and message queues.

Each lab includes:
- A **README** with learning objectives, step-by-step instructions, and exercises
- Working **Go source code** you can run, read, and modify
- A **docker-compose.yml** (or one per language, in multi-language labs) that spins up everything you need with one command

No prior Go experience is required, but basic programming knowledge is helpful.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- A terminal and a code editor
- `curl` (or an API client like Postman)
- [Go 1.24+](https://go.dev/dl/) (optional — all labs run inside Docker)

## Getting Started

```bash
# Pick a lab
cd lab01-hello-api

# Start all services
docker-compose up --build

# Follow the README inside the lab
# Test with curl, browser, or the tools provided

# Clean up when done
docker-compose down -v
```

Each lab is **independent** — you can start from any lab, though we recommend going in order.

---

## Labs

### Part 1: RESTful API Fundamentals

Build a solid foundation in REST API design — from your first endpoint to a fully authenticated, documented API.

| # | Lab | Description |
|---|-----|-------------|
| 01 | [Hello API](lab01-hello-api/) | Your first Go HTTP server — return `{"message": "Hello, World!"}` and learn `net/http` basics |
| 02 | [JSON Response](lab02-json-response/) | Return structured JSON with Go structs, learn JSON marshaling and struct tags |
| 03 | [Path Parameters](lab03-path-parameter/) | Extract URL parameters like `/items/{id}` using the chi router, handle 404s |
| 04 | [CRUD In-Memory](lab04-crud-in-memory/) | Build a full Create/Read/Update/Delete API with an in-memory data store |
| 05 | [CRUD with Database](lab05-crud-with-database/) | Replace the in-memory store with PostgreSQL — write real SQL queries |
| 06 | [Request Validation](lab06-request-validation/) | Validate incoming data with `go-playground/validator`, return structured errors |
| 07 | [Error Handling](lab07-error-handling/) | Build consistent error responses with centralized error types and middleware |
| 08 | [Swagger Documentation](lab08-swagger-documentation/) | Document your API with OpenAPI 3.0 and serve interactive Swagger UI |
| 09 | [Pagination & Filtering](lab09-pagination-and-filtering/) | Add `?page=`, `?sort=`, `?category=` query parameters with response metadata |
| 10 | [Authentication](lab10-authentication/) | Secure your API with JWT tokens, bcrypt password hashing, and auth middleware |

### Part 2: Advanced REST Topics

Take your REST APIs further with versioning, rate limiting, and file handling.

| # | Lab | Description |
|---|-----|-------------|
| 11 | [API Versioning](lab11-api-versioning/) | Run `/v1` and `/v2` side by side — understand breaking vs non-breaking changes |
| 12 | [Rate Limiting & CORS](lab12-rate-limiting-and-cors/) | Implement token bucket rate limiting and configure CORS for cross-origin access |
| 13 | [File Upload & Download](lab13-file-upload-download/) | Upload files via multipart form, store them in MinIO (S3-compatible), serve downloads |

### Part 3: Beyond REST

Explore alternative API technologies — GraphQL, WebSockets, gRPC, and messaging systems.

| # | Lab | Description |
|---|-----|-------------|
| 14 | [GraphQL](lab14-graphql/) | Build a GraphQL API with schemas, queries, mutations, and a built-in Playground |
| 15 | [Webhooks](lab15-webhook/) | Send event-driven notifications with HMAC signature verification and retry logic |
| 16 | [WebSocket](lab16-websocket/) | Build a real-time chat app with bidirectional WebSocket communication |
| 17 | [gRPC Basics](lab17-grpc-basics/) | Define services with Protocol Buffers, build a gRPC server and client |
| 18 | [gRPC Advanced](lab18-grpc-advanced/) | Server/client/bidirectional streaming and a REST-to-gRPC gateway |
| 19 | [Message Queue](lab19-message-queue/) | Async communication with RabbitMQ — publisher, consumer, exchanges, and queues |
| 20 | [MQTT](lab20-mqtt/) | IoT-style pub/sub with Mosquitto — topics, QoS levels, and wildcard subscriptions |

---

## Learning Path

The labs are designed to be completed in order. Each part builds on patterns from the previous one:

```
Part 1: RESTful Fundamentals
  Lab 01-04  →  Go HTTP basics, routing, CRUD
  Lab 05-07  →  Database, validation, error handling
  Lab 08-10  →  Documentation, pagination, authentication
        ↓
Part 2: Advanced REST
  Lab 11-13  →  Versioning, rate limiting, file handling
        ↓
Part 3: Beyond REST
  Lab 14-16  →  GraphQL, webhooks, WebSocket
  Lab 17-18  →  gRPC (unary + streaming)
  Lab 19-20  →  Message queues (RabbitMQ, MQTT)
```

## Lab Structure

Every lab follows the same structure so you always know where to look:

```
labXX-topic-name/
├── README.md              # Instructions, explanations, exercises
├── docker-compose.yml     # Run everything with one command
├── Dockerfile             # Container build for the Go app
├── go.mod / go.sum        # Go dependencies
├── main.go                # Application entry point
└── ...                    # Additional files (handlers, configs, static assets)
```

> **Note:** Some labs (like Lab 11) contain sub-labs with multiple language implementations. In those cases, the `docker-compose.yml` lives inside `golang/` or `dotnet/` subdirectories rather than at the lab root. See individual lab READMEs for details.

## Services & Ports Quick Reference

| Labs | Services | Ports |
|------|----------|-------|
| 01-04 | Go API | `8080` |
| 05-10 | Go API, PostgreSQL | `8080`, `5432` |
| 08 | Go API, PostgreSQL, Swagger UI | `8080`, `5432`, `8081` |
| 11 (sub-labs 01-07) | Go/.NET API, PostgreSQL | `8080`, `5432` |
| 11 (sub-lab 08) | Go/.NET API, PostgreSQL, Prometheus, Grafana | `8080`, `5432`, `9090`, `3000` |
| 12 | Go API, PostgreSQL | `8080`, `5432` |
| 13 | Go API, PostgreSQL, MinIO | `8080`, `5432`, `9000`, `9001` |
| 14 | Go API (GraphQL Playground), PostgreSQL | `8080`, `5432` |
| 15 | Webhook Sender, Receiver, PostgreSQL | `8080`, `9090`, `5432` |
| 16 | Go API (Chat Web UI) | `8080` |
| 17 | gRPC Server, Client, gRPCUI | `50051`, `8080` |
| 18 | gRPC Server, Client, Gateway, gRPCUI | `50051`, `8080`, `8081` |
| 19 | Publisher, Consumer, RabbitMQ | `8080`, `5672`, `15672` |
| 20 | Publisher, Subscriber, Mosquitto | `1883` |

## Troubleshooting

- **Port conflict**: If a port is already in use, stop the previous lab with `docker-compose down` before starting a new one.
- **Build cache**: If something looks stale, rebuild with `docker-compose up --build --force-recreate`.
- **Database issues**: Reset the database volume with `docker-compose down -v` to start fresh.
