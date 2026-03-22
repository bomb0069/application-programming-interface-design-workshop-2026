# Lab 11 — API Versioning Workshop Context

> This file captures the full knowledge base discussed in session. Use it as the foundation for building workshop exercises, slides, and code labs.

---

## 1. Why API versioning exists — benefits & pitfalls

### Benefits

**For API producers (your team):**

- Evolve the API (change response shapes, add breaking changes) without forcing all consumers to upgrade at once
- Deprecate old behaviour on your own schedule
- Each version can have its own SLA, rate limits, and caching strategy

**For API consumers (your clients):**

- Stability guarantee — `v1` won't break just because `v2` shipped
- Time to migrate at their own pace
- Can pin to a version and test migration separately from production changes

### Pitfalls

- **Code duplication** — two code paths doubles bug surface
- **Database drift** — v1 and v2 using different schemas/transformations
- **Zombie versions** — old versions forgotten but not removed; you keep supporting clients you didn't know existed
- **Documentation debt** — two versions, two docs, two changelogs
- **Testing burden** — every change must be tested across all live versions
- **Deprecation pain** — you think nobody uses v1 anymore… but you don't really know (this is why observability matters)

---

## 2. Observability — monitoring which version is called by whom

### The three pillars for version tracking

**Metrics (Prometheus / CloudWatch / Datadog)**

Add a `version` label to every counter and histogram at the gateway or middleware level:

```
api_requests_total{version="v1", endpoint="/point", method="GET", client_id="mobile-app"} 4521
api_requests_total{version="v2", endpoint="/point", method="GET", client_id="web-frontend"} 9820
```

Build a dashboard showing the v1/v2 traffic split over time — the key signal for deciding when it's safe to sunset v1.

**Structured logs**

Every log entry should emit `api_version` as a field (not buried in a URL string), along with `client_id` or `x-api-key`:

```json
{
  "ts": "2026-03-22T10:00:00Z",
  "api_version": "v1",
  "client": "partner-xyz",
  "endpoint": "/point",
  "latency_ms": 42
}
```

Log aggregators (Loki, Elasticsearch, Splunk) can answer "give me every unique client still calling v1".

**Distributed traces (OpenTelemetry / Jaeger / Tempo)**

Tag every span with `api.version` and `client.id`. This helps answer deeper questions: "Are v1 callers experiencing higher latency than v2?" or "Which downstream service is the bottleneck specifically for v1 paths?"

### Deprecation workflow using observability

1. **Announce sunset date** → add a `Sunset` response header to v1 responses
2. **Set a deprecation alert** → fire when `api_requests_total{version="v1"}` is still non-zero within 30 days of sunset
3. **Build a "v1 callers" report** → use log queries to extract unique `client_id` values still hitting v1; reach out to each one
4. **Track migration progress** → the v1/v2 traffic ratio on your dashboard is your KPI
5. **Kill switch** → once v1 traffic hits zero for 2 consecutive weeks, remove the handler with confidence

---

## 3. Implementing API versioning in .NET Core

### Setup

```bash
dotnet add package Asp.Versioning.Mvc
dotnet add package Asp.Versioning.Mvc.ApiExplorer
```

```csharp
// Program.cs
builder.Services.AddApiVersioning(opt => {
    opt.DefaultApiVersion = new ApiVersion(1, 0);
    opt.AssumeDefaultVersionWhenUnspecified = true;
    opt.ReportApiVersions = true; // adds "api-supported-versions" header
    opt.ApiVersionReader = new UrlSegmentApiVersionReader();
})
.AddApiExplorer(opt => {
    opt.GroupNameFormat = "'v'VVV";
    opt.SubstituteApiVersionInUrl = true;
});
```

### Scenario 1: Duplicate controller, shared service

**Use when:** Only the response shape / DTO mapping changes. Business logic is identical.

