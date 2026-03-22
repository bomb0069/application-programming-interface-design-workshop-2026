# Workshop Gap Analysis & Implementation Plan

> **Status:** PENDING — Not yet started
>
> **Source:** Cross-reference between `new-api-design-class.md` (2-day course outline) and existing labs (lab01–lab20 + lab11 sub-labs + lab21 plan).
>
> **Instructions for AI agents / developers:** Read this plan before working on this project. Check the tracker below to know what's done and what's next. Update the checkboxes (`- [ ]` → `- [x]`) and status as you complete each item.

---

## Gap Analysis Summary

### Coverage Matrix: Course Topics vs Existing Labs

| # | Course Topic | Existing Lab | Coverage | Action Needed |
|---|---|---|---|---|
| 1.1 | REST Intro (REST vs GraphQL vs gRPC) | Lab 01, 14, 17 | ✅ Full | None — covered across multiple labs |
| 1.2 | Resource Naming Convention | None | ❌ Missing | **New lab needed** |
| 1.3a | HTTP Methods (GET/POST/PUT/PATCH/DELETE) | Lab 04 | ✅ Full | None — lab04 CRUD covers all methods |
| 1.3b | Path Variable vs Query Parameter decisions | Lab 03, 09 | ⚠️ Partial | Exists but no **decision rules** or **design guidance** |
| 1.3c | Complex Search (10-20 params, POST /search) | None (planned lab21) | ❌ Planned | **Lab 21 plan exists** — implement it |
| 1.3d | Pagination, Filtering, Sorting | Lab 09 | ✅ Full | None |
| 1.4a | Status Code Standards (200 vs 404 for empty) | Lab 07 | ⚠️ Partial | Covers error codes but NOT the "empty search = 200 or 404?" decision |
| 1.4b | Standard Response Envelope | Lab 07, 09 | ⚠️ Partial | Lab 09 has `{data, metadata}` but no **unified envelope** across labs |
| 1.4c | Data Formats (ISO 8601, nulls, enums) | None | ❌ Missing | **New lab needed** |
| 1.5 | API Versioning | Lab 11 (8 sub-labs) | ✅ Comprehensive | None — excellent coverage |
| 1.6 | API Documentation (OpenAPI/Swagger) | Lab 08 | ✅ Full | None — covers both Design-First and Code-First |
| 2.1 | Auth (OAuth2 / JWT) | Lab 10 | ⚠️ Partial | JWT done, **OAuth2 missing** |
| 2.2 | API Key Management | Lab 10 (exercise only) | ❌ Missing | **New lab needed** — lifecycle, rotation, revocation |
| 2.3 | Sensitive Data Handling | None | ❌ Missing | **New lab needed** — masking, field-level security |
| 2.4 | Rate Limiting | Lab 12 | ✅ Full | None — token bucket, headers, 429 all covered |
| 2.5 | Defense Patterns (OWASP, input sanitization) | Lab 06 | ⚠️ Partial | Validation done, **OWASP overview missing** |
| 2.6 | API Gateway | None | ❌ Missing | **New lab needed** |
| 3.1 | Architecture (Monolith → Microservices) | None | ❌ Missing | Conceptual — may not need a lab |
| 3.2a | Caching (Cache-aside, Redis) | None | ❌ Missing | **New lab needed** |
| 3.2b | DB Design (normalization, connection pool) | Lab 05 | ⚠️ Partial | Basic SQL only, no design patterns |
| 3.3 | Resilience (Circuit Breaker, Health Check) | None | ❌ Missing | **New lab needed** |
| 4.1 | Centralized Logging (structured, correlation ID) | Lab 11-08 (partial) | ⚠️ Partial | Lab 11-08 has structured logs but **no correlation ID**, not a dedicated lab |
| 4.2 | Metrics (Prometheus + Grafana) | Lab 11-08 | ⚠️ Partial | Exists within versioning context, not standalone |
| 4.3 | Distributed Tracing (OpenTelemetry, Jaeger) | None | ❌ Missing | **New lab needed** |

---

## Action Categories

### A. Already covered — no action needed
- REST Intro, HTTP Methods, Pagination/Filtering/Sorting, API Versioning, Swagger/OpenAPI, Rate Limiting

### B. Exists but needs expansion (enhance existing labs or add guidance docs)
- Path Variable vs Query Parameter — add decision rules
- Status Codes — add empty search result guidance
- Response Envelope — standardize across labs
- Auth — add OAuth2 / API Key management
- Defense — add OWASP overview
- Logging/Metrics — already in lab11-08, needs standalone lab

