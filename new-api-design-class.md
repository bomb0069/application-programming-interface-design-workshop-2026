# REST API Design Essentials

### Principles, Security & Observability — Implementation Workshop 2026

---

## 🎯 Course Overview

| Item              | Detail                                                                       |
| ----------------- | ---------------------------------------------------------------------------- |
| **Course Name**   | REST API Design Essentials: Principles, Security & Observability             |
| **Duration**      | 2 Days                                                                       |
| **Audience Size** | 50–70 people                                                                 |
| **Focus**         | API Design Practice (not programming-heavy)                                  |
| **Target**        | Teams with inconsistent API standards across internal/external collaboration |

---

## 🧩 Real Pain Points (Requirements from Business)

These are the actual problems the team faces — the course must address all of them:

| Pain Point                                                    | Priority  |
| ------------------------------------------------------------- | --------- |
| Naming Convention — ไม่เคยตรงกัน                              | 🔴 High   |
| Path Variable vs Query Parameter — ไม่รู้ว่าควรใช้แบบไหน      | 🔴 High   |
| Search Parameter 10–20 ตัว — URL Limit ปัญหา                  | 🔴 High   |
| Sensitive Data — จัดการยังไง                                  | 🔴 High   |
| API Versioning — ทำแล้วแต่แก้ยังกระทบอยู่                     | 🔴 High   |
| Response Code Standards — เช่น search ไม่เจอ ควร 200 หรือ 404 | 🟡 Medium |
| API Key Management                                            | 🟡 Medium |
| Rate Limiting — ทำที่ไหน ทำยังไง                              | 🟡 Medium |
| ปัจจุบันใช้แค่ GET/POST — ควรใช้ Method อื่นไหม               | 🟡 Medium |

---

## 📚 Course Outline

### Module 1: REST API Design Best Practices _(Day 1 — Full Day)_

#### 1.1 Introduction to REST API

- REST Constraints & Principles
- REST vs GraphQL vs gRPC — why REST still matters
- Course goal: establishing team-wide API standards

#### 1.2 Resource Modeling & Naming Convention ⭐

- Resource naming rules: nouns, plural, lowercase, kebab-case
- Nested resources — when and how deep
- Standard prefixes: `/api/v1/`, `/internal/`, `/external/`
- **Workshop exercise:** Review & fix real naming examples from the team

#### 1.3 Request Design ⭐

- HTTP Methods — GET, POST, PUT, PATCH, DELETE
  - When to use each (address current GET/POST-only culture)
- **Path Variable vs Query Parameter — Decision Rules**
  - Rule: Path = Identity (`/users/{id}`), Query = Filter/Option (`?status=active`)
  - Decision tree with real-world examples
- **Complex Search: 10–20 Parameters Pattern**
  - URL length limits problem
  - Solution: `POST /resources/search` with JSON body
  - Saved Search pattern for reusable filters
- Pagination, Filtering, Sorting standards

#### 1.4 Response Design ⭐

- **Status Code Standards with Scenarios**
  - 200 vs 201 vs 204
  - Search not found: 200 with empty vs 404 — team decision
  - 4xx client errors with consistent error body structure
  - 5xx: what NOT to expose to clients
- Standard response envelope schema
  ```json
  {
    "data": {},
    "error": null,
    "meta": {},
    "pagination": {}
  }
  ```
- Data formats: dates (ISO 8601), nulls, enums

#### 1.5 API Versioning ⭐

- Versioning strategies: URI path / Header / Query Param
- Semantic versioning for APIs
- **Breaking vs Non-breaking changes** (core problem)
  - Safe: adding new optional fields
  - Breaking: renaming fields, removing fields, changing types
- Deprecation strategy & timeline
- How to evolve API without breaking consumers

#### 1.6 API Documentation

- Design-First vs Code-First (with recommendation)
- OpenAPI / Swagger — writing good specs
- Documentation as a contract between teams

#### 🛠️ Workshop 1: Define Your Team API Standard