```csharp
// Controllers/V1/PointController.cs
[ApiController, ApiVersion("1.0")]
[Route("api/v{version:apiVersion}/point")]
public class PointController : ControllerBase
{
    private readonly IPointService _svc;
    public PointController(IPointService svc) => _svc = svc;

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(int id)
    {
        var point = await _svc.GetByIdAsync(id);
        // v1: flat response shape
        return Ok(new { point.Id, point.Name, point.Score });
    }
}

// Controllers/V2/PointController.cs
[ApiController, ApiVersion("2.0")]
[Route("api/v{version:apiVersion}/point")]
public class PointController : ControllerBase
{
    private readonly IPointService _svc;
    public PointController(IPointService svc) => _svc = svc;

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(int id)
    {
        var point = await _svc.GetByIdAsync(id);
        // v2: enriched nested response
        return Ok(new PointV2Response {
            Id = point.Id,
            Name = point.Name,
            Metadata = new { point.Score, point.CreatedAt, point.Tags }
        });
    }
}

// Shared service — registered once
public interface IPointService
{
    Task<Point> GetByIdAsync(int id);
    Task<IEnumerable<Point>> GetAllAsync();
}

// Program.cs
builder.Services.AddScoped<IPointService, PointService>();
```

**Pitfall:** Both controllers inject the same `IPointService`, so any service change affects both versions at once.

---

### Scenario 2: Duplicate controller + service, shared repository

**Use when:** Business logic itself differs — v2 has new scoring algorithm, new validation rules, or new methods (like pagination).

```csharp
// V1 service — keeps legacy algorithm intact
public interface IPointServiceV1
{
    Task<Point> GetByIdAsync(int id);
    Task<IEnumerable<Point>> GetAllAsync();
}

public class PointServiceV1 : IPointServiceV1
{
    private readonly IPointRepository _repo;
    public PointServiceV1(IPointRepository repo) => _repo = repo;

    public async Task<Point> GetByIdAsync(int id)
    {
        var point = await _repo.GetByIdAsync(id);
        point.Score = LegacyScoreCalc.Calculate(point); // old formula, never touch
        return point;
    }
}

// V2 service — new logic, new method (GetPaged doesn't exist in V1)
public interface IPointServiceV2
{
    Task<PointV2> GetByIdAsync(int id);
    Task<PagedResult<PointV2>> GetPagedAsync(int page, int size); // new!
}

public class PointServiceV2 : IPointServiceV2
{
    private readonly IPointRepository _repo;
    public PointServiceV2(IPointRepository repo) => _repo = repo;

    public async Task<PointV2> GetByIdAsync(int id)
    {
        var point = await _repo.GetByIdAsync(id);
        return PointMapper.ToV2(point, NewScoreCalc.Calculate(point));
    }

    public async Task<PagedResult<PointV2>> GetPagedAsync(int page, int size)
    {
        var items = await _repo.GetPagedAsync(page, size);
        return new PagedResult<PointV2>(items.Select(PointMapper.ToV2), page, size);
    }
}

// Shared repository — data access only, zero business logic
public interface IPointRepository
{
    Task<Point> GetByIdAsync(int id);
    Task<IEnumerable<Point>> GetAllAsync();
    Task<IEnumerable<Point>> GetPagedAsync(int page, int size);
}

// Program.cs
builder.Services.AddScoped<IPointServiceV1, PointServiceV1>();
builder.Services.AddScoped<IPointServiceV2, PointServiceV2>();
builder.Services.AddScoped<IPointRepository, PointRepository>(); // one registration
```

**Pitfall:** Two services that diverge significantly over time. Business logic must NOT leak into the repository.

---

### Scenario 3: Gateway / version router

**Use when:** You want one controller entry point that dispatches internally.