### C. Completely missing — new labs needed
- Resource Naming Convention
- Sensitive Data Handling (masking, field-level security)
- Data Formats (ISO 8601, nulls, enums)
- API Gateway patterns
- Caching Strategies (Redis)
- Circuit Breaker / Resilience
- Distributed Tracing (OpenTelemetry + Jaeger)
- API Key Lifecycle Management

---

## Implementation Tracker

### Lab 22: REST API Design Conventions (NEW)
> Covers: Resource Naming, Path vs Query decisions, Response Envelope, Status Code decisions, Data Formats.
> These are "design guidance" topics that the course marks as ⭐ high priority. Grouped into one lab because they're conceptual patterns, not standalone services.

**Sub-lab structure (like lab11):**

#### Sub-lab 22-01: Resource Naming Convention
- [ ] 22-01.1 Create `lab22-rest-api-design-conventions/lab22-01-resource-naming/golang/` (all files)
- [ ] 22-01.2 Write `lab22-01-resource-naming/README.md`
- [ ] 22-01.3 Verify with docker compose

**Scope:**
- API showing good vs bad naming: `/api/v1/users` vs `/api/v1/getUser`, plural vs singular, kebab-case
- Nested resources: `/users/{id}/orders` vs flat `/orders?userId=123`
- Standard prefixes: `/api/`, `/internal/`, `/health`
- Discovery endpoint listing all routes with naming rationale
- Exercise: "fix these bad endpoint names"

#### Sub-lab 22-02: Path Variable vs Query Parameter
- [ ] 22-02.1 Create `lab22-rest-api-design-conventions/lab22-02-path-vs-query/golang/` (all files)
- [ ] 22-02.2 Write `lab22-02-path-vs-query/README.md`
- [ ] 22-02.3 Verify with docker compose

**Scope:**
- Decision tree: Path = identity (`/users/{id}`), Query = filter/option (`?status=active`)
- Same API demonstrating both patterns with explanation of when to use which
- Anti-patterns: putting filters in path, putting identity in query
- Exercise: given 10 scenarios, decide path vs query

#### Sub-lab 22-03: Response Envelope & Status Codes
- [ ] 22-03.1 Create `lab22-rest-api-design-conventions/lab22-03-response-envelope/golang/` (all files)
- [ ] 22-03.2 Write `lab22-03-response-envelope/README.md`
- [ ] 22-03.3 Verify with docker compose

**Scope:**
- Unified envelope: `{"data": ..., "error": null, "meta": {}, "pagination": {}}`
- Status code decisions with scenarios:
  - `GET /users` returns empty list → 200 with `{"data": [], "pagination": {...}}`
  - `GET /users/{id}` not found → 404 with `{"data": null, "error": {...}}`
  - `POST /users` created → 201 with `{"data": {...}}`
  - `DELETE /users/{id}` → 204 no content
- Consistent error body: `{"error": {"code": "NOT_FOUND", "message": "..."}}`
- Exercise: "what status code for these 15 scenarios?"

#### Sub-lab 22-04: Data Formats & Standards
- [ ] 22-04.1 Create `lab22-rest-api-design-conventions/lab22-04-data-formats/golang/` (all files)
- [ ] 22-04.2 Write `lab22-04-data-formats/README.md`
- [ ] 22-04.3 Verify with docker compose

**Scope:**
- ISO 8601 dates (`2026-03-22T10:00:00Z`) — parsing and returning
- Null handling: `omitempty` vs explicit null vs absent field
- Enum design: string enums, unknown value handling, extensibility contract
- Money/decimal: why not float64 for money, string or integer cents pattern
- UUID as resource identifiers

#### Root docs for Lab 22
- [ ] 22-00.1 Write `lab22-rest-api-design-conventions/CLAUDE.md` (knowledge base)
- [ ] 22-00.2 Write `lab22-rest-api-design-conventions/README.md` (learning path)
- [ ] 22-00.3 Update main workshop `README.md`

### Issues / Notes
<!-- Record any issues, blockers, or decisions made during implementation here -->
- (none yet)

---

### Lab 23: API Security Advanced (NEW)
> Covers: Sensitive Data Handling, API Key Management, API Gateway patterns, OWASP basics.
> Extends lab10 (JWT auth) and lab12 (rate limiting) into production security patterns.

#### Sub-lab 23-01: Sensitive Data Handling
- [ ] 23-01.1 Create `lab23-api-security-advanced/lab23-01-sensitive-data/golang/` (all files)
- [ ] 23-01.2 Write `lab23-01-sensitive-data/README.md`
- [ ] 23-01.3 Verify with docker compose

**Scope:**
- Data masking in responses: `card: "****1234"`, `email: "j***@example.com"`
- Field-level security: different response fields per role (public/internal/admin)
- Rules: never in URL, never in logs, never over-expose
- Middleware that scrubs sensitive fields from logs
- HTTPS/TLS basics (conceptual + docker-compose with self-signed cert)
- PII classification exercise

