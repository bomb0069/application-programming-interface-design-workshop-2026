# Lab 04-08: Version Lifecycle and Observability

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

## Code Walkthrough: The Metrics Middleware

All version metrics come from one small middleware. In the Go edition it lives in `golang/main.go`; the .NET edition's equivalent is `dotnet/Middleware/PrometheusVersionMiddleware.cs`. It has three jobs: **define the metrics**, **time every request**, and **attach the right labels**.

### 1. Define the metrics (once, at startup)

```go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "api_requests_total",       // the metric name you query in PromQL
            Help: "Total number of API requests",
        },
        []string{"version", "endpoint", "method", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "api_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"version", "endpoint", "method"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}
```

Two metric types, two different questions:

- **Counter** (`api_requests_total`) — only ever goes up. Prometheus turns it into traffic *rates* with `rate(...)`; you never read the raw number directly.
- **Histogram** (`api_request_duration_seconds`) — records each duration into buckets, which is what lets the dashboard compute average latency per version (`_sum / _count`) and percentiles later.

They are package-level and registered once in `init()` — a metric must exist **once per process**, with every request incrementing the same series.

### 2. Time the request by wrapping the handler

```go
func prometheusMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &statusRecorder{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(wrapped, r)      // run the actual handler (DB query, JSON encode…)

        version := getVersion(r)
        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(version, r.URL.Path, r.Method,
            strconv.Itoa(wrapped.statusCode)).Inc()
        httpRequestDuration.WithLabelValues(version, r.URL.Path, r.Method).Observe(duration)
    })
}
```

Classic middleware sandwich: capture the clock before `next.ServeHTTP`, record after it returns. One Go-specific trick: `http.ResponseWriter` doesn't expose the status code after the fact, so the middleware wraps it in a `statusRecorder` that remembers whatever `WriteHeader` was called with.

### 3. Attach the version label — and only on versioned routes

The middleware is mounted **inside** the versioned route groups, next to a tiny middleware that stamps the version into the request context:

```go
r.Route("/api/v1", func(r chi.Router) {
    r.Use(setVersion("v1"))          // context value: "v1"
    r.Use(prometheusMiddleware)      // reads it back via getVersion(r)
    ...
})

r.Route("/api/v2", func(r chi.Router) {
    r.Use(setVersion("v2"))
    r.Use(prometheusMiddleware)
    ...
})
```

Two details here carry the whole lab:

- **The label value must match the dashboard.** Grafana's gauges query `version="v1"` — so the middleware must emit exactly `v1`/`v2`. A mismatch is silent: nothing errors, the panels just show *No data*. (The .NET edition once emitted `"1"`/`"2"` here — that was the bug behind empty Traffic Share gauges.)
- **Unversioned paths are never counted.** Because the middleware only exists inside `/api/v1` and `/api/v2`, requests to `/metrics`, `/health`, and `/api/lifecycle` don't pollute the traffic-share math: `sum(rate(api_requests_total{version="v1"})) / sum(rate(api_requests_total))` sums to 100% across v1+v2. (The .NET middleware is global, so it checks `GetRequestedApiVersion()` and skips requests that resolved no version — same effect, different mechanism.)

### A word on label cardinality

Every distinct label combination becomes its own time series in Prometheus. This lab labels by raw `endpoint` path, which is safe only because the ID space is tiny (3 seeded products). In production, label by **route template** (`/api/v1/products/{id}`), never by raw path — otherwise every distinct ID mints a new series and Prometheus memory explodes. That is also why there is no `user_id` label: high-cardinality questions ("which clients still call v1?") belong to structured logs, not metrics.

### From labels to dashboard

The two Grafana gauges are just this pipeline seen end-to-end:

```
middleware label            PromQL on the dashboard
version="v2"       ─────▶   sum(rate(api_requests_total{version="v2"}[1m]))
                            ─────────────────────────────────────────────── × 100
all versioned reqs ─────▶   sum(rate(api_requests_total[1m]))
```

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

### Generate traffic manually

```bash
# V1 traffic
for i in $(seq 1 10); do curl -s -o /dev/null -w "%{http_code} " http://localhost:8080/api/v1/products; done; echo

# V2 traffic
for i in $(seq 1 20); do curl -s -o /dev/null -w "%{http_code} " http://localhost:8080/api/v2/products; done; echo
```

### Run the Load Test

For sustained traffic that produces meaningful Grafana dashboards, use the built-in load test:

```bash
# Start the load test (runs for 3 minutes by default)
docker compose --profile loadtest up loadtest
```

The load test sends ~5 requests/second with a 30/70 split between v1 and v2, simulating a real migration scenario where most clients have already moved to v2.

**Customize the traffic ratio:**

```bash
# Equal traffic split
V1_WEIGHT=50 V2_WEIGHT=50 docker compose --profile loadtest up loadtest

# Almost fully migrated — only 5% on v1
V1_WEIGHT=5 V2_WEIGHT=95 docker compose --profile loadtest up loadtest

# Longer run with higher throughput
DURATION_SECONDS=300 REQUESTS_PER_SEC=10 docker compose --profile loadtest up loadtest
```

The script generates GET (list + by ID) and POST requests across both versions. It prints progress every 30 seconds and a summary at the end.

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

### Grafana Dashboard

The Grafana datasource and dashboard are **auto-provisioned** — no manual setup required.

1. Open http://localhost:3000 (admin/admin)
2. Go to Dashboards -- the **"API Version Traffic"** dashboard is already loaded

The dashboard includes 5 panels:

| Panel | PromQL | What It Shows |
|-------|--------|---------------|
| Request Rate by Version | `sum by (version) (rate(api_requests_total[1m]))` | Requests/sec for v1 vs v2 over time |
| V1 Traffic Share (%) | `sum(rate(...{version="v1"}[1m])) / sum(rate(...[1m])) * 100` | Gauge showing how much traffic is still on deprecated v1 |
| V2 Traffic Share (%) | Same formula for v2 | Gauge showing current version adoption |
| Average Latency by Version | `sum by (version) (rate(duration_sum[1m])) / sum by (version) (rate(duration_count[1m]))` | Are v1 callers experiencing different latency than v2? |
| Requests by Endpoint | `sum by (version, endpoint, method) (rate(...[1m]))` | Which specific endpoints still receive v1 traffic? |

> **Tip:** Run the load test, then watch the dashboard update in real time. Try changing `V1_WEIGHT` and `V2_WEIGHT` to simulate different migration stages.

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
