package rotation

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTmp(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRotate_CreatesBackup(t *testing.T) {
	tmp := t.TempDir()
	src := writeTmp(t, tmp, ".env", "SECRET=hello")
	bkDir := filepath.Join(tmp, "backups")

	r := New(bkDir, 5)
	if err := r.Rotate(src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(bkDir, ".env.*.bak"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(matches))
	}

	data, _ := os.ReadFile(matches[0])
	if string(data) != "SECRET=hello" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	tmp := t.TempDir()
	src := writeTmp(t, tmp, ".env", "SECRET=val")
	bkDir := filepath.Join(tmp, "backups")
	os.MkdirAll(bkDir, 0o700)

	// Pre-create 5 old backups with sorted names.
	for i := 0; i < 5; i++ {
		ts := time.Now().UTC().Add(time.Duration(-5+i) * time.Second).Format("20060102T150405Z")
		name := filepath.Join(bkDir, ".env."+ts+".bak")
		os.WriteFile(name, []byte("old"), 0o600)
		time.Sleep(1 * time.Millisecond)
	}

	r := New(bkDir, 5)
	if err := r.Rotate(src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(bkDir, ".env.*.bak"))
	if len(matches) != 5 {
		t.Errorf("expected 5 backups after pruning, got %d", len(matches))
	}
}

func TestRotate_MissingSource(t *testing.T) {
	tmp := t.TempDir()
	r := New(filepath.Join(tmp, "bk"), 3)
	err := r.Rotate(filepath.Join(tmp, "nonexistent.env"))
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestNew_DefaultMaxBackups(t *testing.T) {
	r := New("/tmp", 0)
	if r.MaxBackups != 5 {
		t.Errorf("expected default 5, got %d", r.MaxBackups)
	}
}
