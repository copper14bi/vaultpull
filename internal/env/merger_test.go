package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	return p
}

func TestMerge_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "DB_HOST=old-host\nKEEP=me\n")

	incoming := map[string]string{
		"DB_HOST": "new-host",
		"API_KEY": "secret",
	}

	result, err := Merge(path, incoming, MergeOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["DB_HOST"] != "new-host" {
		t.Errorf("expected new-host, got %s", result["DB_HOST"])
	}
	if result["API_KEY"] != "secret" {
		t.Errorf("expected secret, got %s", result["API_KEY"])
	}
	// KEEP should not appear since it is not in incoming.
	if _, ok := result["KEEP"]; ok {
		t.Error("KEEP should not be present in overwrite mode")
	}
}

func TestMerge_KeepExisting(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "DB_HOST=old-host\n")

	incoming := map[string]string{
		"DB_HOST": "new-host",
		"API_KEY": "secret",
	}

	result, err := Merge(path, incoming, MergeKeepExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["DB_HOST"] != "old-host" {
		t.Errorf("expected old-host to be preserved, got %s", result["DB_HOST"])
	}
	if result["API_KEY"] != "secret" {
		t.Errorf("expected secret, got %s", result["API_KEY"])
	}
}

func TestMerge_MissingFile(t *testing.T) {
	incoming := map[string]string{"KEY": "value"}

	result, err := Merge("/nonexistent/.env", incoming, MergeOverwrite)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if result["KEY"] != "value" {
		t.Errorf("expected value, got %s", result["KEY"])
	}
}

func TestParseEnvFile_IgnoresComments(t *testing.T) {
	dir := t.TempDir()
	path := writeEnvFile(t, dir, "# comment\nKEY=val\n\nOTHER=123\n")

	result, err := parseEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}
