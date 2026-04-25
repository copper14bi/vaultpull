package multienv_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/multienv"
)

func TestWriteAll_CreatesFiles(t *testing.T) {
	dir := t.TempDir()
	fileA := filepath.Join(dir, "app.env")
	fileB := filepath.Join(dir, "db.env")

	targets := []multienv.Target{
		{OutputFile: fileA, Prefixes: []string{"APP_"}},
		{OutputFile: fileB, Prefixes: []string{"DB_"}},
	}

	secrets := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.internal",
		"OTHER":    "ignored",
	}

	w := multienv.New(targets)
	counts, err := w.WriteAll(secrets)
	if err != nil {
		t.Fatalf("WriteAll error: %v", err)
	}

	if counts[fileA] != 2 {
		t.Errorf("expected 2 keys in app.env, got %d", counts[fileA])
	}
	if counts[fileB] != 1 {
		t.Errorf("expected 1 key in db.env, got %d", counts[fileB])
	}

	contentA, _ := os.ReadFile(fileA)
	if !strings.Contains(string(contentA), "APP_HOST") {
		t.Error("app.env missing APP_HOST")
	}
	contentB, _ := os.ReadFile(fileB)
	if !strings.Contains(string(contentB), "DB_HOST") {
		t.Error("db.env missing DB_HOST")
	}
}

func TestWriteAll_NoPrefixWritesAll(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "all.env")

	targets := []multienv.Target{{OutputFile: out}}
	secrets := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}

	w := multienv.New(targets)
	counts, err := w.WriteAll(secrets)
	if err != nil {
		t.Fatalf("WriteAll error: %v", err)
	}
	if counts[out] != 3 {
		t.Errorf("expected 3 keys, got %d", counts[out])
	}
}

func TestWriteAll_EmptySecretsSkipsFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "empty.env")

	targets := []multienv.Target{{OutputFile: out, Prefixes: []string{"NOMATCH_"}}}
	secrets := map[string]string{"OTHER": "val"}

	w := multienv.New(targets)
	counts, err := w.WriteAll(secrets)
	if err != nil {
		t.Fatalf("WriteAll error: %v", err)
	}
	if counts[out] != 0 {
		t.Errorf("expected 0 keys written, got %d", counts[out])
	}
	if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
		t.Error("expected file not to be created for empty secret set")
	}
}

func TestWriteAll_MultipleTargetsSamePrefix(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "one.env")
	file2 := filepath.Join(dir, "two.env")

	targets := []multienv.Target{
		{OutputFile: file1, Prefixes: []string{"SHARED_"}},
		{OutputFile: file2, Prefixes: []string{"SHARED_"}},
	}
	secrets := map[string]string{"SHARED_KEY": "value"}

	w := multienv.New(targets)
	counts, err := w.WriteAll(secrets)
	if err != nil {
		t.Fatalf("WriteAll error: %v", err)
	}
	if counts[file1] != 1 || counts[file2] != 1 {
		t.Errorf("expected 1 key each, got %d and %d", counts[file1], counts[file2])
	}
}
