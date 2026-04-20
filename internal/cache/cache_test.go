package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGet_MissOnEmpty(t *testing.T) {
	c := New(5 * time.Minute)
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSet_ThenGet(t *testing.T) {
	c := New(5 * time.Minute)
	data := map[string]string{"KEY": "value"}
	c.Set("secret/foo", data)
	got, ok := c.Get("secret/foo")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["KEY"] != "value" {
		t.Errorf("got %q, want %q", got["KEY"], "value")
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("secret/foo", map[string]string{"X": "y"})
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("secret/foo", map[string]string{"A": "b"})
	c.Invalidate("secret/foo")
	_, ok := c.Get("secret/foo")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestSaveAndLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "cache.json")

	c1 := New(5 * time.Minute)
	c1.Set("secret/bar", map[string]string{"TOKEN": "abc123"})
	if err := c1.SaveToFile(file); err != nil {
		t.Fatalf("SaveToFile: %v", err)
	}

	c2 := New(5 * time.Minute)
	if err := c2.LoadFromFile(file); err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}
	got, ok := c2.Get("secret/bar")
	if !ok {
		t.Fatal("expected cache hit after load")
	}
	if got["TOKEN"] != "abc123" {
		t.Errorf("got %q, want %q", got["TOKEN"], "abc123")
	}
}

func TestLoadFromFile_SkipsExpired(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "cache.json")

	expired := map[string]Entry{
		"secret/old": {
			Data:      map[string]string{"K": "v"},
			FetchedAt: time.Now().Add(-10 * time.Minute),
			TTL:       1 * time.Minute,
		},
	}
	f, _ := os.Create(file)
	json.NewEncoder(f).Encode(expired)
	f.Close()

	c := New(5 * time.Minute)
	if err := c.LoadFromFile(file); err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}
	_, ok := c.Get("secret/old")
	if ok {
		t.Fatal("expected expired entry to be skipped on load")
	}
}

func TestLoadFromFile_MissingFile(t *testing.T) {
	c := New(5 * time.Minute)
	if err := c.LoadFromFile("/nonexistent/cache.json"); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}