> Output: A living document — Naming Convention Guide, Response Format Standard, Status Code Rules — signed off by the team

---

### Module 2: API Security _(Day 2 — Morning)_

#### 2.1 Authentication & Authorization

- OAuth2 / JWT — practical patterns
- Internal vs External API auth differences

#### 2.2 API Key Management ⭐

- API Key lifecycle: creation, rotation, revocation
- Where to store: never in URL params or logs
- How to pass: `Authorization` header vs Query Param (and why header wins)
- Key rotation strategy without downtime

#### 2.3 Sensitive Data Handling ⭐ _(new section)_

- What counts as sensitive: PII, financial data, credentials
- **Rules:**
  - Never in URL / query params
  - Never in logs
  - Never over-expose in response body
- Data masking in responses (e.g., `card: "****1234"`)
- Field-level security: return only what the consumer needs
- HTTPS / TLS basics

#### 2.4 API Rate Limiting ⭐

- Why rate limiting matters: abuse, cost, system stability
- Strategies: Fixed window, Sliding window, Token bucket
- **Where to implement:** App layer vs API Gateway vs both
- Rate limit response headers: `X-RateLimit-Limit`, `Retry-After`
- Graceful response: `429 Too Many Requests`

#### 2.5 Defense Patterns

- Input validation & sanitization
- OWASP API Security Top 10 (overview)
- Health check endpoint — what to expose vs what to hide

#### 2.6 API Gateway

- Role in security architecture
- Centralized auth, rate limiting, logging at gateway level
- Internal vs External gateway patterns

#### 🛠️ Workshop 2: Security Review Checklist

> Output: Security checklist applied to the team's existing APIs

---

### Module 3: API Architecture Design _(Day 2 — Early Afternoon)_

#### 3.1 Architecture Evolution (Brief)

- Monolith → N-tier → Microservices trade-offs
- Choosing the right architecture for your context

#### 3.2 Performance & Scalability Design

- Database design: Read vs Write operations
- Normalization vs Denormalization
- SQL tuning basics & Connection Pool
- Caching strategies: Cache-aside, Write-through, Read-through
- Redis as Key-Value cache

#### 3.3 Resilience Patterns

- Circuit Breaker pattern
- Health check design

#### 🛠️ Workshop 3: Architecture Design for Your Use Case

---

### Module 4: API Observability _(Day 2 — Late Afternoon)_

#### 4.1 Centralized Logging

- Structured log format (JSON)
- What to log — and what NOT to log (sensitive data!)
- Log levels and when to use them
- Correlation ID for request tracing

#### 4.2 Application Metrics

- Prometheus + Grafana setup
- Key API metrics: latency, error rate, throughput (RED method)

#### 4.3 Distributed Tracing

- OpenTelemetry — instrumentation basics
- Jaeger for trace visualization
- Connecting logs + metrics + traces

#### 🛠️ Workshop 4: Observability Setup for Your API

---

## 🎯 วัตถุประสงค์ (Course Objectives)

1. ฝึกการวิเคราะห์ ออกแบบ และกำหนดมาตรฐานการออกแบบ REST API ให้สอดคล้องกับหลักการและ Best Practices ที่เป็นสากล
   1. Resource Modeling และ Naming Convention
   2. การเลือกใช้ HTTP Methods ให้ถูกต้องตามวัตถุประสงค์
   3. การออกแบบ Request (Path Variable, Query Parameter, Request Body)
   4. การออกแบบ Response และ Status Codes
   5. API Versioning และการจัดการ Breaking Changes

2. ฝึกวิธีปฏิบัติในการออกแบบและจัดทำเอกสาร API ด้วยแนวทาง API Documentation ผ่านมาตรฐาน OpenAPI / Swagger ทั้งในรูปแบบ Design-first และ Code-first