```csharp
// Single controller, accepts both versions
[ApiController]
[ApiVersion("1.0")]
[ApiVersion("2.0")]
[Route("api/v{version:apiVersion}/point")]
public class PointController : ControllerBase
{
    private readonly GetPointV1Handler _v1;
    private readonly GetPointV2Handler _v2;

    public PointController(GetPointV1Handler v1, GetPointV2Handler v2)
        => (_v1, _v2) = (v1, v2);

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(int id)
    {
        var version = HttpContext.GetRequestedApiVersion();
        return version?.MajorVersion switch
        {
            1 => Ok(await _v1.HandleAsync(id)),
            2 => Ok(await _v2.HandleAsync(id)),
            _ => BadRequest("Unsupported version")
        };
    }
}

// Handlers — thin, focused on one version's concerns
public class GetPointV1Handler(IPointService svc)
{
    public async Task<object> HandleAsync(int id)
    {
        var p = await svc.GetByIdAsync(id);
        return new { p.Id, p.Name, p.Score };
    }
}

public class GetPointV2Handler(IPointService svc)
{
    public async Task<object> HandleAsync(int id)
    {
        var p = await svc.GetByIdAsync(id);
        return new PointV2Response { Id = p.Id, Name = p.Name,
            Metadata = new { p.Score, p.CreatedAt } };
    }
}

// Program.cs
builder.Services.AddScoped<GetPointV1Handler>();
builder.Services.AddScoped<GetPointV2Handler>();
builder.Services.AddScoped<IPointService, PointService>();
```

**Pitfall:** The gateway controller can grow into a god class fast.

### Scenario comparison table

| Scenario                                    | Use when                                | Watch out for                                           |
| ------------------------------------------- | --------------------------------------- | ------------------------------------------------------- |
| Duplicate controller, shared service        | Only the response DTO/shape changes     | Service changes break both versions silently            |
| Duplicate controller + service, shared repo | Business logic genuinely differs        | Two services diverging — keep repo interface stable     |
| Gateway controller                          | You want one entry point, smaller teams | Controller becomes a god class; Swagger setup is fiddly |

---

## 4. URL versioning strategies — all six approaches

### Strategy 1: `/v1/api/point` — Version as root prefix

```csharp
[ApiVersion("1.0")]
[Route("v{version:apiVersion}/api/point")]
public class PointControllerV1 : ControllerBase { }
```

- **Benefits:** Immediately obvious; easy to route at reverse proxy level
- **Drawbacks:** Version prefix before `/api` breaks REST convention; Swagger grouping is awkward
- **Pitfall:** Gateway routing rules written as `/api/*` prefix will silently miss these
- **Verdict:** Uncommon. Only use if you own a very stable `/api` prefix and want the version to gate the entire surface.

---

### Strategy 2: `/api/v1/point` — Version after /api ⭐ Recommended

```csharp
[ApiVersion("1.0")]
[Route("api/v{version:apiVersion}/point")]
public class PointControllerV1 : ControllerBase { }

[ApiVersion("2.0")]
[Route("api/v{version:apiVersion}/point")]
public class PointControllerV2 : ControllerBase { }
```

- **Benefits:** Industry convention (GitHub, Stripe, Twitter); version is visible but `/api/` prefix is stable; trivial to bookmark, test in browser, share in tickets
- **Drawbacks:** URLs change when version changes; violates strict REST purism
- **Pitfall:** Forgetting to add `[MapToApiVersion]` causes both controllers to respond to both versions
- **Verdict:** Default choice for most teams.

---

### Strategy 3: `/api/point/v1` — Version as resource suffix

```csharp
[ApiVersion("1.0")]
[Route("api/point/v{version:apiVersion}")]
public class PointControllerV1 : ControllerBase { }
```

- **Benefits:** Resource name stays prominent and comes first
- **Drawbacks:** Clashes with RESTful sub-resources; `/api/point/v1` looks like a nested resource named `v1`; almost no popular public API uses this
- **Pitfall:** Route constraint conflicts — `/api/point/{id:int}` and `/api/point/v{ver}` fight each other
- **Verdict:** Avoid.

---

### Strategy 4: `/api/point-v1` — Version baked into resource name

```csharp
[Route("api/point-v1")]
public class PointV1Controller : ControllerBase { }

[Route("api/point-v2")]
public class PointV2Controller : ControllerBase { }
```

