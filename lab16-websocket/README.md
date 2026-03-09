# Lab 16 - WebSocket

Real-time bidirectional communication with WebSocket in Go - building a simple chat system.

## Learning Objectives

- Implement a WebSocket server in Go using `gorilla/websocket`
- Understand real-time bidirectional communication between client and server
- Apply the Hub/Client pattern for managing multiple concurrent connections
- Build a simple chat application with live message broadcasting

## Getting Started

```bash
docker-compose up --build
```

Open [http://localhost:8080](http://localhost:8080) in **multiple browser tabs** to simulate different chat users. Enter a username in each tab and start chatting.

## Using from CLI

You can also test with `wscat` (install via `npm install -g wscat`):

```bash
wscat -c "ws://localhost:8080/ws?username=CLI"
```

Then send messages as JSON:

```json
{"content": "Hello from the terminal!"}
```

## Stats Endpoint

Check the stats endpoint to see online users and message count:

```bash
curl http://localhost:8080/stats
```

Example response:

```json
{
  "online_users": 2,
  "usernames": ["Alice", "Bob"],
  "message_count": 5
}
```

## Code Walkthrough

### WebSocket Upgrader

The `websocket.Upgrader` handles the HTTP-to-WebSocket protocol upgrade. `CheckOrigin` is set to allow all origins for development:

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}
```

### Hub Pattern

The `Hub` is the central coordinator that manages all active connections using three channels:

- **register** - New clients connect and are added to the client map
- **unregister** - Disconnected clients are removed and their send channel is closed
- **broadcast** - Messages are fanned out to every connected client

The Hub runs in its own goroutine (`go hub.Run()`) and uses a `select` loop to process events sequentially, avoiding race conditions.

### Client Read/Write Pumps

Each connected client spawns two goroutines:

- **readPump** - Reads incoming messages from the WebSocket connection and forwards them to the Hub's broadcast channel
- **writePump** - Reads from the client's `send` channel and writes messages to the WebSocket connection

This separation ensures that reading and writing never block each other.

### Chat History

The Hub maintains the last 100 messages in memory. When a new client connects, it receives the full history so it can see recent conversation context.

## WebSocket vs REST Comparison

| Feature | REST | WebSocket |
|---------|------|-----------|
| Connection | New connection per request | Persistent connection |
| Direction | Client-initiated only | Bidirectional |
| Overhead | HTTP headers on every request | Minimal framing after handshake |
| Use Case | CRUD operations | Real-time data (chat, games, live feeds) |
| Scaling | Stateless, easy to scale | Stateful, requires sticky sessions or pub/sub |

## Exercises

1. **Add typing indicators** - Broadcast a "user is typing..." message when a user starts typing, and clear it when they stop or send a message.

2. **Add private messaging** - Implement a `/dm username message` command that sends a message only to the specified user.

3. **Add chat rooms/channels** - Allow users to create and join different chat rooms, with messages scoped to the room.

4. **Add message persistence with PostgreSQL** - Store messages in a database so chat history survives server restarts.

## Key Concepts

- **WebSocket Protocol** - A full-duplex communication protocol over a single TCP connection, initiated via an HTTP upgrade handshake.
- **Hub/Client Pattern** - A central hub manages all connections and coordinates message routing; each client has its own read and write goroutines.
- **Goroutines for Concurrent Connections** - Each client connection is handled by dedicated goroutines, enabling the server to manage thousands of simultaneous connections efficiently.
- **Real-time Communication** - Messages are delivered instantly to all connected clients without polling, providing a responsive user experience.

## Cleanup

```bash
docker-compose down
```
