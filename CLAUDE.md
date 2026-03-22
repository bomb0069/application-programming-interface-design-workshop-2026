# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

API Design Workshop — a hands-on learning path for API design, starting from simple RESTful APIs and progressively introducing more advanced API technologies.

## Structure

Labs use `labXX-YY-topic-name` naming where XX = group number, YY = sub-lab number. Each group folder (`labXX-group-name/`) contains sub-lab folders, each with a `golang/` (and optionally `dotnet/`) directory and its own `docker-compose.yml`.

> **NOTE:** The repo is being restructured from flat `labXX` to grouped `labXX-YY` format. See [workshop-reorganization-plan.md](workshop-reorganization-plan.md) for the current state and mapping.

## Running a Lab

```bash
cd labXX-group-name/labXX-YY-topic/golang
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

## Implementation Plans

Before starting new work, **read the reorganization plan first** — it defines the new lab numbering and what needs to be built. Update checkboxes (`- [ ]` → `- [x]`) as you complete each item.

| Plan | Status | Description |
|------|--------|-------------|
| [workshop-reorganization-plan.md](workshop-reorganization-plan.md) | PENDING | **START HERE** — Master plan: restructure all labs into `labXX-YY` groups, rename 20 existing labs, build 16 new labs |
| [lab21-search-api-plan.md](lab21-search-api-plan.md) | PENDING | Detailed spec for Group 05 (Advanced Search) sub-labs — endpoints, DB schemas, code patterns |
| [workshop-gap-analysis-plan.md](workshop-gap-analysis-plan.md) | PENDING | Gap analysis vs 2-day course — detailed scope for each new lab topic |
