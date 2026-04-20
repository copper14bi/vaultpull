// Package cache provides a simple TTL-based local cache for Vault secrets
// to reduce redundant API calls during sync operations.
package cache

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry holds a cached secret value and its expiry time.
type Entry struct {
	Data      map[string]string `json:"data"`
	FetchedAt time.Time         `json:"fetched_at"`
	TTL       time.Duration     `json:"ttl"`
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e Entry) IsExpired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a thread-safe in-memory store for secret entries.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given default TTL.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores secrets for the given path.
func (c *Cache) Set(path string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = Entry{
		Data:      data,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Get retrieves secrets for the given path. Returns nil, false if missing or expired.
func (c *Cache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	if !ok || e.IsExpired() {
		return nil, false
	}
	return e.Data, true
}

// Invalidate removes the cached entry for the given path.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// SaveToFile persists the current cache entries to a JSON file.
func (c *Cache) SaveToFile(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c.entries)
}

// LoadFromFile restores cache entries from a JSON file, skipping expired entries.
func (c *Cache) LoadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	var entries map[string]Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range entries {
		if !v.IsExpired() {
			c.entries[k] = v
		}
	}
	return nil
}
