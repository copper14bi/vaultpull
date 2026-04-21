package watch_test

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/watch"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpull-watch-*.yaml")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0o600)
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTmp(t, "initial")

	var calls int64
	w := watch.New(path, 20*time.Millisecond, func(p string) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- w.Start(ctx) }()

	time.Sleep(50 * time.Millisecond)
	_ = os.WriteFile(path, []byte("updated"), 0o600)

	<-done

	if atomic.LoadInt64(&calls) == 0 {
		t.Error("expected handler to be called at least once")
	}
}

func TestWatcher_NoCallWhenUnchanged(t *testing.T) {
	path := writeTmp(t, "stable")

	var calls int64
	w := watch.New(path, 20*time.Millisecond, func(p string) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- w.Start(ctx) }()
	<-done

	if atomic.LoadInt64(&calls) != 0 {
		t.Errorf("expected 0 handler calls, got %d", atomic.LoadInt64(&calls))
	}
}

func TestWatcher_MissingFileReturnsError(t *testing.T) {
	w := watch.New("/nonexistent/path.yaml", 10*time.Millisecond, func(string) error { return nil })
	err := w.Start(context.Background())
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
