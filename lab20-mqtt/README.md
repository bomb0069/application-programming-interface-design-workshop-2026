# Lab 20 - MQTT

IoT-style pub/sub messaging with an MQTT broker (Mosquitto) in Go.

## Learning Objectives

- Understand IoT-style pub/sub messaging with MQTT
- Configure and run an MQTT broker (Mosquitto)
- Publish and subscribe to MQTT topics in Go
- Work with MQTT QoS levels (0, 1, 2)
- Use MQTT wildcards (`+` single-level and `#` multi-level)

## Architecture

```
Publisher (IoT Sensors) ---> Mosquitto Broker (:1883) ---> Subscriber
```

The publisher simulates three IoT sensors that periodically send temperature and humidity readings. The subscriber listens using wildcard subscriptions and processes incoming data. The Mosquitto broker handles message routing between publishers and subscribers.

## Getting Started

Start all services:

```bash
docker-compose up --build
```

Watch subscriber output:

```bash
docker-compose logs -f subscriber
```

## Topics Used

| Topic                        | Description                    | QoS |
|------------------------------|--------------------------------|-----|
| `sensors/{sensor_id}/data`   | Sensor readings                | 1   |
| `sensors/alerts`             | High temperature alerts        | 2   |

## Test Manually with Mosquitto CLI

Subscribe to all sensor topics from inside the broker container:

```bash
docker-compose exec mosquitto mosquitto_sub -t "sensors/#" -v
```

Publish a test message:

```bash
docker-compose exec mosquitto mosquitto_pub -t "sensors/test/data" -m '{"sensor_id":"test","temperature":25.5}'
```

## MQTT vs RabbitMQ

| Feature            | MQTT (Mosquitto)                  | RabbitMQ (AMQP)                    |
|--------------------|-----------------------------------|------------------------------------|
| Protocol           | MQTT (lightweight)                | AMQP (feature-rich)               |
| Design Goal        | IoT, constrained devices          | Enterprise messaging               |
| Broker Complexity  | Simple                            | Complex (exchanges, bindings)      |
| Message Routing    | Topic-based only                  | Exchange types (direct, fanout, topic, headers) |
| QoS Levels         | 0, 1, 2                           | Acknowledgments + confirms         |
| Message Size       | Optimized for small payloads      | No specific optimization           |
| Wildcards          | `+` (single) and `#` (multi)      | `*` (single) and `#` (multi)       |
| Retained Messages  | Built-in                          | Not native (plugin required)       |
| Last Will (LWT)    | Built-in                          | Not native                         |
| Best For           | IoT, sensors, mobile              | Microservices, task queues          |

## QoS Levels

| Level | Name            | Description                                                                 |
|-------|-----------------|-----------------------------------------------------------------------------|
| 0     | At most once    | Fire and forget. No acknowledgment. Messages may be lost.                   |
| 1     | At least once   | Acknowledged delivery. Messages may be delivered more than once.            |
| 2     | Exactly once    | Assured delivery with a four-step handshake. Highest overhead.              |

## MQTT Wildcards

MQTT uses two wildcard characters for topic subscriptions:

### `+` Single-Level Wildcard

Matches exactly one topic level.

```
sensors/+/data
```

Matches:
- `sensors/sensor-01/data`
- `sensors/sensor-02/data`
- `sensors/any-id/data`

Does NOT match:
- `sensors/data`
- `sensors/floor1/sensor-01/data`

### `#` Multi-Level Wildcard

Matches zero or more topic levels. Must be the last character in the subscription.

```
sensors/#
```

Matches:
- `sensors`
- `sensors/sensor-01/data`
- `sensors/alerts`
- `sensors/floor1/sensor-01/data`

## Code Walkthrough

### MQTT Client (paho.mqtt.golang)

The lab uses the Eclipse Paho MQTT Go client library.

**Connect options:**

```go
opts := mqtt.NewClientOptions().
    AddBroker(broker).
    SetClientID("sensor-publisher").
    SetAutoReconnect(true).
    SetConnectionLostHandler(func(c mqtt.Client, err error) {
        log.Printf("Connection lost: %v", err)
    }).
    SetOnConnectHandler(func(c mqtt.Client) {
        log.Println("Connected to MQTT broker")
    })
```

**Publishing messages:**

```go
token := client.Publish(topic, 1, false, payload) // topic, QoS, retained, payload
token.Wait()
```

**Subscribing with callbacks:**

```go
client.Subscribe("sensors/+/data", 1, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Topic: %s | Payload: %s\n", msg.Topic(), msg.Payload())
})
```

**Wildcard subscriptions:**

```go
// Single-level wildcard: matches one level
client.Subscribe("sensors/+/data", 1, handler)

// Multi-level wildcard: matches everything under sensors/
client.Subscribe("sensors/#", 0, handler)
```

## Exercises

1. **Retained Messages** - Modify the publisher to send retained messages for the last known state of each sensor. When a new subscriber connects, it should immediately receive the latest reading for each sensor.

2. **Last Will and Testament (LWT)** - Add a Last Will and Testament message to the publisher so that when a sensor disconnects unexpectedly, the broker automatically publishes a status message to `sensors/{sensor_id}/status` with a payload of `offline`.

3. **Dashboard Subscriber** - Create a new subscriber that calculates running averages of temperature and humidity per sensor and periodically prints a summary table.

4. **Device Commands** - Implement bidirectional communication by having the publisher also subscribe to `sensors/{id}/commands`. Create a separate command publisher that sends control messages (e.g., change reporting interval, recalibrate sensor).

## Key Concepts

| Concept            | Description                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| MQTT Protocol      | Lightweight publish/subscribe messaging protocol designed for IoT           |
| Pub/Sub Pattern    | Decoupled messaging where publishers and subscribers communicate via topics |
| QoS Levels         | Quality of Service guarantees (0: at most once, 1: at least once, 2: exactly once) |
| Topic Wildcards    | `+` for single-level and `#` for multi-level topic matching                |
| Retained Messages  | Broker stores the last message on a topic for new subscribers              |

## Cleanup

Stop and remove all containers and volumes:

```bash
docker-compose down -v
```
