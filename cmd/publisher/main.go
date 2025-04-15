package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-1", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer sc.Close()

	// Отправляем тестовые сообщения
	for i := 0; i < 5; i++ {
		msg := map[string]interface{}{
			"id":    fmt.Sprintf("msg-%d", i),
			"value": fmt.Sprintf("Test message %d", i),
		}
		data, _ := json.Marshal(msg)
		err := sc.Publish("test-channel", data)
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Printf("Published: %s", msg["id"])
		}
		time.Sleep(1 * time.Second)
	}
}