- **Benefits:** No routing conflicts possible; works without any versioning framework
- **Drawbacks:** Not actually versioning — just naming different resources; no tooling support; loses ALL framework versioning features; no `Sunset` headers, no version negotiation
- **Pitfall:** You build your own deprecation logic from scratch
- **Verdict:** Anti-pattern. Only as absolute last resort.

---

### Strategy 5: `?api-version=1` — Query string versioning

```csharp
builder.Services.AddApiVersioning(opt => {
    opt.ApiVersionReader = new QueryStringApiVersionReader("api-version");
    opt.DefaultApiVersion = new ApiVersion(1, 0);
    opt.AssumeDefaultVersionWhenUnspecified = true;
});

[ApiVersion("1.0")]
[ApiVersion("2.0")]
[Route("api/point")]
public class PointController : ControllerBase {
    [HttpGet, MapToApiVersion("1.0")]
    public IActionResult GetV1() => Ok("v1");

    [HttpGet, MapToApiVersion("2.0")]
    public IActionResult GetV2() => Ok("v2");
}
```

- **Benefits:** Resource URL stays stable; easy to test in browser; natively supported by `Asp.Versioning`
- **Drawbacks:** Version leaks into logs/bookmarks; breaks HTTP caching; version can be accidentally omitted
- **Pitfall:** CDN/proxy caches store separate responses per `?api-version=` value — cache miss rate inflation
- **Verdict:** Excellent as a secondary reader alongside URL versioning. Avoid as the sole strategy for CDN-cached public APIs.

---

### Strategy 6: `X-Api-Version: 1.0` — Custom request header

```csharp
builder.Services.AddApiVersioning(opt => {
    opt.ApiVersionReader = new HeaderApiVersionReader("X-Api-Version");
    opt.DefaultApiVersion = new ApiVersion(1, 0);
    opt.AssumeDefaultVersionWhenUnspecified = true;
});

[ApiVersion("1.0")]
[ApiVersion("2.0")]
[Route("api/point")]
public class PointController : ControllerBase {
    [HttpGet, MapToApiVersion("1.0")]
    public IActionResult GetV1() => Ok();
    [HttpGet, MapToApiVersion("2.0")]
    public IActionResult GetV2() => Ok();
}
```

- **Benefits:** URL stays completely clean; no caching pollution; ideal for internal service-to-service APIs
- **Drawbacks:** Not testable in a plain browser; easy to forget; not visible in logs unless explicitly captured
- **Pitfall:** If `AssumeDefaultVersionWhenUnspecified=true` and client forgets the header, they silently get v1 forever — silent API contract drift
- **Verdict:** Best for internal microservice APIs where all consumers are controlled services.

---

### Strategy 7: `Accept: application/vnd.myapi.v1+json` — Media type / content negotiation

```csharp
// Custom reader — not built-in to Asp.Versioning
public class MediaTypeApiVersionReader : IApiVersionReader {
    public void AddParameters(IApiVersionParameterDescriptionContext ctx) { }
    public IReadOnlyList<string> Read(HttpRequest req) {
        var accept = req.Headers.Accept.ToString();
        var match = Regex.Match(accept, @"vnd\.myapi\.v(\d+)\+json");
        return match.Success
            ? [match.Groups[1].Value + ".0"]
            : [];
    }
}
```

- **Benefits:** True REST compliance; URL is perfectly stable forever; GitHub API v3 uses this
- **Drawbacks:** High developer friction; minimal tooling/Swagger support; middleware often strips custom media types
- **Pitfall:** CORS, caching proxies, and API gateways often strip Accept header variants
- **Verdict:** Architecturally correct per REST purists. Practically painful for most teams.

---

### Combining multiple readers (recommended production pattern)

