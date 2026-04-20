package cache_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/cache"
)

// TestRoundTrip_MultipleKeys verifies that multiple paths survive a
// save/load cycle without data loss or cross-contamination.
func TestRoundTrip_MultipleKeys(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "cache.json")

	secrets := map[string]map[string]string{
		"secret/app/db":  {"DB_PASS": "hunter2", "DB_USER": "admin"},
		"secret/app/api": {"API_KEY": "deadbeef"},
	}

	c1 := cache.New(10 * time.Minute)
	for path, data := range secrets {
		c1.Set(path, data)
	}
	if err := c1.SaveToFile(file); err != nil {
		t.Fatalf("save: %v", err)
	}

	c2 := cache.New(10 * time.Minute)
	if err := c2.LoadFromFile(file); err != nil {
		t.Fatalf("load: %v", err)
	}

	for path, want := range secrets {
		got, ok := c2.Get(path)
		if !ok {
			t.Errorf("path %q: expected hit after round-trip", path)
			continue
		}
		for k, v := range want {
			if got[k] != v {
				t.Errorf("path %q key %q: got %q, want %q", path, k, got[k], v)
			}
		}
	}
}

// TestConcurrentAccess ensures Set/Get are safe under concurrent use.
func TestConcurrentAccess(t *testing.T) {
	c := cache.New(5 * time.Minute)
	done := make(chan struct{})

	for i := 0; i < 20; i++ {
		go func(i int) {
			c.Set("secret/concurrent", map[string]string{"N": "v"})
			c.Get("secret/concurrent")
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
