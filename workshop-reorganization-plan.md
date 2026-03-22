# Workshop Reorganization Plan

> **Status:** PENDING — Not yet started
>
> **Goal:** Restructure all labs into `labXX-YY-topic-name` format where XX = group number, YY = sub-lab number. All REST API groups finish before non-REST technologies. Related topics are grouped together.
>
> **Instructions for AI agents / developers:** This is a major restructuring. Execute Phase 0 (renames) first, then update all docs, then implement new labs phase by phase. Update checkboxes as you go.
>
> **Supersedes:** `lab21-search-api-plan.md` and `workshop-gap-analysis-plan.md` — those plans are merged into this one with updated lab numbers.

---

## New Structure Overview

```
Part 1: REST API (Groups 01–07)
  Group 01: REST Fundamentals ............... learn to build APIs
  Group 02: REST Design Conventions ......... learn to design APIs well
  Group 03: API Security .................... protect your APIs
  Group 04: API Versioning .................. evolve your APIs
  Group 05: Advanced Search ................. complex query patterns
  Group 06: Performance & Resilience ........ make APIs fast and reliable
  Group 07: Observability ................... monitor and debug APIs

Part 2: Beyond REST (Groups 08–11)
  Group 08: GraphQL ......................... flexible queries
  Group 09: Real-Time APIs .................. push-based communication
  Group 10: gRPC ............................ high-performance RPC
  Group 11: Messaging ....................... async communication
```

---

## Complete Lab Mapping (Old → New)

### Group 01: REST Fundamentals

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab01-01 | `lab01-01-hello-api` | lab01-hello-api | ✅ Rename only |
| lab01-02 | `lab01-02-json-response` | lab02-json-response | ✅ Rename only |
| lab01-03 | `lab01-03-crud-in-memory` | lab04-crud-in-memory | ✅ Rename only |
| lab01-04 | `lab01-04-crud-with-database` | lab05-crud-with-database | ✅ Rename only |
| lab01-05 | `lab01-05-file-upload-download` | lab13-file-upload-download | ✅ Rename only |

### Group 02: REST Design Conventions

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab02-01 | `lab02-01-resource-naming` | — | ❌ NEW |
| lab02-02 | `lab02-02-path-and-query-parameters` | lab03-path-parameter | ⚠️ Rename + expand with design decision rules |
| lab02-03 | `lab02-03-request-validation` | lab06-request-validation | ✅ Rename only |
| lab02-04 | `lab02-04-error-handling` | lab07-error-handling | ✅ Rename only |
| lab02-05 | `lab02-05-response-envelope-and-status-codes` | — | ❌ NEW (includes "empty search 200 vs 404" topic) |
| lab02-06 | `lab02-06-data-formats` | — | ❌ NEW |
| lab02-07 | `lab02-07-pagination-and-filtering` | lab09-pagination-and-filtering | ✅ Rename only |
| lab02-08 | `lab02-08-swagger-documentation` | lab08-swagger-documentation | ✅ Rename only |

### Group 03: API Security

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab03-01 | `lab03-01-authentication` | lab10-authentication | ✅ Rename only |
| lab03-02 | `lab03-02-rate-limiting-and-cors` | lab12-rate-limiting-and-cors | ✅ Rename only |
| lab03-03 | `lab03-03-sensitive-data` | — | ❌ NEW |
| lab03-04 | `lab03-04-api-key-management` | — | ❌ NEW |
| lab03-05 | `lab03-05-api-gateway` | — | ❌ NEW |

### Group 04: API Versioning (existing lab11 sub-labs)

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab04-01 | `lab04-01-url-path-versioning` | lab11-01-url-path-versioning | ✅ Rename only |
| lab04-02 | `lab04-02-query-parameter-versioning` | lab11-02-query-parameter-versioning | ✅ Rename only |
| lab04-03 | `lab04-03-header-versioning` | lab11-03-header-versioning | ✅ Rename only |
| lab04-04 | `lab04-04-content-negotiation` | lab11-04-content-negotiation | ✅ Rename only |
| lab04-05 | `lab04-05-evolving-api` | lab11-05-evolving-api | ✅ Rename only |
| lab04-06 | `lab04-06-combining-strategies` | lab11-06-combining-strategies | ✅ Rename only |
| lab04-07 | `lab04-07-breaking-changes-and-deprecation` | lab11-07-breaking-changes-and-deprecation | ✅ Rename only |
| lab04-08 | `lab04-08-version-lifecycle-and-observability` | lab11-08-version-lifecycle-and-observability | ✅ Rename only |

