package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func TestNew_CopiesSecrets(t *testing.T) {
	orig := map[string]string{"KEY": "val"}
	s := snapshot.New("secret/app", orig)
	orig["KEY"] = "mutated"
	if s.Secrets["KEY"] != "val" {
		t.Errorf("expected original value, got %s", s.Secrets["KEY"])
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "snap.json")
	s := snapshot.New("secret/app", map[string]string{"A": "1", "B": "2"})
	if err := s.Save(tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Path != "secret/app" {
		t.Errorf("path mismatch: %s", loaded.Path)
	}
	if loaded.Secrets["A"] != "1" || loaded.Secrets["B"] != "2" {
		t.Errorf("secrets mismatch: %v", loaded.Secrets)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New("secret/app", map[string]string{})
	err := s.Save("/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
}

func TestCompare_Added(t *testing.T) {
	s := snapshot.New("p", map[string]string{"A": "1"})
	d := s.Compare(map[string]string{"A": "1", "B": "2"})
	if len(d.Added) != 1 || d.Added[0] != "B" {
		t.Errorf("expected B added, got %v", d.Added)
	}
}

func TestCompare_Removed(t *testing.T) {
	s := snapshot.New("p", map[string]string{"A": "1", "B": "2"})
	d := s.Compare(map[string]string{"A": "1"})
	if len(d.Removed) != 1 || d.Removed[0] != "B" {
		t.Errorf("expected B removed, got %v", d.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	s := snapshot.New("p", map[string]string{"A": "old"})
	d := s.Compare(map[string]string{"A": "new"})
	if len(d.Changed) != 1 || d.Changed[0] != "A" {
		t.Errorf("expected A changed, got %v", d.Changed)
	}
}

func TestCompare_NoDrift(t *testing.T) {
	s := snapshot.New("p", map[string]string{"A": "1"})
	d := s.Compare(map[string]string{"A": "1"})
	if d.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestHasDrift_True(t *testing.T) {
	s := snapshot.New("p", map[string]string{})
	d := s.Compare(map[string]string{"NEW": "val"})
	if !d.HasDrift() {
		t.Error("expected drift")
	}
}

func init() {
	// ensure temp dir cleanup works cross-platform
	_ = os.MkdirTemp
	_ = filepath.Join
}
