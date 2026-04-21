// Package watch provides file system watching to trigger automatic
// secret re-sync when the vaultpull config file changes.
package watch

import (
	"context"
	"log"
	"os"
	"time"
)

// Handler is called when a watched file changes.
type Handler func(path string) error

// Watcher monitors a file for modifications and invokes a Handler.
type Watcher struct {
	path     string
	interval time.Duration
	handler  Handler
	logger   *log.Logger
}

// New creates a Watcher with the given file path, poll interval, and change handler.
func New(path string, interval time.Duration, handler Handler) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		handler:  handler,
		logger:   log.New(os.Stderr, "[watch] ", log.LstdFlags),
	}
}

// Start begins polling the file for changes. It blocks until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) error {
	info, err := os.Stat(w.path)
	if err != nil {
		return err
	}
	lastMod := info.ModTime()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Printf("watching %s every %s", w.path, w.interval)

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("watcher stopped")
			return ctx.Err()
		case <-ticker.C:
			info, err := os.Stat(w.path)
			if err != nil {
				w.logger.Printf("stat error: %v", err)
				continue
			}
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				w.logger.Printf("change detected in %s", w.path)
				if err := w.handler(w.path); err != nil {
					w.logger.Printf("handler error: %v", err)
				}
			}
		}
	}
}
