package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")

	w := NewWriter(outPath, "")
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWrite_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")
	backupDir := filepath.Join(dir, "backups")

	// Pre-create the .env file to trigger backup logic.
	if err := os.WriteFile(outPath, []byte("OLD=value\n"), 0600); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	w := NewWriter(outPath, backupDir)
	if err := w.Write(map[string]string{"NEW": "value"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("backup dir not created: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 backup file, got %d", len(entries))
	}
}

func TestQuoteValue_NoQuotesNeeded(t *testing.T) {
	if got := quoteValue("simple"); got != "simple" {
		t.Errorf("expected simple, got %s", got)
	}
}

func TestQuoteValue_SpaceRequiresQuotes(t *testing.T) {
	got := quoteValue("hello world")
	if got != `"hello world"` {
		t.Errorf("expected quoted value, got %s", got)
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")

	w := NewWriter(outPath, "")
	if err := w.Write(map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
