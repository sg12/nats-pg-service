package main

import (
	"encoding/json"
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

	// Отправляем тестовые сообщения в формате model.json
	for i := 0; i < 5; i++ {
		msg := map[string]interface{}{
			"id":           "b563feb7b2b84b6test" + string(rune('0'+i)), // Уникальный order_uid
			"track_number": "WBILMTESTTRACK" + string(rune('0'+i)),
			"entry":        "WBIL",
			"delivery": map[string]interface{}{
				"name":    "Test Testov",
				"phone":   "+9720000000",
				"zip":     "2639809",
				"city":    "Kiryat Mozkin",
				"address": "Ploshad Mira 15",
				"region":  "Kraiot",
				"email":   "test@gmail.com",
			},
			"payment": map[string]interface{}{
				"transaction":   "b563feb7b2b84b6test" + string(rune('0'+i)),
				"request_id":    "",
				"currency":      "USD",
				"provider":      "wbpay",
				"amount":        1817,
				"payment_dt":    1637907727,
				"bank":          "alpha",
				"delivery_cost": 1500,
				"goods_total":   317,
				"custom_fee":    0,
			},
			"items": []map[string]interface{}{
				{
					"chrt_id":      9934930,
					"track_number": "WBILMTESTTRACK" + string(rune('0'+i)),
					"price":        453,
					"rid":          "ab4219087a764ae0btest" + string(rune('0'+i)),
					"name":         "Mascaras",
					"sale":         30,
					"size":         "0",
					"total_price":  317,
					"nm_id":        2389212,
					"brand":        "Vivienne Sabo",
					"status":       202,
				},
			},
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        "test",
			"delivery_service":   "meest",
			"shardkey":           "9",
			"sm_id":              99,
			"date_created":       "2021-11-26T06:22:19Z",
			"oof_shard":          "1",
		}
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}
		err = sc.Publish("test-channel", data)
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Printf("Published: %s", msg["id"])
		}
		time.Sleep(1 * time.Second)
	}
}
