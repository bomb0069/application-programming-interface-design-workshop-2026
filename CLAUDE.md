# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

API Design Workshop — a hands-on learning path for API design, starting from simple RESTful APIs and progressively introducing more advanced API technologies.

## Structure

Each lab lives in its own folder following the naming convention `labXX-topic-name` (e.g., `lab01-simple-rest-api`). Most labs contain a `docker-compose.yml` at the lab root. Labs with sub-labs (e.g., `lab11-api-versioning`) contain multiple sub-directories, each with its own `golang/` and `dotnet/` implementations and their own `docker-compose.yml`.

## Running a Lab

```bash
# Standard lab
cd labXX-topic-name
docker compose up --build

# Sub-lab (e.g., lab11)
cd lab11-api-versioning/lab11-01-url-path-versioning/golang
docker compose up --build
```

## Learning Path

The labs are ordered by complexity, designed for beginners and progressing to advanced topics:

1. **RESTful API** — simple API design with Swagger/OpenAPI documentation, growing in complexity across multiple labs
2. **GraphQL**
3. **Webhook**
4. **WebSocket**
5. **gRPC**
6. **Message Queue**
7. **MQTT**

## Conventions

- Each lab is self-contained with its own `docker-compose.yml` (at the lab root, or inside `golang/`/`dotnet/` for multi-language sub-labs)
- Lab numbering (`labXX`) defines the learning order
- Topic name in the folder slug should be kebab-case and descriptive
