# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

API Design Workshop — a hands-on learning path for API design, starting from simple RESTful APIs and progressively introducing more advanced API technologies.

## Structure

Each lab lives in its own folder following the naming convention `labXX-topic-name` (e.g., `lab01-simple-rest-api`). Each lab folder contains a `docker-compose.yml` to start all required tools and services for that lab.

## Running a Lab

```bash
cd labXX-topic-name
docker-compose up
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

- Each lab is self-contained with its own `docker-compose.yml`
- Lab numbering (`labXX`) defines the learning order
- Topic name in the folder slug should be kebab-case and descriptive
