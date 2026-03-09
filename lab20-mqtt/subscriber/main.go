package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	SensorID    string  `json:"sensor_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}

func main() {
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("sensor-subscriber").
		SetAutoReconnect(true).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			log.Printf("Connection lost: %v", err)
		}).
		SetOnConnectHandler(func(c mqtt.Client) {
			log.Println("Connected to MQTT broker")
			subscribe(c)
		})

	client := mqtt.NewClient(opts)
	for i := 0; i < 30; i++ {
		token := client.Connect()
		token.Wait()
		if token.Error() == nil {
			break
		}
		log.Printf("Waiting for MQTT broker... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if !client.IsConnected() {
		log.Fatal("Failed to connect to MQTT broker")
	}
	defer client.Disconnect(250)

	log.Println("Subscriber started. Listening for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down subscriber...")
}

func subscribe(client mqtt.Client) {
	// Subscribe to all sensor data using wildcard
	// sensors/+/data matches sensors/sensor-01/data, sensors/sensor-02/data, etc.
	client.Subscribe("sensors/+/data", 1, func(c mqtt.Client, msg mqtt.Message) {
		var data SensorData
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Error parsing message: %v", err)
			return
		}
		fmt.Printf("[DATA] %s | Sensor: %s | Temp: %.1f°C | Humidity: %.1f%% | Topic: %s\n",
			data.Timestamp, data.SensorID, data.Temperature, data.Humidity, msg.Topic())
	})

	// Subscribe to alerts with QoS 2
	client.Subscribe("sensors/alerts", 2, func(c mqtt.Client, msg mqtt.Message) {
		var alert map[string]interface{}
		json.Unmarshal(msg.Payload(), &alert)
		fmt.Printf("[ALERT] %s | Sensor: %s | Type: %s | Value: %.1f | Threshold: %.0f\n",
			alert["timestamp"], alert["sensor_id"], alert["type"],
			alert["value"].(float64), alert["threshold"].(float64))
	})

	// Subscribe to ALL sensors topics using multi-level wildcard
	// sensors/# matches everything under sensors/
	client.Subscribe("sensors/#", 0, func(c mqtt.Client, msg mqtt.Message) {
		log.Printf("[ALL] Topic: %s | Payload size: %d bytes", msg.Topic(), len(msg.Payload()))
	})
}