3. ฝึกวิธีปฏิบัติในการออกแบบ API ให้มีความมั่นคงปลอดภัย (API Security) ครอบคลุมตั้งแต่ระดับการออกแบบจนถึงการนำไปใช้งานจริง
   1. Authentication และ Authorization
   2. การจัดการ API Key อย่างถูกต้องและปลอดภัย
   3. การจัดการ Sensitive Data ไม่ให้รั่วไหลผ่าน API
   4. API Rate Limiting และการป้องกันการใช้งานในทางที่ผิด
   5. การใช้งาน API Gateway เพื่อรวมศูนย์การจัดการความปลอดภัย

4. ฝึกวิธีปฏิบัติในการออกแบบ Architecture เพื่อรองรับประสิทธิภาพและความสามารถในการขยายตัวของระบบ (Performance & Scalability)
   1. การออกแบบฐานข้อมูลสำหรับ Read และ Write Operations
   2. การใช้งาน Caching Strategy และ Key-Value Database
   3. แนวคิด Circuit Breaker และ Health Check

5. ฝึกวิธีปฏิบัติในการออกแบบและติดตั้งระบบ API Observability เพื่อให้สามารถติดตาม วิเคราะห์ และแก้ไขปัญหาของ API ในระบบ Production ได้อย่างมีประสิทธิภาพ
   1. Centralized Logging และ Structured Log
   2. Application Metrics ด้วย Prometheus และ Grafana
   3. Distributed Tracing ด้วย OpenTelemetry และ Jaeger

---

## 🗓️ 2-Day Schedule

| Time        | Day 1                                                  | Day 2                                        |
| ----------- | ------------------------------------------------------ | -------------------------------------------- |
| 09:00–10:30 | REST Intro + Naming Convention                         | API Security (Auth, API Key, Sensitive Data) |
| 10:30–10:45 | Break                                                  | Break                                        |
| 10:45–12:00 | Request Design (Method, Path vs Query, Complex Search) | Rate Limiting + API Gateway                  |
| 12:00–13:00 | Lunch                                                  | Lunch                                        |
| 13:00–14:30 | Response Design + Versioning                           | Architecture Design                          |
| 14:30–14:45 | Break                                                  | Break                                        |
| 14:45–16:00 | Documentation + Workshop 1                             | Observability + Workshop 4                   |
| 16:00–17:00 | Workshop 1 (continued)                                 | Q&A + Retrospective                          |

---

## 🛠️ Workshop Deliverables

| Workshop   | Expected Output                                                    |
| ---------- | ------------------------------------------------------------------ |
| Workshop 1 | Team API Standard Document (Naming, Response format, Status codes) |
| Workshop 2 | Security Review Checklist for existing APIs                        |
| Workshop 3 | Architecture decision diagram for a given scenario                 |
| Workshop 4 | Observability dashboard draft + logging format spec                |

---

## 📎 Implementation Notes

- **Tech Stack:** .NET (ASP.NET Core Web API)
- **Documentation:** Swagger / OpenAPI (Swashbuckle)
- **Security Demo:** JWT + API Key middleware
- **Observability Stack:** OpenTelemetry + Jaeger + Prometheus + Grafana
- **Workshop Format:** Group-based (teams of 5–7 people), present findings to class
- **Sample Project:** Shared base API project in the repo for all workshops to build upon

---

## 📁 Project Structure (Proposed)

```
application-programming-interface-design-workshop-2026/
├── dotnet-api-design-course.md       ← this file
├── slides/
│   ├── module-1-design-best-practices/
│   ├── module-2-security/
│   ├── module-3-architecture/
│   └── module-4-observability/
├── workshops/
│   ├── workshop-1-api-standard/
│   ├── workshop-2-security-checklist/
│   ├── workshop-3-architecture/
│   └── workshop-4-observability/
├── sample-api/                        ← base .NET project
│   ├── src/
│   └── docker-compose.yml
└── resources/
    ├── api-naming-cheatsheet.md
    ├── status-code-guide.md
    └── security-checklist.md
```

---

_Last updated: March 2026 — summarized from course design discussion_
