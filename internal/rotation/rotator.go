package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rotator manages backup rotation for .env files.
type Rotator struct {
	MaxBackups int
	BackupDir  string
}

// New creates a Rotator with sensible defaults.
func New(backupDir string, maxBackups int) *Rotator {
	if maxBackups <= 0 {
		maxBackups = 5
	}
	return &Rotator{MaxBackups: maxBackups, BackupDir: backupDir}
}

// Rotate copies src to a timestamped backup file and prunes old backups.
func (r *Rotator) Rotate(src string) error {
	if err := os.MkdirAll(r.BackupDir, 0o700); err != nil {
		return fmt.Errorf("rotation: mkdir %s: %w", r.BackupDir, err)
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("rotation: read %s: %w", src, err)
	}

	base := filepath.Base(src)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dest := filepath.Join(r.BackupDir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return fmt.Errorf("rotation: write backup %s: %w", dest, err)
	}

	return r.pruneOldBackups(base)
}

// pruneOldBackups removes oldest backups exceeding MaxBackups for a given base name.
func (r *Rotator) pruneOldBackups(base string) error {
	pattern := filepath.Join(r.BackupDir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) <= r.MaxBackups {
		return nil
	}
	// Glob returns sorted order; remove oldest (first entries).
	for _, old := range matches[:len(matches)-r.MaxBackups] {
		if err := os.Remove(old); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rotation: remove %s: %w", old, err)
		}
	}
	return nil
}
