# API Design Workshop 2026

Welcome to the **API Design Workshop**! A hands-on, progressive workshop that takes you from building your first HTTP endpoint to mastering advanced API technologies.

Each lab includes:
- A **README** with learning objectives, step-by-step instructions, and exercises
- Working **Go source code** you can run, read, and modify
- A **docker-compose.yml** that spins up everything you need with one command

No prior Go experience is required, but basic programming knowledge is helpful.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- A terminal and a code editor
- `curl` (or an API client like Postman)
- [Go 1.24+](https://go.dev/dl/) (optional — all labs run inside Docker)

## Getting Started

```bash
# Pick a lab group and sub-lab
cd lab01-rest-fundamentals/lab01-01-hello-api
docker compose up --build

# For multi-language labs (e.g., API versioning)
cd lab04-api-versioning/lab04-01-url-path-versioning/golang
docker compose up --build

# Clean up when done
docker compose down -v
```

---

## Labs

### Part 1: REST API

#### Group 01: [REST Fundamentals](lab01-rest-fundamentals/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 01-01 | [Hello API](lab01-rest-fundamentals/lab01-01-hello-api/) | Your first Go HTTP server — `net/http` basics | ✅ |
| 01-02 | [JSON Response](lab01-rest-fundamentals/lab01-02-json-response/) | Structured JSON with Go structs and struct tags | ✅ |
| 01-03 | [CRUD In-Memory](lab01-rest-fundamentals/lab01-03-crud-in-memory/) | Full CRUD API with in-memory data store | ✅ |
| 01-04 | [CRUD with Database](lab01-rest-fundamentals/lab01-04-crud-with-database/) | PostgreSQL-backed CRUD with real SQL | ✅ |
| 01-05 | [File Upload & Download](lab01-rest-fundamentals/lab01-05-file-upload-download/) | Multipart upload, MinIO (S3-compatible) storage | ✅ |

#### Group 02: [REST Design Conventions](lab02-rest-design-conventions/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 02-01 | Resource Naming | Nouns, plural, kebab-case, nested resources | ❌ |
| 02-02 | [Path & Query Parameters](lab02-rest-design-conventions/lab02-02-path-and-query-parameters/) | URL parameters, path vs query decision rules | ✅ |
| 02-03 | [Request Validation](lab02-rest-design-conventions/lab02-03-request-validation/) | `go-playground/validator` with structured errors | ✅ |
| 02-04 | [Error Handling](lab02-rest-design-conventions/lab02-04-error-handling/) | Centralized error types and middleware | ✅ |
| 02-05 | Response Envelope & Status Codes | Unified envelope, 200 vs 404 for empty search | ❌ |
| 02-06 | Data Formats | ISO 8601 dates, nulls, enums, money/decimal | ❌ |
| 02-07 | [Pagination & Filtering](lab02-rest-design-conventions/lab02-07-pagination-and-filtering/) | `?page=`, `?sort=`, `?category=` with metadata | ✅ |
| 02-08 | [Swagger Documentation](lab02-rest-design-conventions/lab02-08-swagger-documentation/) | OpenAPI 3.0 with interactive Swagger UI | ✅ |

#### Group 03: [API Security](lab03-api-security/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 03-01 | [Authentication](lab03-api-security/lab03-01-authentication/) | JWT tokens, bcrypt, auth middleware | ✅ |
| 03-02 | [Rate Limiting & CORS](lab03-api-security/lab03-02-rate-limiting-and-cors/) | Token bucket rate limiting, CORS headers | ✅ |
| 03-03 | Sensitive Data Handling | Data masking, field-level security, PII rules | ❌ |
| 03-04 | API Key Management | Key lifecycle, rotation, revocation | ❌ |
| 03-05 | API Gateway | Centralized auth, rate limiting, routing | ❌ |

#### Group 04: [API Versioning](lab04-api-versioning/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 04-01 | [URL Path Versioning](lab04-api-versioning/lab04-01-url-path-versioning/) | `/api/v1/resource` — the industry default | ✅ |
| 04-02 | [Query Parameter](lab04-api-versioning/lab04-02-query-parameter-versioning/) | `?api-version=1` versioning | ✅ |
| 04-03 | [Header Versioning](lab04-api-versioning/lab04-03-header-versioning/) | `X-Api-Version` header | ✅ |
| 04-04 | [Content Negotiation](lab04-api-versioning/lab04-04-content-negotiation/) | Media type versioning via Accept header | ✅ |
| 04-05 | [Evolving API](lab04-api-versioning/lab04-05-evolving-api/) | Additive changes without versioning | ✅ |
| 04-06 | [Combining Strategies](lab04-api-versioning/lab04-06-combining-strategies/) | URL + query + header with priority | ✅ |
| 04-07 | [Breaking Changes](lab04-api-versioning/lab04-07-breaking-changes-and-deprecation/) | Deprecation/Sunset headers, 410 Gone | ✅ |
| 04-08 | [Lifecycle & Observability](lab04-api-versioning/lab04-08-version-lifecycle-and-observability/) | Prometheus metrics, Grafana dashboards | ✅ |

#### Group 05: Advanced Search

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 05-01 | POST /search | Complex filters with JSON body | ❌ |
| 05-02 | Saved Search | Two-phase pattern, TTL, 410 Gone | ❌ |
| 05-03 | Search Security | Scoping, field projection, rate limiting | ❌ |
| 05-04 | Capstone | Multi-tenant search with audit trail | ❌ |

#### Group 06: Performance & Resilience

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 06-01 | Caching with Redis | Cache-aside, write-through, ETags | ❌ |
| 06-02 | Circuit Breaker | Closed/Open/Half-Open, fallbacks, health checks | ❌ |

#### Group 07: Observability

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 07-01 | Structured Logging | JSON logs, correlation ID, log levels | ❌ |
| 07-02 | Metrics | Prometheus + Grafana, RED method | ❌ |
| 07-03 | Distributed Tracing | OpenTelemetry + Jaeger | ❌ |

---

### Part 2: Beyond REST

#### Group 08: [GraphQL](lab08-graphql/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 08-01 | [GraphQL](lab08-graphql/lab08-01-graphql/) | Schemas, queries, mutations, Playground | ✅ |

#### Group 09: [Real-Time APIs](lab09-real-time-apis/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 09-01 | [Webhook](lab09-real-time-apis/lab09-01-webhook/) | HMAC verification, retry logic | ✅ |
| 09-02 | [WebSocket](lab09-real-time-apis/lab09-02-websocket/) | Bidirectional real-time chat | ✅ |

#### Group 10: [gRPC](lab10-grpc/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 10-01 | [gRPC Basics](lab10-grpc/lab10-01-grpc-basics/) | Protocol Buffers, server and client | ✅ |
| 10-02 | [gRPC Advanced](lab10-grpc/lab10-02-grpc-advanced/) | Streaming, REST-to-gRPC gateway | ✅ |

#### Group 11: [Messaging](lab11-messaging/)

| # | Lab | Description | Status |
|---|-----|-------------|--------|
| 11-01 | [Message Queue](lab11-messaging/lab11-01-message-queue/) | RabbitMQ — publisher, consumer, exchanges | ✅ |
| 11-02 | [MQTT](lab11-messaging/lab11-02-mqtt/) | Mosquitto — topics, QoS, wildcards | ✅ |

---

## Learning Path

```
Part 1: REST API
  Group 01  →  Fundamentals (Hello API → CRUD → File handling)
  Group 02  →  Design Conventions (Naming, Params, Validation, Errors, Docs)
  Group 03  →  Security (Auth, Rate Limiting, Sensitive Data, API Keys, Gateway)
  Group 04  →  Versioning (URL, Header, Query, Deprecation, Observability)
  Group 05  →  Advanced Search (POST /search, Saved Search, Security)
  Group 06  →  Performance (Caching, Circuit Breaker)
  Group 07  →  Observability (Logging, Metrics, Tracing)
        ↓
Part 2: Beyond REST
  Group 08  →  GraphQL
  Group 09  →  Real-Time (Webhook, WebSocket)
  Group 10  →  gRPC (Basics, Streaming)
  Group 11  →  Messaging (RabbitMQ, MQTT)
```

## Lab Structure

```
labXX-group-name/
├── README.md                    # Group overview and lab list
├── labXX-YY-topic-name/
│   ├── README.md                # Instructions and exercises
│   ├── docker-compose.yml       # Run everything with one command
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   ├── main.go
│   └── ...
└── labXX-YY-another-topic/
    └── ...
```

Multi-language labs (e.g., Group 04) have `golang/` and `dotnet/` subdirectories inside each sub-lab.

## Services & Ports Quick Reference

| Group | Labs | Services | Ports |
|-------|------|----------|-------|
| 01 | 01-01 to 01-03 | Go API | `8080` |
| 01 | 01-04 | Go API, PostgreSQL | `8080`, `5432` |
| 01 | 01-05 | Go API, PostgreSQL, MinIO | `8080`, `5432`, `9000`, `9001` |
| 02 | 02-02 to 02-04 | Go API | `8080` |
| 02 | 02-07 | Go API, PostgreSQL | `8080`, `5432` |
| 02 | 02-08 | Go API, PostgreSQL, Swagger UI | `8080`, `5432`, `8081` |
| 03 | 03-01, 03-02 | Go API, PostgreSQL | `8080`, `5432` |
| 04 | 04-01 to 04-07 | Go/.NET API, PostgreSQL | `8080`, `5432` |
| 04 | 04-08 | Go/.NET API, PostgreSQL, Prometheus, Grafana | `8080`, `5432`, `9090`, `3000` |
| 08 | 08-01 | Go API (GraphQL Playground), PostgreSQL | `8080`, `5432` |
| 09 | 09-01 | Webhook Sender, Receiver, PostgreSQL | `8080`, `9090`, `5432` |
| 09 | 09-02 | Go API (Chat Web UI) | `8080` |
| 10 | 10-01 | gRPC Server, Client, gRPCUI | `50051`, `8080` |
| 10 | 10-02 | gRPC Server, Client, Gateway, gRPCUI | `50051`, `8080`, `8081` |
| 11 | 11-01 | Publisher, Consumer, RabbitMQ | `8080`, `5672`, `15672` |
| 11 | 11-02 | Publisher, Subscriber, Mosquitto | `1883` |

## Troubleshooting

- **Port conflict**: Stop the previous lab with `docker compose down` before starting a new one.
- **Build cache**: Rebuild with `docker compose up --build --force-recreate`.
- **Database issues**: Reset with `docker compose down -v` to start fresh.