```csharp
builder.Services.AddApiVersioning(opt => {
    opt.ApiVersionReader = ApiVersionReader.Combine(
        new UrlSegmentApiVersionReader(),               // /api/v1/point  ← primary
        new QueryStringApiVersionReader("api-version"), // ?api-version=1  ← fallback for tools
        new HeaderApiVersionReader("X-Api-Version")    // header ← internal services
    );
    opt.DefaultApiVersion = new ApiVersion(1, 0);
    opt.AssumeDefaultVersionWhenUnspecified = true;
    opt.ReportApiVersions = true;
});
```

### Quick decision guide

| Your situation                     | Best choice                           |
| ---------------------------------- | ------------------------------------- |
| Public API, external developers    | `/api/v1/resource` — URL path mid     |
| Internal microservices only        | `X-Api-Version` header                |
| Strict REST purism required        | Media type (Accept header)            |
| Need URL stable + browser testable | Query string `?api-version=`          |
| CDN heavy, cache-sensitive         | URL path (avoids cache fragmentation) |
| All of the above                   | Combine URL + query string + header   |

---

## 5. Breaking change classification

### Safe — no new version required

- Add an optional request field
- Add a new field to the response body
- Add a new optional query parameter
- Add a new HTTP method to existing resource
- Add a new endpoint / resource entirely
- Widen an accepted value range (e.g. string max length 50 → 200)
- Fix a bug that makes behaviour match documented spec

### Breaking — new version required

- Remove or rename a request or response field
- Change a field type (string → int, array → object)
- Change an HTTP status code for the same outcome
- Make a previously optional request field required
- Narrow an accepted value range (max length 200 → 50)
- Remove a supported HTTP method
- Change authentication mechanism
- Change pagination structure (cursor vs offset)
- Remove an endpoint entirely

### Context-dependent

- **Add a new enum value** — safe if clients use exhaustive switch without default; breaking if they crash on unknown values
- **Change error response shape** — safe if clients only check HTTP status; breaking if they parse error body
- **Add required header** — safe for internal clients you control; breaking for external consumers
- **Change URL casing or trailing slash** — technically safe by HTTP spec; breaks many hardcoded clients in practice

> **The new enum value rule of thumb:** Document your extensibility contract explicitly. If your API docs say "clients MUST handle unknown enum values gracefully", adding a new value is always safe.

---

## 6. Version numbering schemes

### Integer versioning — `v1`, `v2`, `v3`

Used by: Stripe, GitHub, Twilio, most REST APIs.

**Rule:** Increment only when you have a breaking change. Non-breaking additions stay in the same version.

**Pitfall:** `v1` becomes a frozen-in-time snapshot that accumulates technical debt. Teams end up with `v1` routes that can never be refactored.

---

### Semantic versioning — `v1.2.0`

Major.Minor.Patch — used more in SDK/library versioning than HTTP APIs.

**Rule:** major = breaking, minor = additive, patch = bug fix.

**Pitfall:** Exposing all three in the URL creates an explosion of URL variants. Usually only the major is in the URL; minor/patch are changelog-only.

---

### Date-based versioning — `2024-01-15`

Used by: Stripe (new system), Azure REST APIs, Anthropic API.

**Rule:** Each date is a complete API snapshot. Clients pin to a date; the changelog IS the diff between two dates.

```csharp
// Custom date-based version reader
public class DateVersionReader : IApiVersionReader {
    public IReadOnlyList<string> Read(HttpRequest req) {
        var date = req.Headers["Api-Version"].ToString();
        // map "2024-01-15" → "20240115" for ApiVersion
        return DateTime.TryParse(date, out var d)
            ? [d.ToString("yyyyMMdd")]
            : [];
    }
}

// Controller attribute
[ApiVersion("20240101")]
[ApiVersion("20240601")]
[Route("api/point")]
public class PointController : ControllerBase { }
```

**Pitfall:** Many active versions accumulate quickly. MUST enforce a sunset window (e.g. retire versions older than 18 months).

---

## 7. Deprecation & Sunset headers (RFC 8594)

HTTP defines two standard response headers for communicating version deprecation to clients.

### Deprecation header

