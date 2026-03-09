# Lab 15 - Webhooks

## Learning Objectives

- Implement a webhook sender and receiver
- Event-driven notifications between services
- HMAC signature verification for webhook security
- Retry logic with exponential backoff
- Webhook registration and management

## Architecture

```
┌─────────────────────────┐         Webhook POST         ┌─────────────────────────┐
│                         │  ──────────────────────────>  │                         │
│   Sender                │    X-Webhook-Event            │   Receiver              │
│   (Order Service)       │    X-Webhook-Signature        │   (Webhook Consumer)    │
│   :8080                 │    JSON Payload               │   :9090                 │
│                         │                               │                         │
│  - POST /orders         │                               │  - POST /webhook        │
│  - GET  /orders         │                               │  - GET  /events         │
│  - PUT  /orders/:id     │                               │  - GET  /health         │
│  - POST /webhooks       │                               │                         │
│  - GET  /webhooks       │                               │                         │
│  - DELETE /webhooks/:id │                               │                         │
└───────────┬─────────────┘                               └─────────────────────────┘
            │
            │
    ┌───────┴───────┐
    │  PostgreSQL   │
    │  :5432        │
    └───────────────┘
```

## Getting Started

Start all services with Docker Compose:

```bash
docker-compose up --build
```

This starts three containers:
- **sender** - Order service on port 8080
- **receiver** - Webhook consumer on port 9090
- **db** - PostgreSQL on port 5432

## Test Workflow

### Step 1: Register a webhook

Register the receiver to get notified of order events:

```bash
curl -s -X POST http://localhost:8080/webhooks \
  -H "Content-Type: application/json" \
  -d '{"url":"http://receiver:9090/webhook"}' | jq
```

### Step 2: Create an order

```bash
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"item":"Laptop","amount":999.99}' | jq
```

### Step 3: Check receiver for events

```bash
curl -s http://localhost:9090/events | jq
```

You should see the `order.created` event with the order data.

### Step 4: Update order status

```bash
curl -s -X PUT http://localhost:8080/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{"status":"shipped"}' | jq
```

### Step 5: Check events again

```bash
curl -s http://localhost:9090/events | jq
```

You should now see both `order.created` and `order.shipped` events.

### Additional commands

List all orders:

```bash
curl -s http://localhost:8080/orders | jq
```

List registered webhooks:

```bash
curl -s http://localhost:8080/webhooks | jq
```

Delete a webhook:

```bash
curl -s -X DELETE http://localhost:8080/webhooks/1
```

## Code Walkthrough

### Webhook Registration

The sender stores webhook URLs in a `webhooks` table. When a client registers a webhook, they provide a URL and optionally which events to listen for:

```go
r.Post("/webhooks", registerWebhook)
```

### HMAC Signing

Every webhook delivery is signed with HMAC-SHA256. The sender computes the signature over the JSON payload:

```go
mac := hmac.New(sha256.New, []byte(webhookSecret))
mac.Write(body)
signature := hex.EncodeToString(mac.Sum(nil))
```

The signature is sent in the `X-Webhook-Signature` header.

### Retry Logic

Failed deliveries are retried up to 3 times with increasing delays (2s, 4s):

```go
for attempt := 1; attempt <= maxRetries; attempt++ {
    // ... send request ...
    if err == nil && statusCode >= 200 && statusCode < 300 {
        return
    }
    time.Sleep(time.Duration(attempt) * 2 * time.Second)
}
```

### Signature Verification

The receiver verifies the HMAC signature to ensure the payload was not tampered with:

```go
signature := r.Header.Get("X-Webhook-Signature")
mac := hmac.New(sha256.New, []byte(webhookSecret))
mac.Write(body)
expectedSig := hex.EncodeToString(mac.Sum(nil))
valid := hmac.Equal([]byte(signature), []byte(expectedSig))
```

## Exercises

1. **Add event filtering** - Update the sender to check the `events` field when delivering webhooks. Only deliver events that match the registered pattern (e.g., `order.created` or `order.*`).

2. **Add exponential backoff with jitter** - Replace the linear retry delay with exponential backoff plus random jitter to prevent thundering herd problems.

3. **Add a webhook delivery log endpoint** - Create a `GET /webhooks/{id}/logs` endpoint on the sender that returns the delivery history from the `webhook_logs` table.

4. **Add webhook payload schema versioning** - Include a `version` field in the webhook payload and support multiple payload formats for backward compatibility.

## Key Concepts

| Concept | Description |
|---------|-------------|
| **Webhooks** | HTTP callbacks that notify external services when events occur. The sender POSTs data to a URL registered by the receiver. |
| **HMAC Signatures** | Hash-based Message Authentication Codes ensure payload integrity. The sender and receiver share a secret key used to sign and verify payloads. |
| **Retry Logic** | Failed webhook deliveries are retried with increasing delays. This handles temporary network issues and receiver downtime. |
| **Event-Driven Architecture** | Services communicate through events rather than direct API calls. This decouples the sender from the receiver and allows multiple consumers. |

## Cleanup

Stop and remove all containers:

```bash
docker-compose down -v
```
