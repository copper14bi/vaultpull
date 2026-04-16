package audit

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestRecord_WritesEntry(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "audit-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	l, err := NewLogger(tmp.Name())
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	defer l.Close()

	if err := l.Record("sync", "secret/app", ".env", true, nil); err != nil {
		t.Fatalf("Record: %v", err)
	}
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	var entry Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if entry.Event != "sync" {
		t.Errorf("expected event 'sync', got %q", entry.Event)
	}
	if !entry.Success {
		t.Error("expected success=true")
	}
	if entry.Error != "" {
		t.Errorf("expected no error field, got %q", entry.Error)
	}
}

func TestRecord_CapturesError(t *testing.T) {
	tmp, _ := os.CreateTemp(t.TempDir(), "audit-*.log")
	tmp.Close()

	l, _ := NewLogger(tmp.Name())
	defer l.Close()

	l.Record("sync", "secret/app", ".env", false, errors.New("permission denied"))
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	var entry Entry
	json.NewDecoder(f).Decode(&entry)
	if entry.Error != "permission denied" {
		t.Errorf("expected error field, got %q", entry.Error)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	tmp, _ := os.CreateTemp(t.TempDir(), "audit-*.log")
	tmp.Close()

	l, _ := NewLogger(tmp.Name())
	l.Record("sync", "secret/a", ".env", true, nil)
	l.Record("rotate", "secret/b", ".env.bak", true, nil)
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 log lines, got %d", count)
	}
}