Tells clients this version is deprecated. Value = ISO 8601 date when deprecation began, or `true`.

```
Deprecation: Thu, 01 Jan 2026 00:00:00 GMT
```

### Sunset header

Tells clients the exact date this version will stop responding. After that date, return HTTP 410 Gone.

```
Sunset: Wed, 01 Jul 2026 00:00:00 GMT
```

### Link header (RFC 8288)

Points clients to the migration guide.

```
Link: <https://api.example.com/docs/migrate-v1-v2>; rel="successor-version"
```

### .NET implementation — deprecation middleware

```csharp
public class DeprecationMiddleware(RequestDelegate next) {
    private static readonly Dictionary<string, (string dep, string sun)> _schedule = new() {
        ["1.0"] = ("Sat, 01 Mar 2026 00:00:00 GMT", "Wed, 01 Jul 2026 00:00:00 GMT")
    };

    public async Task InvokeAsync(HttpContext ctx) {
        await next(ctx);
        var ver = ctx.GetRequestedApiVersion()?.ToString();
        if (ver != null && _schedule.TryGetValue(ver, out var s)) {
            ctx.Response.Headers["Deprecation"] = s.dep;
            ctx.Response.Headers["Sunset"] = s.sun;
            ctx.Response.Headers["Link"] = "</docs/migrate>;rel=\"successor-version\"";
        }
    }
}

// Register in Program.cs
app.UseMiddleware<DeprecationMiddleware>();
```

> **Key insight:** Clients that consume these headers can automatically log warnings like "You are calling deprecated API v1.0. Sunset: 2026-07-01. Migrate by then." This is far more reliable than manual email campaigns.

---

## 8. Backward compatibility techniques

### Additive-only response fields

Always add new fields, never remove or rename. Clients that don't know about a field simply ignore it.

```json
// v1 response
{ "id": 1, "name": "Gold" }

// v2 response — same endpoint, no new version needed
{ "id": 1, "name": "Gold", "tier": "premium", "createdAt": "2024-01-01" }
```

### Field deprecation annotation

Mark fields deprecated in the response rather than removing them immediately.

```csharp
public class PointResponse {
    public int Id { get; set; }

    [Obsolete("Use TierName instead")]
    public string? LegacyCategory { get; set; }  // still returned

    public string TierName { get; set; }          // new field
}
```

### Tolerant reader pattern

Design clients to ignore unknown fields and unknown enum values — never throw on unexpected data.

```csharp
// Use JsonExtensionData to absorb unknown fields gracefully
public class PointDto {
    public int Id { get; set; }
    public string Name { get; set; }

    [JsonExtensionData]
    public Dictionary<string, JsonElement>? Extra { get; set; }
}
```

> **Warning:** The #1 mistake is removing a field because "nobody uses it" without checking. Use your observability data (log field access patterns) before removing anything.

---

## 9. Contract testing (consumer-driven)

Contract testing solves the problem that unit tests can't detect: your API changed in a way that breaks a specific consumer even though all your own tests pass.

### How it works (Pact / PactNet)

Each consumer records what it actually uses from your API — specific fields, statuses, shapes. That record is the "contract". Your API's CI pipeline verifies every contract still holds before merging.

### Without vs with contract testing

| Without                                                                                                     | With                                                                                                                |
| ----------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| Server team renames `score` to `points`. All server tests pass. Consumer app crashes in prod. Found at 2am. | Consumer contract says "I need field `score`". Server PR runs pact verify. Fails immediately — caught before merge. |

### PactNet implementation

```csharp
// Consumer side — records what fields it needs
pact.UponReceiving("a request for a point")
    .WithRequest(HttpMethod.Get, "/api/v1/point/1")
    .WillRespond()
    .WithStatus(HttpStatusCode.OK)
    .WithJsonBody(new {
        id = Match.Integer(),
        score = Match.Decimal()   // only fields consumer cares about
    });

// Provider side — verifies all contracts in CI
new PactVerifier("PointService", options)
    .WithHttpEndpoint(new Uri("http://localhost:5000"))
    .WithPactBrokerSource(new Uri("https://your-pact-broker"))
    .Verify();
```

