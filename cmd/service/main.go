package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/stan.go"
)

// Cache хранит данные в памяти
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Set(id string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[id] = value
}

func (c *Cache) Get(id string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[id]
	return v, ok
}

func main() {
	// Подключение к PostgreSQL
	ctx := context.Background()
	connStr := "postgres://test_user:pass@localhost:5433/test_db"
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close(ctx)

	// Подключение к NATS Streaming
	sc, err := stan.Connect("test-cluster", "client-1", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer sc.Close()

	// Инициализация кэша
	cache := NewCache()

	// Восстановление кэша из БД
	rows, err := db.Query(ctx, "SELECT message_id, data FROM messages")
	if err != nil {
		log.Fatalf("Failed to restore cache: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var messageID string
		var data json.RawMessage
		if err := rows.Scan(&messageID, &data); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		cache.Set(messageID, data)
	}
	log.Printf("Restored %d messages to cache", len(cache.data))

	// Подписка на канал
	_, err = sc.Subscribe("test-channel", func(m *stan.Msg) {
		var msgData map[string]interface{}
		if err := json.Unmarshal(m.Data, &msgData); err != nil {
			log.Printf("Invalid JSON: %v", err)
			return
		}

		messageID, ok := msgData["id"].(string)
		if !ok {
			log.Println("Missing or invalid ID")
			return
		}

		// Сохранение в БД
		_, err := db.Exec(ctx, `
            INSERT INTO messages (message_id, data)
            VALUES ($1, $2)
            ON CONFLICT (message_id) DO NOTHING`,
			messageID, m.Data)
		if err != nil {
			log.Printf("Failed to save to DB: %v", err)
			return
		}

		// Сохранение в кэш
		cache.Set(messageID, msgData)
		log.Printf("Processed message: %s, Cache size: %d", messageID, len(cache.data))
	}, stan.DurableName("durable"))
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Subscribed to test-channel")

	// Настройка HTTP-сервера
	r := gin.Default()

	// Эндпоинт для получения данных по ID
	r.GET("/message/:id", func(c *gin.Context) {
		id := c.Param("id")
		if data, ok := cache.Get(id); ok {
			c.JSON(http.StatusOK, data)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		}
	})

	// Простой UI
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.LoadHTMLGlob("web/templates/*")

	// Запуск сервера
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	// Держим сервис активным
	select {}
}