### Group 05: Advanced Search (from lab21 plan)

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab05-01 | `lab05-01-post-search` | — | ❌ NEW |
| lab05-02 | `lab05-02-saved-search` | — | ❌ NEW |
| lab05-03 | `lab05-03-search-security` | — | ❌ NEW |
| lab05-04 | `lab05-04-capstone` | — | ❌ NEW |

### Group 06: Performance & Resilience

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab06-01 | `lab06-01-caching-with-redis` | — | ❌ NEW |
| lab06-02 | `lab06-02-circuit-breaker` | — | ❌ NEW |

### Group 07: Observability

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab07-01 | `lab07-01-structured-logging` | — | ❌ NEW |
| lab07-02 | `lab07-02-metrics` | — | ❌ NEW |
| lab07-03 | `lab07-03-distributed-tracing` | — | ❌ NEW |

### Group 08: GraphQL

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab08-01 | `lab08-01-graphql` | lab14-graphql | ✅ Rename only |

### Group 09: Real-Time APIs

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab09-01 | `lab09-01-webhook` | lab15-webhook | ✅ Rename only |
| lab09-02 | `lab09-02-websocket` | lab16-websocket | ✅ Rename only |

### Group 10: gRPC

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab10-01 | `lab10-01-grpc-basics` | lab17-grpc-basics | ✅ Rename only |
| lab10-02 | `lab10-02-grpc-advanced` | lab18-grpc-advanced | ✅ Rename only |

### Group 11: Messaging

| New Number | New Folder Name | Old Lab | Status |
|------------|----------------|---------|--------|
| lab11-01 | `lab11-01-message-queue` | lab19-message-queue | ✅ Rename only |
| lab11-02 | `lab11-02-mqtt` | lab20-mqtt | ✅ Rename only |

---

## Summary Counts

| Category | Count |
|----------|-------|
| Total labs in new structure | 41 |
| Existing labs to rename | 20 |
| Existing labs to rename + expand | 1 (lab02-02: add design decision rules to path params) |
| New labs to implement | 16 |
| New group READMEs to write | 11 |
| New group CLAUDE.md to write | ~6 (groups with new content) |

---

## Implementation Tracker

### Phase 0: Rename existing labs (20 renames)

> **IMPORTANT:** Rename all folders first, then fix all internal references (go.mod module names, README cross-links, docker-compose references, CLAUDE.md). Do this as one atomic step before any new lab work.

#### Group 01 renames
- [ ] 0.1 `mv lab01-hello-api lab01-01-hello-api`
- [ ] 0.2 `mv lab02-json-response lab01-02-json-response`
- [ ] 0.3 `mv lab04-crud-in-memory lab01-03-crud-in-memory`
- [ ] 0.4 `mv lab05-crud-with-database lab01-04-crud-with-database`
- [ ] 0.5 `mv lab13-file-upload-download lab01-05-file-upload-download`

#### Group 02 renames
- [ ] 0.6 `mv lab03-path-parameter lab02-02-path-and-query-parameters`
- [ ] 0.7 `mv lab06-request-validation lab02-03-request-validation`
- [ ] 0.8 `mv lab07-error-handling lab02-04-error-handling`
- [ ] 0.9 `mv lab09-pagination-and-filtering lab02-07-pagination-and-filtering`
- [ ] 0.10 `mv lab08-swagger-documentation lab02-08-swagger-documentation`

#### Group 03 renames
- [ ] 0.11 `mv lab10-authentication lab03-01-authentication`
- [ ] 0.12 `mv lab12-rate-limiting-and-cors lab03-02-rate-limiting-and-cors`