> Contract tests live in CI, not production. Run them as a required check on every PR. They are the safety net that lets you refactor confidently across version boundaries.

---

## 10. Version lifecycle management

Every version goes through the same lifecycle. The discipline is having a **written policy before you ship v1**.

### Lifecycle stages

| Stage           | Description                                                                                                                                     |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| **Current**     | Latest version. Actively developed. No deprecation headers. All new features land here.                                                         |
| **Maintained**  | Older version receiving security patches and critical bug fixes only. No new features. No deprecation headers yet.                              |
| **Deprecated**  | `Deprecation` + `Sunset` headers present on every response. Link header points to migration guide. Active outreach to consumers still using it. |
| **End of life** | Returns HTTP 410 Gone for all requests. Zero maintenance cost. Tombstone endpoint explains where to migrate.                                    |

### Recommended policy template

Write this **before** v1 ships:

- Minimum support window: 12 months from release date
- Deprecation notice: 6 months before sunset
- Sunset headers appear immediately on deprecation
- Sunset triggered when: traffic drops below 1% **OR** support window expires (whichever comes first)

### HTTP 410 Gone tombstone response

```json
{
  "error": "VERSION_SUNSET",
  "message": "API v1 was sunset on 2026-07-01",
  "migrateUrl": "/docs/v1-v2"
}
```

### .NET sunset filter

```csharp
public class SunsetFilter : IActionFilter {
    private static readonly Dictionary<string, DateTime> _sunsets = new() {
        { "1.0", new DateTime(2026, 7, 1) }
    };

    public void OnActionExecuting(ActionExecutingContext ctx) {
        var ver = ctx.HttpContext.GetRequestedApiVersion()?.ToString();
        if (ver != null && _sunsets.TryGetValue(ver, out var sun) &&
            DateTime.UtcNow > sun) {
            ctx.Result = new ObjectResult(new {
                error = "VERSION_SUNSET",
                message = $"v{ver} sunset on {sun:yyyy-MM-dd}",
                migrateUrl = "/docs/migrate"
            }) { StatusCode = 410 };
        }
    }

    public void OnActionExecuted(ActionExecutedContext ctx) { }
}
```

---

## Workshop exercise ideas

### Exercise 1 — Identify breaking changes (15 min)

Given a list of 10 proposed API changes, classify each as safe / breaking / context-dependent. Discuss the enum value edge case.

### Exercise 2 — Implement URL versioning (30 min)

Start from a single `PointController`. Add `Asp.Versioning`. Create V1 and V2 with different response shapes sharing one service.

### Exercise 3 — Add deprecation middleware (20 min)

Add the `DeprecationMiddleware` to the project. Verify with curl that `Deprecation`, `Sunset`, and `Link` headers appear on V1 responses only.

### Exercise 4 — Observability (20 min)

Add Prometheus middleware. Label metrics with `api_version`. Build a Grafana query showing V1 vs V2 request share over time.

### Exercise 5 — Contract test (30 min)

Write a PactNet consumer test for the V1 `/point/{id}` endpoint that pins to the `score` field. Then rename `score` to `points` in the service and watch the pact verify step fail in CI.

### Exercise 6 — Lifecycle planning (15 min, discussion)

Write a version lifecycle policy doc for the workshop API. Set a sunset date for V1. Wire up the `SunsetFilter`. Return 410 Gone and verify behaviour.

---

## Key references

- `Asp.Versioning` NuGet: `Asp.Versioning.Mvc` + `Asp.Versioning.Mvc.ApiExplorer`
- RFC 8594 — The Sunset HTTP Header Field
- RFC 8288 — Web Linking (Link header)
- PactNet — https://github.com/pact-foundation/pact-net
- Real-world examples: Stripe API versioning docs, GitHub REST API, Anthropic API
