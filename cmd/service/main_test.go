package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nats-pg-service/internal/cache"

	"github.com/gin-gonic/gin"
)

func TestMessageHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	cache := cache.NewCache()

	cache.Set("test-1", map[string]string{"value": "Hello"})

	r.GET("/message/:id", func(c *gin.Context) {
		id := c.Param("id")
		if data, ok := cache.Get(id); ok {
			c.JSON(http.StatusOK, data)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		}
	})

	t.Run("Valid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/message/test-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["value"] != "Hello" {
			t.Errorf("Expected value 'Hello', got %v", response)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/message/wrong-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["error"] != "Message not found" {
			t.Errorf("Expected error 'Message not found', got %v", response)
		}
	})
}