#### Group 04 renames (lab11 → lab04, including sub-labs)
- [ ] 0.13 `mv lab11-api-versioning lab04-api-versioning` (rename parent)
- [ ] 0.14 Rename all sub-lab folders inside: `lab11-01-*` → `lab04-01-*` through `lab11-08-*` → `lab04-08-*`
- [ ] 0.15 Update go.mod module names in all 8 sub-labs (golang/)
- [ ] 0.16 Update lab04 CLAUDE.md and README.md internal references

#### Group 08–11 renames (Beyond REST)
- [ ] 0.17 `mv lab14-graphql lab08-01-graphql`
- [ ] 0.18 `mv lab15-webhook lab09-01-webhook`
- [ ] 0.19 `mv lab16-websocket lab09-02-websocket`
- [ ] 0.20 `mv lab17-grpc-basics lab10-01-grpc-basics`
- [ ] 0.21 `mv lab18-grpc-advanced lab10-02-grpc-advanced`
- [ ] 0.22 `mv lab19-message-queue lab11-01-message-queue`
- [ ] 0.23 `mv lab20-mqtt lab11-02-mqtt`

#### Fix all references after renames
- [ ] 0.24 Update go.mod module names in all renamed labs
- [ ] 0.25 Update all README.md cross-references (links between labs)
- [ ] 0.26 Update root README.md with new structure and lab table
- [ ] 0.27 Update root CLAUDE.md with new naming convention
- [ ] 0.28 Write group-level README.md for groups that need one (01, 02, 03, 04, 08, 09, 10, 11)
- [ ] 0.29 Verify: `docker compose up --build` works for at least one lab per group

### Issues / Notes (Phase 0)
- (none yet)

---

### Phase 1: Group 02 — New Design Convention labs

#### lab02-01: Resource Naming Convention (NEW)
- [ ] 1.1 Create `lab02-01-resource-naming/golang/` (docker-compose, Dockerfile, go.mod, main.go)
- [ ] 1.2 Write `lab02-01-resource-naming/README.md`
- [ ] 1.3 Verify with docker compose

**Scope:** API showing good vs bad naming patterns. Nouns, plural, kebab-case. Nested resources. Discovery endpoint listing all routes with rationale. Exercise: "fix these bad endpoint names."

#### lab02-02: Expand path-and-query-parameters (EXPAND existing)
- [ ] 1.4 Add design decision rules to README (path = identity, query = filter/option, decision tree)
- [ ] 1.5 Add examples of anti-patterns (filters in path, identity in query)
- [ ] 1.6 Add exercise: "given 10 scenarios, decide path vs query"

#### lab02-05: Response Envelope & Status Codes (NEW)
- [ ] 1.7 Create `lab02-05-response-envelope-and-status-codes/golang/` (all files)
- [ ] 1.8 Write `lab02-05-response-envelope-and-status-codes/README.md`
- [ ] 1.9 Verify with docker compose

**Scope:** Unified envelope `{"data", "error", "meta", "pagination"}`. Status code scenarios: empty search = 200, not found = 404, created = 201, deleted = 204. Consistent error body. Exercise: "what status code for these 15 scenarios?"

#### lab02-06: Data Formats (NEW)
- [ ] 1.10 Create `lab02-06-data-formats/golang/` (all files)
- [ ] 1.11 Write `lab02-06-data-formats/README.md`
- [ ] 1.12 Verify with docker compose

**Scope:** ISO 8601 dates, null handling (omitempty vs explicit null), enum design, money/decimal patterns, UUID identifiers.

#### Group 02 root docs
- [ ] 1.13 Write `lab02-rest-design-conventions/CLAUDE.md`
- [ ] 1.14 Write `lab02-rest-design-conventions/README.md` (learning path, priority guide)

### Issues / Notes (Phase 1)
- (none yet)

---

### Phase 2: Group 03 — New Security labs

#### lab03-03: Sensitive Data Handling (NEW)
- [ ] 2.1 Create `lab03-03-sensitive-data/golang/` (all files)
- [ ] 2.2 Write `lab03-03-sensitive-data/README.md`
- [ ] 2.3 Verify with docker compose

