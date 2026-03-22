# Implementation Plan: Lab 21 — Advanced Search API

> **Status:** PENDING — Not yet started
>
> **Instructions for AI agents / developers:** Read this plan before working on this project. Check the tracker below to know what's done and what's next. Update the checkboxes and status as you complete each item.

---

## Implementation Tracker

### Phase 1: Foundation
- [ ] 1.1 Create full directory structure (`lab21-search-api/` with all sub-lab folders)
- [ ] 1.2 Write `lab21-search-api/CLAUDE.md` (knowledge base from user's learning sessions)

### Phase 2: Sub-lab 01 — POST /search
- [ ] 2.1 Create `lab21-01-post-search/golang/docker-compose.yml`
- [ ] 2.2 Create `lab21-01-post-search/golang/Dockerfile`
- [ ] 2.3 Create `lab21-01-post-search/golang/go.mod`
- [ ] 2.4 Create `lab21-01-post-search/golang/main.go` (DB setup, seed 30+ products, routes)
- [ ] 2.5 Create `lab21-01-post-search/golang/search_handler.go` (filter operators, buildWhereClause)
- [ ] 2.6 Write `lab21-01-post-search/README.md`
- [ ] 2.7 Verify: `docker compose up --build` and test curl examples

### Phase 3: Sub-lab 02 — Saved Search
- [ ] 3.1 Create `lab21-02-saved-search/golang/docker-compose.yml`
- [ ] 3.2 Create `lab21-02-saved-search/golang/Dockerfile`
- [ ] 3.3 Create `lab21-02-saved-search/golang/go.mod`
- [ ] 3.4 Create `lab21-02-saved-search/golang/main.go` (DB setup, search_sessions + saved_searches tables, cleanup goroutine)
- [ ] 3.5 Create `lab21-02-saved-search/golang/search_handler.go` (reuse from sub-lab 01)
- [ ] 3.6 Create `lab21-02-saved-search/golang/saved_search_handler.go` (session CRUD, 410 Gone)
- [ ] 3.7 Write `lab21-02-saved-search/README.md`
- [ ] 3.8 Verify: test searchId flow, expiry → 410 Gone, named saved searches

### Phase 4: Sub-lab 03 — Search Security
- [ ] 4.1 Create `lab21-03-search-security/golang/docker-compose.yml`
- [ ] 4.2 Create `lab21-03-search-security/golang/Dockerfile`
- [ ] 4.3 Create `lab21-03-search-security/golang/go.mod`
- [ ] 4.4 Create `lab21-03-search-security/golang/main.go` (DB setup, users + audit_log tables, routes)
- [ ] 4.5 Create `lab21-03-search-security/golang/auth.go` (JWT middleware, register, login — reuse lab10 pattern)
- [ ] 4.6 Create `lab21-03-search-security/golang/search_handler.go` (forced scope, field projection per role)
- [ ] 4.7 Create `lab21-03-search-security/golang/middleware.go` (rate limiting, audit logging, query timeout)
- [ ] 4.8 Write `lab21-03-search-security/README.md`
- [ ] 4.9 Verify: auth flow, scoped results, rate limit 429, field projection by role

### Phase 5: Sub-lab 04 — Capstone
- [ ] 5.1 Create `lab21-04-capstone/golang/docker-compose.yml`
- [ ] 5.2 Create `lab21-04-capstone/golang/Dockerfile`
- [ ] 5.3 Create `lab21-04-capstone/golang/go.mod`
- [ ] 5.4 Create `lab21-04-capstone/golang/main.go` (all tables combined)
- [ ] 5.5 Create `lab21-04-capstone/golang/auth.go`
- [ ] 5.6 Create `lab21-04-capstone/golang/search_handler.go`
- [ ] 5.7 Create `lab21-04-capstone/golang/saved_search_handler.go` (with user ownership)
- [ ] 5.8 Create `lab21-04-capstone/golang/middleware.go`
- [ ] 5.9 Create `lab21-04-capstone/golang/audit.go` (admin audit log viewer)
- [ ] 5.10 Write `lab21-04-capstone/README.md`
- [ ] 5.11 Verify: multi-tenant isolation, saved search ownership, admin audit log

### Phase 6: Root Documentation
- [ ] 6.1 Write `lab21-search-api/README.md` (learning path, priority guide, project structure)
- [ ] 6.2 Update main workshop `README.md` (add Part 4 section, learning path diagram, ports table)

### Issues / Notes
<!-- Record any issues, blockers, or decisions made during implementation here -->
- (none yet)

---

## Context

The user has been learning API search design in another Claude thread, covering 3 topics: POST /search with complex filters, Saved Search (two-phase pattern), and Search API Security. These patterns are not covered in the current workshop (lab09 only has basic GET query param filtering). This plan adds a new sub-lab series teaching advanced search progressively, following the same structure as lab11.

**Placement:** `lab21-search-api` — avoids renumbering existing labs. Main README gets a new "Part 4: Advanced REST Deep Dives" section. The lab builds on patterns from lab09 (pagination/filtering) and lab10 (JWT auth).

**Language:** Go only (no .NET for now — can be added later like lab11).

---

## Directory Structure

```
lab21-search-api/
├── CLAUDE.md                              # Knowledge base from user's learning sessions
├── README.md                              # Root README with learning path + priority guide
├── lab21-01-post-search/
│   ├── README.md
│   └── golang/
│       ├── docker-compose.yml
│       ├── Dockerfile
│       ├── go.mod
│       ├── main.go
│       └── search_handler.go
├── lab21-02-saved-search/
│   ├── README.md
│   └── golang/
│       ├── docker-compose.yml
│       ├── Dockerfile
│       ├── go.mod
│       ├── main.go
│       ├── search_handler.go
│       └── saved_search_handler.go
├── lab21-03-search-security/
│   ├── README.md
│   └── golang/
│       ├── docker-compose.yml
│       ├── Dockerfile
│       ├── go.mod
│       ├── main.go
│       ├── auth.go
│       ├── search_handler.go
│       └── middleware.go
└── lab21-04-capstone/
    ├── README.md
    └── golang/
        ├── docker-compose.yml
        ├── Dockerfile
        ├── go.mod
        ├── main.go
        ├── auth.go
        ├── search_handler.go
        ├── saved_search_handler.go
        ├── middleware.go
        └── audit.go
```

---

## Sub-Lab Details

### Sub-lab 01: POST /search with Body (~300 lines Go)

**What it teaches:** Why POST for complex search, structured filter operators, dynamic WHERE clause building.

**Database:** Products table (id, name, price, category, brand, in_stock, created_at) seeded with 30+ products.

**Endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/search` | Complex search with filters, sort, pagination |
| GET | `/products/{id}` | Get single product |
| GET | `/health` | Health check |

**Request body:**
```json
{
  "filters": [
    {"field": "category", "operator": "eq", "value": "electronics"},
    {"field": "price", "operator": "gte", "value": 50},
    {"field": "price", "operator": "lte", "value": 500},
    {"field": "brand", "operator": "in", "value": ["Apple", "Samsung"]},
    {"field": "name", "operator": "like", "value": "Pro"}
  ],
  "sort": {"field": "price", "order": "desc"},
  "pagination": {"page": 1, "limit": 10}
}
```

**Key patterns:**
- Supported operators: `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `in`, `like`
- `buildWhereClause(filters)` — extends lab09's dynamic WHERE with `argIdx` counter, parameterized queries
- Whitelist of allowed filter fields and sortable fields (map lookups)
- `like` operator blocks leading wildcards on raw value
- Response shape matches lab09: `{"data": [...], "metadata": {"current_page", "page_size", "total_items", "total_pages"}}`

**Dependencies:** `chi/v5`, `lib/pq`

---

### Sub-lab 02: Saved Search Pattern (~350 lines Go)

**What it teaches:** Two-phase search (save criteria → execute with GET), TTL-based sessions, 410 Gone handling, named saved searches.

**Additional tables:**
```sql
CREATE TABLE IF NOT EXISTS search_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    criteria JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS saved_searches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    criteria JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/search` | Direct search (from sub-lab 01) |
| POST | `/searches` | Save criteria → returns searchId + TTL |
| GET | `/products?searchId={id}&page=1&limit=10` | Execute saved search with pagination |
| POST | `/saved-searches` | Persist named search permanently |
| GET | `/saved-searches` | List saved searches |
| GET | `/saved-searches/{id}` | Get saved search by ID |
| DELETE | `/saved-searches/{id}` | Delete saved search |

**Key patterns:**
- `POST /searches` stores criteria as JSONB, sets `expires_at = NOW() + 30min`, returns `{"searchId": "uuid", "expiresIn": 1800, "resultsUrl": "/products?searchId=uuid"}`
- `GET /products?searchId=xxx` — if expired → **410 Gone** with `{"error": "SEARCH_EXPIRED", "message": "...", "action": "Re-submit via POST /searches"}`
- Background goroutine cleans expired sessions every 5 minutes
- Saved searches persist without TTL, support name and list/delete

**Dependencies:** `chi/v5`, `lib/pq`

---

### Sub-lab 03: Search API Security (~400 lines Go)

**What it teaches:** 6 security concerns — injection prevention, BOLA/access control, field projection by role, DoS protection, error masking, audit logging.

**Additional tables:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'public',
    company_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    action TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    criteria JSONB,
    result_count INT,
    duration_ms INT,
    ip_address TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

Products table adds `company_id INT NOT NULL`.

**Endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/register` | Register with role + company_id |
| POST | `/login` | Login → JWT with user_id, role, company_id |
| POST | `/search` | Auth + scoped + field-projected search |
| GET | `/products/{id}` | Scoped by company |

**JWT claims** (extended from lab10):
```json
{"user_id": 1, "username": "alice", "role": "manager", "company_id": 42, "exp": ...}
```

**Key patterns:**
- **Forced scope injection:** Always append `company_id = $N` from JWT, client cannot override
- **Field projection per role:**
  - `public`: id, name, price, category
  - `manager`: + brand, in_stock
  - `admin`: + company_id, created_at
- **Rate limiting:** 20 searches/minute per user, returns 429 + Retry-After
- **Input validation:** Block leading wildcards, whitelist filter/sort fields, max 20 filters
- **Query timeout:** `context.WithTimeout(ctx, 5*time.Second)`, use `db.QueryContext()`
- **Audit logging:** middleware writes to audit_log table
- **Error masking:** Generic errors externally, full errors in server logs

**Dependencies:** `chi/v5`, `lib/pq`, `golang-jwt/jwt/v5`, `golang.org/x/crypto`

---

### Sub-lab 04: Capstone — Combining Everything (~450 lines Go)

**What it teaches:** Full integration of all patterns — authenticated saved search with multi-tenant isolation, ownership enforcement, and audit trail.

**All tables combined.** search_sessions and saved_searches gain `user_id` and `company_id` columns.

**Endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/register` | Register |
| POST | `/login` | Login → JWT |
| POST | `/search` | Auth + scoped search |
| POST | `/searches` | Save criteria (user-owned) |
| GET | `/products?searchId={id}&page=1` | Execute saved search (validates ownership) |
| POST | `/saved-searches` | Persist named search (user-owned) |
| GET | `/saved-searches` | List current user's saved searches only |
| DELETE | `/saved-searches/{id}` | Delete only if owned |
| GET | `/audit-log` | Admin-only: view audit trail |

**Key additions over sub-lab 03:**
- Saved search ownership enforcement (`WHERE user_id = $1 AND company_id = $2`)
- Two users in different companies see different data
- Admin-only audit log viewer with date range filtering
- Full end-to-end multi-tenant data isolation

---

## CLAUDE.md Knowledge Base

`lab21-search-api/CLAUDE.md` — structured reference from user's learning sessions:

1. Why POST /search instead of GET (URL limits, 4 solution comparison + decision matrix)
2. Filter operator design (eq, neq, gt, gte, lt, lte, in, like + type safety)
3. Dynamic WHERE clause building (parameterized queries, arg indexing)
4. Saved Search pattern (two-phase, criteria-only, TTL, 410 Gone, Jira/Salesforce/Grafana refs)
5. DB schema for saved searches (search_sessions vs saved_searches, JSONB, indexes)
6. Security: injection attacks (SQL, NoSQL, Elasticsearch)
7. Security: BOLA/access control (forced scope from JWT)
8. Security: data over-exposure (field projection per role)
9. Security: DoS protection (rate limit, leading wildcard block, query timeout, max filters)
10. Security: error masking + audit logging

---

## Root README

`lab21-search-api/README.md` — follows lab11's structure:

- Learning Path table (4 sub-labs with descriptions and time estimates)
- "Which Labs Should I Do?" tiers:
  - **Must-Do:** Sub-lab 01 (POST /search) — the core pattern (~25 min)
  - **Recommended:** Sub-lab 03 (Security) — critical for production APIs (~25 min)
  - **Full experience:** Sub-labs 01 → 02 → 03 → 04 in order (~2 hours)
- "How to Run" instructions
- Project structure tree
- Key Concepts summary
- Reference link to CLAUDE.md
- Prerequisites: labs 09 and 10

---

## Main Workshop README Update

**File:** `README.md` (project root)

1. Add "Part 4: Advanced REST Deep Dives" section after Part 3:
```markdown
### Part 4: Advanced REST Deep Dives

Go deeper into specific REST API design challenges. These labs build on patterns from Part 1 and Part 2.

| # | Lab | Description |
|---|-----|-------------|
| 21 | [Advanced Search API](lab21-search-api/) | Design complex search APIs with POST /search, saved search patterns, security hardening, and multi-tenant scoping |
```

2. Update Learning Path diagram to add Part 4
3. Add to Services & Ports table: `| 21 (sub-labs 01-04) | Go API, PostgreSQL | 8080, 5432 |`

---

## Critical Existing Files to Reference

```
# Patterns to reuse/extend
lab09-pagination-and-filtering/main.go         — dynamic WHERE, PaginatedResponse, product seed
lab10-authentication/auth.go                   — JWT AuthMiddleware, context user, bcrypt
lab10-authentication/handlers.go               — register/login handlers
lab11-api-versioning/README.md                 — sub-lab README structure, learning path tiers
lab11-api-versioning/CLAUDE.md                 — knowledge base format
lab11-api-versioning/lab11-01-*/golang/        — simple sub-lab file structure

# Files to modify
README.md                                      — add Part 4 section, update learning path + ports table
```
