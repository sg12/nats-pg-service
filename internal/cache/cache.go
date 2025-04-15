package cache

import "sync"

// Cache хранит данные в памяти
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewCache создаёт новый кэш
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Set сохраняет значение по ключу
func (c *Cache) Set(id string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[id] = value
}

// Get возвращает значение по ключу
func (c *Cache) Get(id string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[id]
	return v, ok
}
