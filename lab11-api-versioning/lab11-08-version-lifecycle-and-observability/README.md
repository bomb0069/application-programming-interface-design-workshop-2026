# Lab 11-08: Version Lifecycle and Observability

## Overview

This lab demonstrates the complete version lifecycle management and the three pillars of version observability: metrics, structured logs, and (conceptually) distributed traces. It includes Prometheus and Grafana for visualizing API version traffic.

## Version Lifecycle Stages

| Stage | Description | Headers | HTTP Status |
|-------|-------------|---------|-------------|
| **Current** | Latest version. Actively developed. All new features land here. | None | 200 |
| **Maintained** | Older version. Security patches and critical bug fixes only. | None | 200 |
| **Deprecated** | `Deprecation` + `Sunset` headers on every response. Active outreach to consumers. | Deprecation, Sunset, Link | 200 |
| **End of Life** | Returns 410 Gone. Zero maintenance cost. | - | 410 |

## Recommended Policy Template

Write this **before** v1 ships:

- Minimum support window: 12 months from release date
- Deprecation notice: 6 months before sunset
- Sunset headers appear immediately on deprecation
- Sunset triggered when: traffic drops below 1% OR support window expires (whichever comes first)

## The Three Pillars of Version Observability

### 1. Metrics (Prometheus)

```
api_requests_total{version="v1", endpoint="/api/v1/products", method="GET", status="200"} 4521
api_requests_total{version="v2", endpoint="/api/v2/products", method="GET", status="200"} 9820
```

Build a dashboard showing the V1/V2 traffic split over time -- the key signal for deciding when to sunset V1.

### 2. Structured Logs

```json
{
  "ts": "2026-03-22T10:00:00Z",
  "api_version": "v1",
  "method": "GET",
  "endpoint": "/api/v1/products",
  "status": 200,
  "latency_ms": 42,
  "user_agent": "curl/8.4.0"
}
```

Log aggregators can answer "give me every unique client still calling v1".

### 3. Distributed Traces (Conceptual)

Tag every span with `api.version` and `client.id`. This helps answer: "Are v1 callers experiencing higher latency than v2?"

## Getting Started

```bash
cd golang  # or cd dotnet
docker compose up --build
```

This starts 4 services:

| Service | URL | Purpose |
|---------|-----|---------|
| API | http://localhost:8080 | The versioned API |
| Prometheus | http://localhost:9090 | Metrics collection |
| Grafana | http://localhost:3000 | Dashboards (admin/admin) |
| PostgreSQL | localhost:5432 | Database |

## Try It Out

### Generate traffic for both versions

```bash
# V1 traffic
for i in $(seq 1 10); do curl -s http://localhost:8080/api/v1/products > /dev/null; done

# V2 traffic
for i in $(seq 1 20); do curl -s http://localhost:8080/api/v2/products > /dev/null; done
```

### Check Prometheus metrics

```bash
curl http://localhost:8080/metrics | grep api_requests_total
```

### View lifecycle information

```bash
curl http://localhost:8080/api/lifecycle | jq
```

Expected:

```json
{
  "versions": [
    {"version": "v1", "stage": "deprecated", "released_at": "2025-01-01"},
    {"version": "v2", "stage": "current", "released_at": "2026-01-01"}
  ],
  "policy": {
    "minimum_support_window": "12 months from release",
    "deprecation_notice": "6 months before sunset",
    "sunset_trigger": "Traffic below 1% OR support window expires"
  }
}
```

### Prometheus Queries

Open Prometheus at http://localhost:9090 and try:

```promql
# Total requests by version
sum by (version) (api_requests_total)

# Request rate by version (per second)
sum by (version) (rate(api_requests_total[5m]))

# V1 traffic percentage
sum(api_requests_total{version="v1"}) / sum(api_requests_total) * 100

# Average request duration by version
sum by (version) (rate(api_request_duration_seconds_sum[5m])) / sum by (version) (rate(api_request_duration_seconds_count[5m]))
```

### Grafana Setup

1. Open http://localhost:3000 (admin/admin)
2. Add Prometheus data source: http://prometheus:9090
3. Create a dashboard with the queries above

## Deprecation Workflow Using Observability

1. **Announce sunset date** -- add Sunset response header to v1 responses
2. **Set deprecation alert** -- fire when `api_requests_total{version="v1"}` is still non-zero within 30 days of sunset
3. **Build a "v1 callers" report** -- use log queries to extract unique user agents still hitting v1
4. **Track migration progress** -- the v1/v2 traffic ratio is your KPI
5. **Kill switch** -- once v1 traffic hits zero for 2 consecutive weeks, remove the handler

## Contract Testing (Concept)

Contract testing solves the problem that unit tests can't detect: your API changed in a way that breaks a specific consumer.

| Without Contract Testing | With Contract Testing |
|--------------------------|----------------------|
| Server renames `score` to `points`. All server tests pass. Consumer crashes in prod. | Consumer contract says "I need `score`". Server PR fails pact verify. Caught before merge. |

## Key Concepts

- Version lifecycle must be planned before v1 ships
- Metrics with version labels are the primary signal for sunset decisions
- Structured logs enable per-client migration tracking
- Prometheus + Grafana provide real-time version traffic visibility
- Contract testing catches cross-version breakage in CI
- The deprecation workflow is observability-driven, not calendar-driven

## Cleanup

```bash
docker compose down -v
```
