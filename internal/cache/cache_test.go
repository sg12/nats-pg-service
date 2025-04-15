package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	c := NewCache()

	// Тест добавления и получения
	c.Set("key1", map[string]string{"value": "test"})
	if data, ok := c.Get("key1"); !ok {
		t.Error("Expected to find key1")
	} else if data.(map[string]string)["value"] != "test" {
		t.Errorf("Expected value 'test', got %v", data)
	}

	// Тест отсутствующего ключа
	if _, ok := c.Get("key2"); ok {
		t.Error("Expected key2 to be missing")
	}

	// Тест перезаписи
	c.Set("key1", map[string]string{"value": "updated"})
	if data, ok := c.Get("key1"); !ok {
		t.Error("Expected to find key1")
	} else if data.(map[string]string)["value"] != "updated" {
		t.Errorf("Expected value 'updated', got %v", data)
	}
}