#### Sub-lab 23-02: API Key Management
- [ ] 23-02.1 Create `lab23-api-security-advanced/lab23-02-api-key-management/golang/` (all files)
- [ ] 23-02.2 Write `lab23-02-api-key-management/README.md`
- [ ] 23-02.3 Verify with docker compose

**Scope:**
- API Key lifecycle: create → use → rotate → revoke
- DB schema: `api_keys` table with hashed key, client_name, scopes, created_at, expires_at, revoked_at
- Key passed via `Authorization: ApiKey xxx` header (never in URL)
- Dual-key rotation: new key active while old key still valid for grace period
- Rate limiting per API key
- Audit: log which key accessed what

#### Sub-lab 23-03: API Gateway Patterns
- [ ] 23-03.1 Create `lab23-api-security-advanced/lab23-03-api-gateway/golang/` (all files)
- [ ] 23-03.2 Write `lab23-03-api-gateway/README.md`
- [ ] 23-03.3 Verify with docker compose

**Scope:**
- Simple reverse proxy gateway (Go or nginx) in front of API service
- Centralized concerns: auth validation, rate limiting, request logging, CORS — all at gateway
- Internal vs external gateway: different auth for internal services vs public consumers
- docker-compose with gateway + 2 backend services
- Request routing and path rewriting

#### Root docs for Lab 23
- [ ] 23-00.1 Write `lab23-api-security-advanced/CLAUDE.md` (knowledge base)
- [ ] 23-00.2 Write `lab23-api-security-advanced/README.md` (learning path)
- [ ] 23-00.3 Update main workshop `README.md`

### Issues / Notes
- (none yet)

---

### Lab 24: Performance & Resilience (NEW)
> Covers: Caching (Redis), Circuit Breaker, Health Check design.

#### Sub-lab 24-01: Caching Strategies with Redis
- [ ] 24-01.1 Create `lab24-performance-and-resilience/lab24-01-caching-with-redis/golang/` (all files)
- [ ] 24-01.2 Write `lab24-01-caching-with-redis/README.md`
- [ ] 24-01.3 Verify with docker compose

**Scope:**
- Cache-aside pattern: check Redis → miss → query DB → store in Redis → return
- Write-through: write to DB + Redis simultaneously
- Cache invalidation on update/delete
- TTL-based expiry
- docker-compose with API + PostgreSQL + Redis
- Cache hit/miss metrics
- HTTP cache headers: Cache-Control, ETag, If-None-Match (304 Not Modified)

#### Sub-lab 24-02: Circuit Breaker & Resilience
- [ ] 24-02.1 Create `lab24-performance-and-resilience/lab24-02-circuit-breaker/golang/` (all files)
- [ ] 24-02.2 Write `lab24-02-circuit-breaker/README.md`
- [ ] 24-02.3 Verify with docker compose

**Scope:**
- Circuit breaker states: Closed → Open → Half-Open
- API calls an unreliable downstream service; circuit breaker prevents cascade failure
- Fallback responses when circuit is open
- Health check design: `/health` (simple) vs `/health/ready` (deep) vs `/health/live` (liveness)
- What to expose in health check vs what to hide (no DB credentials, no internal IPs)
- docker-compose with API + flaky downstream service (simulated failures)

#### Root docs for Lab 24
- [ ] 24-00.1 Write `lab24-performance-and-resilience/CLAUDE.md` (knowledge base)
- [ ] 24-00.2 Write `lab24-performance-and-resilience/README.md` (learning path)
- [ ] 24-00.3 Update main workshop `README.md`

### Issues / Notes
- (none yet)

---

### Lab 25: Observability (NEW)
> Covers: Structured Logging, Correlation ID, Distributed Tracing (OpenTelemetry + Jaeger), Metrics dashboard.
> Extends lab11-08's partial observability into a dedicated, comprehensive lab.

#### Sub-lab 25-01: Structured Logging & Correlation ID
- [ ] 25-01.1 Create `lab25-observability/lab25-01-structured-logging/golang/` (all files)
- [ ] 25-01.2 Write `lab25-01-structured-logging/README.md`
- [ ] 25-01.3 Verify with docker compose

**Scope:**
- Structured JSON logs: timestamp, level, message, requestId, userId, method, path, status, durationMs
- Correlation ID middleware: generate UUID per request, propagate via `X-Request-ID` header
- Log levels: DEBUG, INFO, WARN, ERROR — when to use each
- What NOT to log: passwords, tokens, PII, credit cards
- Log aggregation concept (mention ELK/Loki, but use simple stdout for workshop)