**Scope:** Data masking (`card: "****1234"`), field-level security per role, rules (never in URL/logs), middleware that scrubs sensitive fields from logs, PII classification.

#### lab03-04: API Key Management (NEW)
- [ ] 2.4 Create `lab03-04-api-key-management/golang/` (all files)
- [ ] 2.5 Write `lab03-04-api-key-management/README.md`
- [ ] 2.6 Verify with docker compose

**Scope:** Key lifecycle (create → use → rotate → revoke), DB schema with hashed keys, header-based auth, dual-key rotation, rate limiting per key, audit log.

#### lab03-05: API Gateway Patterns (NEW)
- [ ] 2.7 Create `lab03-05-api-gateway/golang/` (all files)
- [ ] 2.8 Write `lab03-05-api-gateway/README.md`
- [ ] 2.9 Verify with docker compose

**Scope:** Simple reverse proxy gateway, centralized auth + rate limiting + logging at gateway, internal vs external gateway, docker-compose with gateway + 2 backend services.

#### Group 03 root docs
- [ ] 2.10 Write `lab03-api-security/CLAUDE.md`
- [ ] 2.11 Write `lab03-api-security/README.md`

### Issues / Notes (Phase 2)
- (none yet)

---

### Phase 3: Group 05 — Advanced Search (from lab21 plan)

#### lab05-01: POST /search with Body (NEW)
- [ ] 3.1 Create `lab05-01-post-search/golang/` (all files)
- [ ] 3.2 Write `lab05-01-post-search/README.md`
- [ ] 3.3 Verify with docker compose

#### lab05-02: Saved Search Pattern (NEW)
- [ ] 3.4 Create `lab05-02-saved-search/golang/` (all files)
- [ ] 3.5 Write `lab05-02-saved-search/README.md`
- [ ] 3.6 Verify with docker compose

#### lab05-03: Search Security (NEW)
- [ ] 3.7 Create `lab05-03-search-security/golang/` (all files)
- [ ] 3.8 Write `lab05-03-search-security/README.md`
- [ ] 3.9 Verify with docker compose

#### lab05-04: Capstone (NEW)
- [ ] 3.10 Create `lab05-04-capstone/golang/` (all files)
- [ ] 3.11 Write `lab05-04-capstone/README.md`
- [ ] 3.12 Verify with docker compose

#### Group 05 root docs
- [ ] 3.13 Write `lab05-advanced-search/CLAUDE.md`
- [ ] 3.14 Write `lab05-advanced-search/README.md`

> **Detailed spec:** See "Sub-Lab Details" section in `lab21-search-api-plan.md` for endpoints, DB schemas, request/response shapes, and code patterns for each sub-lab.

### Issues / Notes (Phase 3)
- (none yet)

---

### Phase 4: Group 06 — Performance & Resilience (NEW)

#### lab06-01: Caching with Redis (NEW)
- [ ] 4.1 Create `lab06-01-caching-with-redis/golang/` (all files)
- [ ] 4.2 Write `lab06-01-caching-with-redis/README.md`
- [ ] 4.3 Verify with docker compose

**Scope:** Cache-aside pattern, write-through, TTL expiry, cache invalidation on update/delete, HTTP cache headers (Cache-Control, ETag, 304), docker-compose with API + PostgreSQL + Redis.

#### lab06-02: Circuit Breaker (NEW)
- [ ] 4.4 Create `lab06-02-circuit-breaker/golang/` (all files)
- [ ] 4.5 Write `lab06-02-circuit-breaker/README.md`
- [ ] 4.6 Verify with docker compose

**Scope:** Closed → Open → Half-Open states, fallback responses, health check design (`/health`, `/health/ready`, `/health/live`), docker-compose with API + flaky downstream service.

#### Group 06 root docs
- [ ] 4.7 Write `lab06-performance-and-resilience/CLAUDE.md`
- [ ] 4.8 Write `lab06-performance-and-resilience/README.md`

### Issues / Notes (Phase 4)
- (none yet)

---

### Phase 5: Group 07 — Observability (NEW)

