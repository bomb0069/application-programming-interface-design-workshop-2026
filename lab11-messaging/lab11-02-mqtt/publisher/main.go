package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
		SetClientID("sensor-publisher").
		SetAutoReconnect(true).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			log.Printf("Connection lost: %v", err)
		}).
		SetOnConnectHandler(func(c mqtt.Client) {
			log.Println("Connected to MQTT broker")
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

	sensors := []string{"sensor-01", "sensor-02", "sensor-03"}

	// Publish sensor data periodically
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Publisher started. Sending sensor data every 3 seconds...")

	for {
		select {
		case <-ticker.C:
			for _, sensorID := range sensors {
				data := SensorData{
					SensorID:    sensorID,
					Temperature: 20 + rand.Float64()*15,
					Humidity:    40 + rand.Float64()*40,
					Timestamp:   time.Now().UTC().Format(time.RFC3339),
				}

				payload, _ := json.Marshal(data)
				topic := fmt.Sprintf("sensors/%s/data", sensorID)

				// QoS 1 = At least once delivery
				token := client.Publish(topic, 1, false, payload)
				token.Wait()

				log.Printf("Published to %s: temp=%.1f°C humidity=%.1f%%",
					topic, data.Temperature, data.Humidity)
			}

			// Also publish an alert if temperature is high
			for _, sensorID := range sensors {
				temp := 20 + rand.Float64()*20
				if temp > 32 {
					alert := map[string]interface{}{
						"sensor_id": sensorID,
						"type":      "high_temperature",
						"value":     temp,
						"threshold": 32.0,
						"timestamp": time.Now().UTC().Format(time.RFC3339),
					}
					payload, _ := json.Marshal(alert)
					client.Publish("sensors/alerts", 2, false, payload) // QoS 2 = Exactly once
					log.Printf("ALERT published for %s: temp=%.1f°C", sensorID, temp)
				}
			}

		case <-sigChan:
			log.Println("Shutting down publisher...")
			return
		}
	}
}