#### Sub-lab 25-02: Metrics with Prometheus & Grafana
- [ ] 25-02.1 Create `lab25-observability/lab25-02-metrics/golang/` (all files)
- [ ] 25-02.2 Write `lab25-02-metrics/README.md`
- [ ] 25-02.3 Verify with docker compose

**Scope:**
- RED method: Rate (requests/sec), Errors (error rate), Duration (latency)
- Prometheus client: counters, histograms, gauges
- Custom business metrics (e.g., orders_created_total, search_queries_total)
- Grafana dashboard with auto-provisioning (reuse pattern from lab11-08)
- Alerting concepts: when error rate > 5%, when p99 latency > 2s
- docker-compose: API + Prometheus + Grafana

#### Sub-lab 25-03: Distributed Tracing with OpenTelemetry & Jaeger
- [ ] 25-03.1 Create `lab25-observability/lab25-03-distributed-tracing/golang/` (all files)
- [ ] 25-03.2 Write `lab25-03-distributed-tracing/README.md`
- [ ] 25-03.3 Verify with docker compose

**Scope:**
- OpenTelemetry SDK: instrument HTTP server + DB calls
- Spans: parent/child relationships, attributes, events
- Trace propagation: `traceparent` header across services
- Jaeger UI for trace visualization
- Connecting logs + traces: add traceId to structured logs
- docker-compose: API + downstream service + Jaeger + (optional) Prometheus
- Exercise: "find the slow query by looking at the trace"

#### Root docs for Lab 25
- [ ] 25-00.1 Write `lab25-observability/CLAUDE.md` (knowledge base)
- [ ] 25-00.2 Write `lab25-observability/README.md` (learning path)
- [ ] 25-00.3 Update main workshop `README.md`

### Issues / Notes
- (none yet)

---

## Implementation Priority

Based on the course pain points (🔴 High priority from `new-api-design-class.md`):

| Priority | Lab | Course Pain Point Addressed |
|----------|-----|----------------------------|
| 🔴 1st | Lab 21 (Search API) | Search 10-20 params — URL Limit (plan already exists) |
| 🔴 2nd | Lab 22 (Design Conventions) | Naming Convention, Path vs Query, Response Codes |
| 🔴 3rd | Lab 23 (Security Advanced) | Sensitive Data, API Key Management |
| 🟡 4th | Lab 25 (Observability) | Logging, Metrics, Tracing |
| 🟡 5th | Lab 24 (Performance) | Caching, Circuit Breaker |

---

## Main Workshop README Updates (when labs are implemented)

Add to `README.md`:

```markdown
### Part 4: Advanced REST Deep Dives

Go deeper into specific REST API design challenges. These labs build on patterns from Part 1 and Part 2.

| # | Lab | Description |
|---|-----|-------------|
| 21 | [Advanced Search API](lab21-search-api/) | POST /search with complex filters, saved search pattern, security, multi-tenant |
| 22 | [REST API Design Conventions](lab22-rest-api-design-conventions/) | Resource naming, path vs query decisions, response envelope, data formats |
| 23 | [API Security Advanced](lab23-api-security-advanced/) | Sensitive data masking, API key lifecycle, API gateway patterns |
| 24 | [Performance & Resilience](lab24-performance-and-resilience/) | Redis caching strategies, circuit breaker, health check design |
| 25 | [Observability](lab25-observability/) | Structured logging, Prometheus metrics, distributed tracing with OpenTelemetry |
```

Update Learning Path:
```
Part 4: Advanced REST Deep Dives
  Lab 21     →  Search API (builds on Lab 09 + 10)
  Lab 22     →  Design Conventions (foundational, can start anytime)
  Lab 23     →  Security Advanced (builds on Lab 10 + 12)
  Lab 24     →  Performance (builds on Lab 05)
  Lab 25     →  Observability (builds on Lab 11-08)
```

---

## Critical Existing Files to Reference

```
# Existing patterns to reuse
lab03-path-parameter/main.go               — path parameter extraction
lab05-crud-with-database/main.go           — base CRUD + PostgreSQL pattern
lab06-request-validation/main.go           — go-playground/validator
lab07-error-handling/main.go               — error types, error middleware
lab09-pagination-and-filtering/main.go     — dynamic WHERE, PaginatedResponse
lab10-authentication/auth.go               — JWT AuthMiddleware, bcrypt
lab11-api-versioning/                      — sub-lab structure pattern
lab11-08-.../golang/                       — Prometheus/Grafana/structured logs pattern
lab12-rate-limiting-and-cors/main.go       — token bucket rate limiter

# Plans (check status before implementing)
lab21-search-api-plan.md                   — Lab 21 detailed plan (PENDING)

# Source document
new-api-design-class.md                    — 2-day course outline driving these requirements
```