#### lab07-01: Structured Logging & Correlation ID (NEW)
- [ ] 5.1 Create `lab07-01-structured-logging/golang/` (all files)
- [ ] 5.2 Write `lab07-01-structured-logging/README.md`
- [ ] 5.3 Verify with docker compose

**Scope:** JSON structured logs, X-Request-ID middleware, log levels, what NOT to log (PII, tokens).

#### lab07-02: Metrics with Prometheus & Grafana (NEW)
- [ ] 5.4 Create `lab07-02-metrics/golang/` (all files)
- [ ] 5.5 Write `lab07-02-metrics/README.md`
- [ ] 5.6 Verify with docker compose

**Scope:** RED method, Prometheus counters/histograms/gauges, custom business metrics, Grafana auto-provisioned dashboard, docker-compose with API + Prometheus + Grafana.

#### lab07-03: Distributed Tracing with OpenTelemetry & Jaeger (NEW)
- [ ] 5.7 Create `lab07-03-distributed-tracing/golang/` (all files)
- [ ] 5.8 Write `lab07-03-distributed-tracing/README.md`
- [ ] 5.9 Verify with docker compose

**Scope:** OpenTelemetry SDK, spans with parent/child, traceparent header propagation, Jaeger UI, connecting traceId to structured logs, docker-compose with API + downstream service + Jaeger.

#### Group 07 root docs
- [ ] 5.10 Write `lab07-observability/CLAUDE.md`
- [ ] 5.11 Write `lab07-observability/README.md`

### Issues / Notes (Phase 5)
- (none yet)

---

### Phase 6: Final documentation

- [ ] 6.1 Write new root `README.md` with full restructured lab table, learning path, and ports reference
- [ ] 6.2 Update root `CLAUDE.md` with new naming convention (`labXX-YY-topic-name`)
- [ ] 6.3 Write group-level READMEs for Beyond REST groups (08, 09, 10, 11) if needed
- [ ] 6.4 Final verification: walk through learning path, confirm every lab builds/runs

### Issues / Notes (Phase 6)
- (none yet)

---

## New Folder Structure (Final State)

```
application-programming-interface-design-workshop-2026/
├── CLAUDE.md
├── README.md
│
├── Part 1: REST API
│
├── lab01-rest-fundamentals/
│   ├── README.md
│   ├── lab01-01-hello-api/                          ✅ exists (renamed from lab01)
│   ├── lab01-02-json-response/                      ✅ exists (renamed from lab02)
│   ├── lab01-03-crud-in-memory/                     ✅ exists (renamed from lab04)
│   ├── lab01-04-crud-with-database/                 ✅ exists (renamed from lab05)
│   └── lab01-05-file-upload-download/               ✅ exists (renamed from lab13)
│
├── lab02-rest-design-conventions/
│   ├── README.md
│   ├── CLAUDE.md
│   ├── lab02-01-resource-naming/                    ❌ NEW
│   ├── lab02-02-path-and-query-parameters/          ⚠️ exists, needs expansion (from lab03)
│   ├── lab02-03-request-validation/                 ✅ exists (renamed from lab06)
│   ├── lab02-04-error-handling/                     ✅ exists (renamed from lab07)
│   ├── lab02-05-response-envelope-and-status-codes/ ❌ NEW
│   ├── lab02-06-data-formats/                       ❌ NEW
│   ├── lab02-07-pagination-and-filtering/           ✅ exists (renamed from lab09)
│   └── lab02-08-swagger-documentation/              ✅ exists (renamed from lab08)
│
├── lab03-api-security/
│   ├── README.md
│   ├── CLAUDE.md
│   ├── lab03-01-authentication/                     ✅ exists (renamed from lab10)
│   ├── lab03-02-rate-limiting-and-cors/             ✅ exists (renamed from lab12)
│   ├── lab03-03-sensitive-data/                     ❌ NEW
│   ├── lab03-04-api-key-management/                 ❌ NEW
│   └── lab03-05-api-gateway/                        ❌ NEW
│
├── lab04-api-versioning/
│   ├── README.md
│   ├── CLAUDE.md
│   ├── lab04-01-url-path-versioning/                ✅ exists (renamed from lab11-01)
│   ├── lab04-02-query-parameter-versioning/         ✅ exists (renamed from lab11-02)
│   ├── lab04-03-header-versioning/                  ✅ exists (renamed from lab11-03)
│   ├── lab04-04-content-negotiation/                ✅ exists (renamed from lab11-04)
│   ├── lab04-05-evolving-api/                       ✅ exists (renamed from lab11-05)
│   ├── lab04-06-combining-strategies/               ✅ exists (renamed from lab11-06)
│   ├── lab04-07-breaking-changes-and-deprecation/   ✅ exists (renamed from lab11-07)
│   └── lab04-08-version-lifecycle-and-observability/ ✅ exists (renamed from lab11-08)
│
├── lab05-advanced-search/
│   ├── README.md
│   ├── CLAUDE.md
│   ├── lab05-01-post-search/                        ❌ NEW
│   ├── lab05-02-saved-search/                       ❌ NEW
│   ├── lab05-03-search-security/                    ❌ NEW
│   └── lab05-04-capstone/                           ❌ NEW
│
├── lab06-performance-and-resilience/
│   ├── README.md
│   ├── CLAUDE.md
│   ├── lab06-01-caching-with-redis/                 ❌ NEW
│   └── lab06-02-circuit-breaker/                    ❌ NEW
│
├── lab07-observability/
│   ├── README.md
│   ├── CLAUDE.md
│   ├── lab07-01-structured-logging/                 ❌ NEW
│   ├── lab07-02-metrics/                            ❌ NEW
│   └── lab07-03-distributed-tracing/                ❌ NEW
│
├── Part 2: Beyond REST
│
├── lab08-graphql/
│   ├── README.md
│   └── lab08-01-graphql/                            ✅ exists (renamed from lab14)
│
├── lab09-real-time-apis/
│   ├── README.md
│   ├── lab09-01-webhook/                            ✅ exists (renamed from lab15)
│   └── lab09-02-websocket/                          ✅ exists (renamed from lab16)
│
├── lab10-grpc/
│   ├── README.md
│   ├── lab10-01-grpc-basics/                        ✅ exists (renamed from lab17)
│   └── lab10-02-grpc-advanced/                      ✅ exists (renamed from lab18)
│
└── lab11-messaging/
    ├── README.md
    ├── lab11-01-message-queue/                      ✅ exists (renamed from lab19)
    └── lab11-02-mqtt/                               ✅ exists (renamed from lab20)
```

---

## Implementation Priority

| Priority | Phase | What | Effort |
|----------|-------|------|--------|
| 🔴 Do first | Phase 0 | Rename all 20 existing labs + fix references | Large but mechanical |
| 🔴 High | Phase 1 | Group 02: Design Conventions (3 new + 1 expand) | Medium |
| 🔴 High | Phase 2 | Group 03: Security Advanced (3 new) | Medium |
| 🔴 High | Phase 3 | Group 05: Advanced Search (4 new) | Large |
| 🟡 Medium | Phase 5 | Group 07: Observability (3 new) | Medium |
| 🟡 Medium | Phase 4 | Group 06: Performance (2 new) | Medium |
| 🟢 Final | Phase 6 | All documentation updates | Small |

---

## Critical Existing Files to Reference

```
# Patterns to reuse
lab01-04 (was lab05) /main.go              — base CRUD + PostgreSQL
lab02-03 (was lab06) /main.go              — go-playground/validator
lab02-04 (was lab07) /main.go              — error types, error middleware
lab02-07 (was lab09) /main.go              — dynamic WHERE, PaginatedResponse
lab03-01 (was lab10) /auth.go              — JWT AuthMiddleware, bcrypt
lab03-02 (was lab12) /main.go              — token bucket rate limiter
lab04 (was lab11) /                        — sub-lab structure pattern, CLAUDE.md format
lab04-08 (was lab11-08) /golang/           — Prometheus, Grafana, structured logs, loadtest

# Source documents
new-api-design-class.md                    — 2-day course outline driving requirements
lab21-search-api-plan.md                   — detailed spec for Group 05 sub-labs (endpoints, schemas, patterns)
workshop-gap-analysis-plan.md              — gap analysis with detailed scope per new lab
```
